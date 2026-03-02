package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fun-dotto/app-bff-api/generated/external/announcement_api"
	"github.com/fun-dotto/app-bff-api/generated/external/faculty_api"
	"github.com/fun-dotto/app-bff-api/generated/external/subject_api"
	"google.golang.org/api/idtoken"
)

const httpClientTimeout = 30 * time.Second

// ExternalClients 外部APIクライアントをまとめて管理
type ExternalClients struct {
	Announcement *announcement_api.ClientWithResponses
	Subject      *subject_api.ClientWithResponses
	Faculty      *faculty_api.ClientWithResponses
}

// NewExternalClients 全ての外部APIクライアントを初期化
func NewExternalClients(ctx context.Context) (*ExternalClients, error) {
	announcement, err := newAnnouncementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("announcement client: %w", err)
	}

	subject, err := newSubjectClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("subject client: %w", err)
	}

	faculty, err := newFacultyClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("faculty client: %w", err)
	}

	return &ExternalClients{
		Announcement: announcement,
		Subject:      subject,
		Faculty:      faculty,
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

func newSubjectClient(ctx context.Context) (*subject_api.ClientWithResponses, error) {
	url := os.Getenv("SUBJECT_API_URL")
	if url == "" {
		return nil, fmt.Errorf("SUBJECT_API_URL is required")
	}

	authClient, err := newAuthHTTPClient(ctx, url)
	if err != nil {
		return nil, err
	}

	return subject_api.NewClientWithResponses(
		url,
		subject_api.WithHTTPClient(authClient),
	)
}

func newFacultyClient(ctx context.Context) (*faculty_api.ClientWithResponses, error) {
	url := os.Getenv("FACULTY_API_URL")
	if url == "" {
		return nil, fmt.Errorf("FACULTY_API_URL is required")
	}

	authClient, err := newAuthHTTPClient(ctx, url)
	if err != nil {
		return nil, err
	}

	return faculty_api.NewClientWithResponses(
		url,
		faculty_api.WithHTTPClient(authClient),
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
