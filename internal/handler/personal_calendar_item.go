package handler

import (
	"context"
	"fmt"
	"net/http"

	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/internal/domain"
	"github.com/fun-dotto/app-bff-api/internal/middleware"
)

// PersonalCalendarItemsV1List 個人カレンダーアイテム一覧を取得する
func (h *Handler) PersonalCalendarItemsV1List(ctx context.Context, request api.PersonalCalendarItemsV1ListRequestObject) (api.PersonalCalendarItemsV1ListResponseObject, error) {
	if h.academicService == nil {
		return nil, errAcademicServiceNotConfigured
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context: %w", fmt.Errorf("%d", http.StatusUnauthorized))
	}

	items, err := h.academicService.GetPersonalCalendarItems(userID, request.Params.Dates)
	if err != nil {
		return nil, fmt.Errorf("failed to get personal calendar items: %w", err)
	}

	apiItems := make([]api.PersonalCalendarItem, len(items))
	for i, item := range items {
		apiItems[i] = toAPIPersonalCalendarItem(item)
	}

	return api.PersonalCalendarItemsV1List200JSONResponse{
		PersonalCalendarItems: apiItems,
	}, nil
}

func toAPIPersonalCalendarItem(item domain.PersonalCalendarItem) api.PersonalCalendarItem {
	return api.PersonalCalendarItem{
		Date: item.Date,
		Slot: api.DottoFoundationV1TimetableSlot{
			DayOfWeek: api.DottoFoundationV1DayOfWeek(item.Slot.DayOfWeek),
			Period:    api.DottoFoundationV1Period(item.Slot.Period),
		},
		TimetableItem: toApiTimetableItem(item.TimetableItem),
	}
}
