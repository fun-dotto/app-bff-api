package handler

import (
	"time"

	"github.com/fun-dotto/app-bff-api/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AnnouncementsList(c *gin.Context) {
	announcements := h.getAnnouncements()
	c.JSON(200, announcements)
}

func (h *Handler) getAnnouncements() []domain.Announcement {
	return h.newMockAnnouncement()
}

func (h *Handler) newMockAnnouncement() []domain.Announcement {
	return []domain.Announcement{
		{
			ID:       "1",
			Title:    "Announcement 1",
			Date:     time.Now(),
			URL:      "https://example.com",
			IsActive: true,
		},
	}
}
