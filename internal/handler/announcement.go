package handler

import (
	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AnnouncementsList(c *gin.Context) {
	domainAnnouncements, err := h.announcementService.GetAnnouncements()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// ドメインモデル → API型に変換
	apiAnnouncements := make([]api.Announcement, len(domainAnnouncements))
	for i, a := range domainAnnouncements {
		apiAnnouncements[i] = api.Announcement{
			Id:    a.ID,
			Title: a.Title,
			Date:  a.Date,
			Url:   a.URL,
		}
	}

	c.JSON(200, apiAnnouncements)
}
