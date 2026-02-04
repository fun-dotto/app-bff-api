package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/internal/domain"
)

func (h *Handler) AnnouncementsV1List(ctx context.Context, request api.AnnouncementsV1ListRequestObject) (api.AnnouncementsV1ListResponseObject, error) {
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

func (h *Handler) AnnouncementsV1Create(ctx context.Context, request api.AnnouncementsV1CreateRequestObject) (api.AnnouncementsV1CreateResponseObject, error) {
	domainReq := domain.AnnouncementRequest{
		Title:          request.Body.Title,
		AvailableFrom:  request.Body.AvailableFrom,
		AvailableUntil: request.Body.AvailableUntil,
		URL:            request.Body.Url,
	}

	announcement, err := h.announcementService.CreateAnnouncement(domainReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create announcement: %w", err)
	}

	return api.AnnouncementsV1Create201JSONResponse{
		Announcement: toApiAnnouncement(*announcement),
	}, nil
}

func (h *Handler) AnnouncementsV1Detail(ctx context.Context, request api.AnnouncementsV1DetailRequestObject) (api.AnnouncementsV1DetailResponseObject, error) {
	announcement, err := h.announcementService.GetAnnouncement(request.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get announcement: %w", err)
	}

	return api.AnnouncementsV1Detail200JSONResponse{
		Announcement: toApiAnnouncement(*announcement),
	}, nil
}

func (h *Handler) AnnouncementsV1Update(ctx context.Context, request api.AnnouncementsV1UpdateRequestObject) (api.AnnouncementsV1UpdateResponseObject, error) {
	domainReq := domain.AnnouncementRequest{
		Title:          request.Body.Title,
		AvailableFrom:  request.Body.AvailableFrom,
		AvailableUntil: request.Body.AvailableUntil,
		URL:            request.Body.Url,
	}

	announcement, err := h.announcementService.UpdateAnnouncement(request.Id, domainReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update announcement: %w", err)
	}

	return api.AnnouncementsV1Update200JSONResponse{
		Announcement: toApiAnnouncement(*announcement),
	}, nil
}

func (h *Handler) AnnouncementsV1Delete(ctx context.Context, request api.AnnouncementsV1DeleteRequestObject) (api.AnnouncementsV1DeleteResponseObject, error) {
	if err := h.announcementService.DeleteAnnouncement(request.Id); err != nil {
		return nil, fmt.Errorf("failed to delete announcement: %w", err)
	}

	return api.AnnouncementsV1Delete204Response{}, nil
}
