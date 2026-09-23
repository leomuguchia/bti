package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func respondJSON(c *gin.Context, status int, payload any) {
	c.JSON(status, payload)
}

func respondError(c *gin.Context, status int, msg string) {
	fmt.Println("respondError:", status, msg)
	c.JSON(status, gin.H{"error": msg})
}
