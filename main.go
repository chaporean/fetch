package main

import (
	"fetch/receipts"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/receipts/:id/points", receipts.HandleGetPoints)
	router.POST("/receipts/process", receipts.HandleProcessReceipt)

	router.Run("localhost:8080")
}
