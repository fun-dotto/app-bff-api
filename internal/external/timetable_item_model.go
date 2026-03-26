package external

import (
	"time"

	"github.com/fun-dotto/app-bff-api/generated/external/academic_api"
	"github.com/fun-dotto/app-bff-api/internal/domain"
)

// ToDomainTimetableItem は外部APIのTimetableItemをDomainのTimetableItemに変換する
func ToDomainTimetableItem(m academic_api.TimetableItem) domain.TimetableItem {
	var slot *domain.TimetableSlot
	if m.Slot != nil {
		slot = &domain.TimetableSlot{
			DayOfWeek: domain.DayOfWeek(m.Slot.DayOfWeek),
			Period:    domain.Period(m.Slot.Period),
		}
	}

	return domain.TimetableItem{
		ID:      m.Id,
		Slot:    slot,
		Rooms:   toDomainRooms(m.Rooms),
		Subject: ToDomainSubjectSummary(m.Subject),
	}
}

// ToExternalTimetableItemQuery はDomainのTimetableItemQueryを外部APIのTimetableItemsV1ListParamsに変換する
func ToExternalTimetableItemQuery(q domain.TimetableItemQuery) *academic_api.TimetableItemsV1ListParams {
	params := &academic_api.TimetableItemsV1ListParams{
		Semesters: toExternalSemesters(q.Semesters),
		Year:      q.Year,
	}

	return params
}

// ToDomainPersonalCalendarItem は外部APIのPersonalCalendarItemをDomainに変換する
func ToDomainPersonalCalendarItem(m academic_api.PersonalCalendarItem) domain.PersonalCalendarItem {
	return domain.PersonalCalendarItem{
		Date: m.Date,
		Slot: domain.TimetableSlot{
			DayOfWeek: domain.DayOfWeek(m.Slot.DayOfWeek),
			Period:    domain.Period(m.Slot.Period),
		},
		TimetableItem: ToDomainTimetableItem(m.TimetableItem),
	}
}

// ToExternalPersonalCalendarItemParams は domain の検索条件を外部APIのクエリに変換する
func ToExternalPersonalCalendarItemParams(userID string, dates []time.Time) *academic_api.PersonalCalendarItemsV1ListParams {
	return &academic_api.PersonalCalendarItemsV1ListParams{
		UserId: userID,
		Dates:  dates,
	}
}

func toDomainRooms(rooms []academic_api.Room) []domain.Room {
	result := make([]domain.Room, len(rooms))
	for i, room := range rooms {
		result[i] = domain.Room{
			ID:    room.Id,
			Name:  room.Name,
			Floor: domain.Floor(room.Floor),
		}
	}
	return result
}
