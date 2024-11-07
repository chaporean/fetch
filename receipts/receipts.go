package receipts

import (
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string `json:"retailer,omitempty"`
	PurchaseDate string `json:"purchaseDate,omitempty"`
	PurchaseTime string `json:"purchaseTime,omitempty"`
	Total        string `json:"total,omitempty"`
	Items        []Item `json:"items,omitempty"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var receiptsPoints = make(map[string]int)

func HandleProcessReceipt(c *gin.Context) {
	var receipt Receipt

	if err := c.BindJSON(&receipt); err != nil {
		return
	}
	id := uuid.New()
	points, err := receipt.CalulatePoints()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error calculating points"})
		return
	}
	receiptsPoints[id.String()] = points
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})

}

func HandleGetPoints(c *gin.Context) {
	id := c.Param("id")

	val, found := receiptsPoints[id]
	if !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "receipt not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"total": val})
}

func (r Receipt) CalulatePoints() (int, error) {
	points := 0

	points += calulateRetailerNamePoints(r.Retailer)
	totalDollarPoints, err := calculateTotalDollarPoints(r.Total)
	if err != nil {
		return 0, err
	}
	points += totalDollarPoints

	itemsListPoints, err := calculateItemsListPoints(r.Items)
	if err != nil {
		return 0, err
	}
	points += itemsListPoints

	purchaseDatePoints, err := calculatePurchaseDatePoints(r.PurchaseDate)
	if err != nil {
		return 0, err
	}
	points += purchaseDatePoints

	purchaseTimePoints, err := calculatePurchaseTimePoints(r.PurchaseTime)
	if err != nil {
		return 0, err
	}
	points += purchaseTimePoints

	return points, nil
}

func calculatePurchaseTimePoints(timeString string) (int, error) {
	tm, err := time.Parse("15:04", timeString)
	if err != nil {
		return 0, err
	}
	// This is assuming that 2:00pm is not after 2:00pm
	if (tm.Hour() == 14 && tm.Minute() != 0) || (tm.Hour() > 14 && tm.Hour() < 16) {
		return 10, nil
	}

	return 0, nil
}

func calculatePurchaseDatePoints(date string) (int, error) {
	day := strings.Split(date, "-")[2]
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		return 0, err
	}
	if dayInt%2 == 1 {
		return 6, nil
	}
	return 0, nil
}

func calulateRetailerNamePoints(name string) int {
	count := 0
	rx := regexp.MustCompile(`[a-zA-Z0-9]`)

	for _, char := range name {
		if rx.MatchString(string(char)) {
			count++
		}
	}

	return count
}

func calculateTotalDollarPoints(total string) (int, error) {
	points := 0
	cents, err := strconv.ParseInt(strings.Split(total, ".")[1], 10, 64)
	if err != nil {
		return 0, err
	}
	if cents == 0 {
		points += 50
	}
	if cents%25 == 0 {
		points += 25
	}
	return points, nil
}

func calculateItemsListPoints(items []Item) (int, error) {

	points := (len(items) / 2) * 5
	for _, item := range items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, err
			}
			println(price)
			points += int(math.Round(price*0.2 + .5))
		}
	}
	return points, nil
}
