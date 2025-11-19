package handler

import (
	"time"

	"github.com/fun-dotto/app-bff-api/internal/domain"
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) AnnouncementsList(c *gin.Context) {
	announcements := []domain.Announcement{
		{
			ID:       "1",
			Title:    "Announcement 1",
			Date:     time.Now(),
			URL:      "https://example.com",
			IsActive: true,
		},
	}

	c.JSON(200, announcements)
}
