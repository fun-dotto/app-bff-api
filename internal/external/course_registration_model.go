package external

import (
	"github.com/fun-dotto/app-bff-api/generated/external/academic_api"
	"github.com/fun-dotto/app-bff-api/internal/domain"
)

// ToDomainCourseRegistration は外部APIのCourseRegistrationをDomainのCourseRegistrationに変換する
func ToDomainCourseRegistration(m academic_api.CourseRegistration) domain.CourseRegistration {
	return domain.CourseRegistration{
		ID:      m.Id,
		Subject: ToDomainSubjectSummary(m.Subject),
	}
}
