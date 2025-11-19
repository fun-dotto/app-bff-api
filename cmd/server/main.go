package main

import (
	"log"

	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/internal/handler"
	"github.com/fun-dotto/app-bff-api/internal/repository"
	"github.com/fun-dotto/app-bff-api/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	announcementRepository := repository.NewMockAnnouncementRepository()

	announcementService := service.NewAnnouncementService(announcementRepository)

	h := handler.NewHandler(announcementService)

	api.RegisterHandlers(router, h)

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
