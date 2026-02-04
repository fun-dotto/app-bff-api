package handler

import (
	"context"
	"testing"

	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/internal/repository"
	"github.com/fun-dotto/app-bff-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnnouncementsV0List(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T, resp api.AnnouncementsV0ListResponseObject, err error)
	}{
		{
			name: "正常にお知らせ一覧が取得できる",
			validate: func(t *testing.T, resp api.AnnouncementsV0ListResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.AnnouncementsV0List200JSONResponse)
				require.True(t, ok, "レスポンスが200 JSONレスポンスではありません")
				assert.NotEmpty(t, result, "アナウンスメントが空です")
			},
		},
		{
			name: "レスポンスが配列形式である",
			validate: func(t *testing.T, resp api.AnnouncementsV0ListResponseObject, err error) {
				require.NoError(t, err)
				_, ok := resp.(api.AnnouncementsV0List200JSONResponse)
				require.True(t, ok, "レスポンスが配列形式ではありません")
			},
		},
		{
			name: "お知らせのフィールドが正しく返される",
			validate: func(t *testing.T, resp api.AnnouncementsV0ListResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.AnnouncementsV0List200JSONResponse)
				require.True(t, ok)
				assert.Len(t, result, 1, "MockRepositoryは1件返すはずです")
				assert.Equal(t, "1", result[0].Id)
				assert.Equal(t, "Announcement 1", result[0].Title)
				assert.Equal(t, "https://example.com", result[0].Url)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewMockAnnouncementRepository()
			h := NewHandler(service.NewAnnouncementService(mockRepo))

			resp, err := h.AnnouncementsV0List(context.Background(), api.AnnouncementsV0ListRequestObject{})

			if tt.validate != nil {
				tt.validate(t, resp, err)
			}
		})
	}
}
