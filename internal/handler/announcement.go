package handler

import (
	"net/http"

	api "github.com/fun-dotto/app-bff-api/generated"
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
