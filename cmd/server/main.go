package main

import (
	"log"

	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	h := handler.NewHandler()
	api.RegisterHandlers(router, h)

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
