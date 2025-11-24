package external

import (
	"github.com/fun-dotto/app-bff-api/generated/external/announcement_api"
	"github.com/fun-dotto/app-bff-api/internal/domain"
)

func ToDomainAnnouncement(m announcement_api.Announcement) domain.Announcement {
	return domain.Announcement{
		ID:    m.Id,
		Title: m.Title,
		Date:  m.Date,
		URL:   m.Url,
	}
}
