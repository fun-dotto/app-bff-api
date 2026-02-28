package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/internal/domain"
)

func (h *Handler) AnnouncementsV0List(ctx context.Context, request api.AnnouncementsV0ListRequestObject) (api.AnnouncementsV0ListResponseObject, error) {
	announcements, err := h.announcementService.GetAnnouncements()
	if err != nil {
		return nil, fmt.Errorf("failed to get announcements: %w", err)
	}

	apiAnnouncements := make(api.AnnouncementsV0List200JSONResponse, len(announcements))
	for i, announcement := range announcements {
		apiAnnouncements[i] = toApiAnnouncement(announcement)
	}

	return apiAnnouncements, nil
}

// toApiAnnouncement はDomainのお知らせをAPIのお知らせに変換する
func toApiAnnouncement(announcement domain.Announcement) api.Announcement {
	return api.Announcement{
		Id:    announcement.ID,
		Title: announcement.Title,
		Date:  announcement.AvailableFrom,
		Url:   announcement.URL,
	}
}
