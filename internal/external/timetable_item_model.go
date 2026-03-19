package external

import (
	"github.com/fun-dotto/app-bff-api/generated/external/academic_api"
	"github.com/fun-dotto/app-bff-api/internal/domain"
)

// ToDomainTimetableItem は外部APIのTimetableItemをDomainのTimetableItemに変換する
func ToDomainTimetableItem(m academic_api.TimetableItem) domain.TimetableItem {
	return domain.TimetableItem{
		ID:        m.Id,
		DayOfWeek: domain.DayOfWeek(m.DayOfWeek),
		Period:    domain.Period(m.Period),
		Subject:   ToDomainSubjectSummary(m.Subject),
	}
}

// ToExternalTimetableItemQuery はDomainのTimetableItemQueryを外部APIのTimetableItemsV1ListParamsに変換する
func ToExternalTimetableItemQuery(q domain.TimetableItemQuery) *academic_api.TimetableItemsV1ListParams {
	params := &academic_api.TimetableItemsV1ListParams{
		Semester: academic_api.DottoFoundationV1CourseSemester(q.Semester),
		Year:     q.Year,
	}

	if len(q.DayOfWeek) > 0 {
		dayOfWeek := make([]academic_api.DottoFoundationV1DayOfWeek, len(q.DayOfWeek))
		for i, d := range q.DayOfWeek {
			dayOfWeek[i] = academic_api.DottoFoundationV1DayOfWeek(d)
		}
		params.DayOfWeek = &dayOfWeek
	}

	return params
}
