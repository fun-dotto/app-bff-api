package external

import (
	"github.com/fun-dotto/app-bff-api/generated/external/announcement_api"
	"github.com/fun-dotto/app-bff-api/internal/domain"
)

func ToDomainAnnouncement(m announcement_api.Announcement) domain.Announcement {
	return domain.Announcement{
		ID:             m.Id,
		Title:          m.Title,
		AvailableFrom:  m.AvailableFrom,
		AvailableUntil: m.AvailableUntil,
		URL:            m.Url,
	}
}

func ToExternalAnnouncementRequest(r domain.AnnouncementRequest) announcement_api.AnnouncementRequest {
	return announcement_api.AnnouncementRequest{
		Title:          r.Title,
		AvailableFrom:  r.AvailableFrom,
		AvailableUntil: r.AvailableUntil,
		Url:            r.URL,
	}
}
