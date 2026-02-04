package handler

import (
	"net/http"

	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AnnouncementsV1List(c *gin.Context) {
	announcements, err := h.announcementService.GetAnnouncements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	apiAnnouncements := make([]api.Announcement, len(announcements))
	for i, announcement := range announcements {
		apiAnnouncements[i] = toApiAnnouncement(announcement)
	}

	c.JSON(http.StatusOK, gin.H{"announcements": apiAnnouncements})
}

func (h *Handler) AnnouncementsV1Create(c *gin.Context) {
	var req api.AnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainReq := domain.AnnouncementRequest{
		Title:          req.Title,
		AvailableFrom:  req.AvailableFrom,
		AvailableUntil: req.AvailableUntil,
		URL:            req.Url,
	}

	announcement, err := h.announcementService.CreateAnnouncement(domainReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"announcement": toApiAnnouncement(*announcement)})
}

func (h *Handler) AnnouncementsV1Detail(c *gin.Context, id string) {
	announcement, err := h.announcementService.GetAnnouncement(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"announcement": toApiAnnouncement(*announcement)})
}

func (h *Handler) AnnouncementsV1Update(c *gin.Context, id string) {
	var req api.AnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainReq := domain.AnnouncementRequest{
		Title:          req.Title,
		AvailableFrom:  req.AvailableFrom,
		AvailableUntil: req.AvailableUntil,
		URL:            req.Url,
	}

	announcement, err := h.announcementService.UpdateAnnouncement(id, domainReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"announcement": toApiAnnouncement(*announcement)})
}

func (h *Handler) AnnouncementsV1Delete(c *gin.Context, id string) {
	if err := h.announcementService.DeleteAnnouncement(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
