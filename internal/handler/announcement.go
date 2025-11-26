package handler

import (
	"net/http"

	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AnnouncementsList(c *gin.Context) {
	announcements, err := h.announcementService.GetAnnouncements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	apiAnnouncements := make([]api.Announcement, len(announcements))
	for i, announcement := range announcements {
		apiAnnouncements[i] = toApiAnnouncement(announcement)
	}

	c.JSON(http.StatusOK, apiAnnouncements)
}

func toApiAnnouncement(announcement domain.Announcement) api.Announcement {
	return api.Announcement{
		Id:    announcement.ID,
		Title: announcement.Title,
		Date:  announcement.Date,
		Url:   announcement.URL,
	}
}
