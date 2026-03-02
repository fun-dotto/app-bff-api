package external

import (
	"github.com/fun-dotto/app-bff-api/generated/external/faculty_api"
	"github.com/fun-dotto/app-bff-api/internal/domain"
)

// ToDomainFaculty は外部APIのFacultyをDomainのFacultyに変換する
func ToDomainFaculty(f faculty_api.Faculty) domain.Faculty {
	return domain.Faculty{
		ID:    f.Id,
		Name:  f.Name,
		Email: f.Email,
	}
}

// ToDomainFaculties は外部APIのFaculty一覧をDomainのFaculty一覧に変換する
func ToDomainFaculties(faculties []faculty_api.Faculty) []domain.Faculty {
	result := make([]domain.Faculty, len(faculties))
	for i, f := range faculties {
		result[i] = ToDomainFaculty(f)
	}
	return result
}
