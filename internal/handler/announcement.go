package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) AnnouncementsList(c *gin.Context) {
	announcements, err := h.announcementService.GetAnnouncements()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, announcements)
}
