package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/internal/repository"
	"github.com/fun-dotto/app-bff-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAnnouncementsV1List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		setupContext func(c *gin.Context)
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:         "正常にお知らせ一覧が取得できる",
			setupContext: func(c *gin.Context) {},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var result struct {
					Announcements []api.Announcement `json:"announcements"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.NotEmpty(t, result.Announcements, "アナウンスメントが空です")
			},
		},
		{
			name:         "Content-Typeがapplication/jsonである",
			setupContext: func(c *gin.Context) {},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
			},
		},
		{
			name:         "レスポンスがオブジェクト形式である",
			setupContext: func(c *gin.Context) {},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var result map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				_, hasKey := result["announcements"]
				assert.True(t, hasKey, "レスポンスにannouncementsキーがありません")
			},
		},
		{
			name:         "お知らせのフィールドが正しく返される",
			setupContext: func(c *gin.Context) {},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var result struct {
					Announcements []api.Announcement `json:"announcements"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Len(t, result.Announcements, 1, "MockRepositoryは1件返すはずです")
				assert.Equal(t, "1", result.Announcements[0].Id)
				assert.Equal(t, "Announcement 1", result.Announcements[0].Title)
				assert.Equal(t, "https://example.com", result.Announcements[0].Url)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewMockAnnouncementRepository()
			h := NewHandler(service.NewAnnouncementService(mockRepo))
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.setupContext != nil {
				tt.setupContext(c)
			}

			h.AnnouncementsV1List(c)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}
