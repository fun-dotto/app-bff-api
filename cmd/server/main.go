package main

import (
	"context"
	"log"
	"os"

	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/generated/external/announcement_api"
	"github.com/fun-dotto/app-bff-api/internal/handler"
	"github.com/fun-dotto/app-bff-api/internal/repository"
	"github.com/fun-dotto/app-bff-api/internal/service"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

func main() {
	router := gin.Default()

	// 環境変数から外部APIのURLを取得
	announcementAPIURL := os.Getenv("ANNOUNCEMENT_API_URL")
	if announcementAPIURL == "" {
		log.Fatal("ANNOUNCEMENT_API_URL is required")
	}

	// 認証付きHTTPクライアントを作成
	ctx := context.Background()
	authClient, err := idtoken.NewClient(ctx, announcementAPIURL)
	if err != nil {
		log.Fatal("Failed to create auth client:", err)
	}

	// 生成されたクライアントに認証付きHTTPクライアントを注入
	apiClient, err := announcement_api.NewClientWithResponses(
		announcementAPIURL,
		announcement_api.WithHTTPClient(authClient),
	)
	if err != nil {
		log.Fatal("Failed to create API client:", err)
	}

	announcementRepository := repository.NewAnnouncementRepository(apiClient)
	announcementService := service.NewAnnouncementService(announcementRepository)
	h := handler.NewHandler(announcementService)

	api.RegisterHandlers(router, h)

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
