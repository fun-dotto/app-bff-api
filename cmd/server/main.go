package main

import (
	"context"
	"log"
	"os"

	firebaseAdmin "firebase.google.com/go/v4"
	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/generated/external/announcement_api"
	"github.com/fun-dotto/app-bff-api/internal/handler"
	"github.com/fun-dotto/app-bff-api/internal/middleware"
	"github.com/fun-dotto/app-bff-api/internal/repository"
	"github.com/fun-dotto/app-bff-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/idtoken"
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

	router := gin.Default()

	// App Check ミドルウェアを適用
	router.Use(middleware.AppCheckMiddleware(appCheckClient))

	// 環境変数から外部APIのURLを取得
	announcementAPIURL := os.Getenv("ANNOUNCEMENT_API_URL")
	if announcementAPIURL == "" {
		log.Fatal("ANNOUNCEMENT_API_URL is required")
	}

	// 認証付きHTTPクライアントを作成
	announcementAPIAuthClient, err := idtoken.NewClient(ctx, announcementAPIURL)
	if err != nil {
		log.Fatal("Failed to create auth client:", err)
	}

	// 生成されたクライアントに認証付きHTTPクライアントを注入
	announcementAPIClient, err := announcement_api.NewClientWithResponses(
		announcementAPIURL,
		announcement_api.WithHTTPClient(announcementAPIAuthClient),
	)
	if err != nil {
		log.Fatal("Failed to create API client:", err)
	}

	announcementRepository := repository.NewAnnouncementRepository(announcementAPIClient)

	announcementService := service.NewAnnouncementService(announcementRepository)

	h := handler.NewHandler(announcementService)

	strictHandler := api.NewStrictHandler(h, nil)

	api.RegisterHandlers(router, strictHandler)

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
