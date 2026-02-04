package handler

import (
	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/internal/domain"
)

func toApiAnnouncement(announcement domain.Announcement) api.Announcement {
	return api.Announcement{
		Id:    announcement.ID,
		Title: announcement.Title,
		Date:  announcement.AvailableFrom,
		Url:   announcement.URL,
	}
}
