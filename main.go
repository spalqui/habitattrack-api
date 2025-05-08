package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	err := r.Run()
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
