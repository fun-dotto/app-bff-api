package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/app-bff-api/generated"
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
