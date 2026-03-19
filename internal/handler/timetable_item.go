package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/internal/domain"
)

// TimetableItemsV1List 時間割アイテム一覧を取得する
func (h *Handler) TimetableItemsV1List(ctx context.Context, request api.TimetableItemsV1ListRequestObject) (api.TimetableItemsV1ListResponseObject, error) {
	if h.academicService == nil {
		return nil, errAcademicServiceNotConfigured
	}

	query := toTimetableItemQuery(request.Params)

	items, err := h.academicService.GetTimetableItems(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get timetable items: %w", err)
	}

	apiItems := make([]api.TimetableItem, len(items))
	for i, item := range items {
		apiItems[i] = toApiTimetableItem(item)
	}

	return api.TimetableItemsV1List200JSONResponse{
		TimetableItems: apiItems,
	}, nil
}

// toTimetableItemQuery は TimetableItemsV1List の API パラメータを domain.TimetableItemQuery に変換する
func toTimetableItemQuery(params api.TimetableItemsV1ListParams) domain.TimetableItemQuery {
	query := domain.TimetableItemQuery{
		Semester: domain.CourseSemester(params.Semester),
		Year:     params.Year,
	}

	if params.DayOfWeek != nil {
		query.DayOfWeek = make([]domain.DayOfWeek, len(*params.DayOfWeek))
		for i, d := range *params.DayOfWeek {
			query.DayOfWeek[i] = domain.DayOfWeek(d)
		}
	}

	return query
}

// toApiTimetableItem はDomainの時間割アイテムをAPIの時間割アイテムに変換する
func toApiTimetableItem(item domain.TimetableItem) api.TimetableItem {
	return api.TimetableItem{
		Id:        item.ID,
		DayOfWeek: api.DottoFoundationV1DayOfWeek(item.DayOfWeek),
		Period:    api.DottoFoundationV1Period(item.Period),
		Subject:   toApiSubjectSummary(item.Subject),
	}
}
