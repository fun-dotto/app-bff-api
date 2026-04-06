package external

import (
	"github.com/fun-dotto/app-bff-api/generated/external/academic_api"
	"github.com/fun-dotto/app-bff-api/internal/domain"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ToDomainMakeupClass は外部APIのMakeupClassをDomainに変換する
func ToDomainMakeupClass(m academic_api.MakeupClass) domain.MakeupClass {
	return domain.MakeupClass{
		ID:      m.Id,
		Comment: m.Comment,
		Date:    m.Date.Time,
		Period:  domain.Period(m.Period),
		Subject: ToDomainSubjectSummary(m.Subject),
	}
}

// ToExternalMakeupClassQuery はDomainのMakeupClassQueryを外部APIのパラメータに変換する
func ToExternalMakeupClassQuery(q domain.MakeupClassQuery) *academic_api.MakeupClassesV1ListParams {
	params := &academic_api.MakeupClassesV1ListParams{
		SubjectIds: q.SubjectIds,
	}
	if q.From != nil {
		d := openapi_types.Date{Time: *q.From}
		params.From = &d
	}
	if q.Until != nil {
		d := openapi_types.Date{Time: *q.Until}
		params.Until = &d
	}
	return params
}
