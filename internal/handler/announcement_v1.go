package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/app-bff-api/generated"
)

func (h *Handler) AnnouncementsV1List(ctx context.Context, request api.AnnouncementsV1ListRequestObject) (api.AnnouncementsV1ListResponseObject, error) {
	if h.announcementService == nil {
		return nil, errAnnouncementServiceNotConfigured
	}

	announcements, err := h.announcementService.GetAnnouncements()
	if err != nil {
		return nil, fmt.Errorf("failed to get announcements: %w", err)
	}

	apiAnnouncements := make([]api.Announcement, len(announcements))
	for i, announcement := range announcements {
		apiAnnouncements[i] = toApiAnnouncement(announcement)
	}

	return api.AnnouncementsV1List200JSONResponse{
		Announcements: apiAnnouncements,
	}, nil
}
