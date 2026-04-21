package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fun-dotto/app-bff-api/generated/external/academic_api"
	"github.com/fun-dotto/app-bff-api/generated/external/announcement_api"
	"github.com/fun-dotto/app-bff-api/generated/external/funch_api"
	"github.com/fun-dotto/app-bff-api/generated/external/user_api"
	"google.golang.org/api/idtoken"
)

const httpClientTimeout = 30 * time.Second

// ExternalClients 外部APIクライアントをまとめて管理
type ExternalClients struct {
	Announcement *announcement_api.ClientWithResponses
	Academic     *academic_api.ClientWithResponses
	User         *user_api.ClientWithResponses
	Funch        *funch_api.ClientWithResponses
}

// NewExternalClients 全ての外部APIクライアントを初期化
func NewExternalClients(ctx context.Context) (*ExternalClients, error) {
	announcement, err := newAnnouncementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("announcement client: %w", err)
	}

	academic, err := newAcademicClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("academic client: %w", err)
	}

	user, err := newUserClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("user client: %w", err)
	}

	funch, err := newFunchClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("funch client: %w", err)
	}

	return &ExternalClients{
		Announcement: announcement,
		Academic:     academic,
		User:         user,
		Funch:        funch,
	}, nil
}

func newAnnouncementClient(ctx context.Context) (*announcement_api.ClientWithResponses, error) {
	url := os.Getenv("ANNOUNCEMENT_API_URL")
	if url == "" {
		return nil, fmt.Errorf("ANNOUNCEMENT_API_URL is required")
	}

	authClient, err := newAuthHTTPClient(ctx, url)
	if err != nil {
		return nil, err
	}

	return announcement_api.NewClientWithResponses(
		url,
		announcement_api.WithHTTPClient(authClient),
	)
}

func newAcademicClient(ctx context.Context) (*academic_api.ClientWithResponses, error) {
	url := os.Getenv("ACADEMIC_API_URL")
	if url == "" {
		return nil, fmt.Errorf("ACADEMIC_API_URL is required")
	}

	authClient, err := newAuthHTTPClient(ctx, url)
	if err != nil {
		return nil, err
	}

	return academic_api.NewClientWithResponses(
		url,
		academic_api.WithHTTPClient(authClient),
	)
}

func newUserClient(ctx context.Context) (*user_api.ClientWithResponses, error) {
	url := os.Getenv("USER_API_URL")
	if url == "" {
		return nil, fmt.Errorf("USER_API_URL is required")
	}

	authClient, err := newAuthHTTPClient(ctx, url)
	if err != nil {
		return nil, err
	}

	return user_api.NewClientWithResponses(
		url,
		user_api.WithHTTPClient(authClient),
	)
}

func newFunchClient(ctx context.Context) (*funch_api.ClientWithResponses, error) {
	url := os.Getenv("FUNCH_API_URL")
	if url == "" {
		return nil, fmt.Errorf("FUNCH_API_URL is required")
	}

	authClient, err := newAuthHTTPClient(ctx, url)
	if err != nil {
		return nil, err
	}

	return funch_api.NewClientWithResponses(
		url,
		funch_api.WithHTTPClient(authClient),
	)
}

// newAuthHTTPClient Google Cloud認証付きHTTPクライアントを作成
func newAuthHTTPClient(ctx context.Context, targetURL string) (*http.Client, error) {
	client, err := idtoken.NewClient(ctx, targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth client: %w", err)
	}
	client.Timeout = httpClientTimeout
	return client, nil
}
