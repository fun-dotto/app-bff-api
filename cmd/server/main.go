package main

import (
	"context"
	"log"

	firebaseAdmin "firebase.google.com/go/v4"
	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/internal/handler"
	"github.com/fun-dotto/app-bff-api/internal/infrastructure"
	"github.com/fun-dotto/app-bff-api/internal/middleware"
	"github.com/fun-dotto/app-bff-api/internal/repository"
	"github.com/fun-dotto/app-bff-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	ctx := context.Background()

	// Firebase App Check の初期化
	app, err := firebaseAdmin.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v\n", err)
	}

	appCheckClient, err := app.AppCheck(ctx)
	if err != nil {
		log.Fatalf("error initializing App Check client: %v\n", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing Auth client: %v\n", err)
	}

	router := gin.Default()

	// App Check ミドルウェアを適用
	router.Use(middleware.AppCheckMiddleware(appCheckClient))

	// Auth ミドルウェアを適用
	router.Use(middleware.AuthMiddleware(authClient))

	// 外部APIクライアントを初期化
	clients, err := infrastructure.NewExternalClients(ctx)
	if err != nil {
		log.Fatalf("Failed to create external clients: %v", err)
	}

	announcementRepository := repository.NewAnnouncementRepository(clients.Announcement)
	announcementService := service.NewAnnouncementService(announcementRepository)

	academicRepository := repository.NewAcademicRepository(clients.Academic)
	academicService := service.NewAcademicService(academicRepository)

	userRepository := repository.NewUserRepository(clients.User)
	userService := service.NewUserService(userRepository)

	h := handler.NewHandler(
		handler.WithAnnouncementService(announcementService),
		handler.WithAcademicService(academicService),
		handler.WithUserService(userService),
	)

	strictHandler := api.NewStrictHandler(h, nil)

	api.RegisterHandlers(router, strictHandler)

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
