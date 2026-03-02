package repository

import (
	"context"
	"fmt"

	"github.com/fun-dotto/app-bff-api/generated/external/faculty_api"
	"github.com/fun-dotto/app-bff-api/internal/domain"
	"github.com/fun-dotto/app-bff-api/internal/external"
)

type FacultyRepository struct {
	client *faculty_api.ClientWithResponses
}

func NewFacultyRepository(client *faculty_api.ClientWithResponses) *FacultyRepository {
	return &FacultyRepository{client: client}
}

// GetFaculties は全教員一覧を取得する
func (r *FacultyRepository) GetFaculties() ([]domain.Faculty, error) {
	response, err := r.client.FacultiesV1ListWithResponse(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to call faculty API: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get faculties: status %d", response.StatusCode())
	}

	return external.ToDomainFaculties(response.JSON200.Faculties), nil
}

// GetFaculty は教員詳細を取得する
func (r *FacultyRepository) GetFaculty(id string) (*domain.Faculty, error) {
	response, err := r.client.FacultiesV1DetailWithResponse(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to call faculty API: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get faculty: status %d", response.StatusCode())
	}

	faculty := external.ToDomainFaculty(response.JSON200.Faculty)
	return &faculty, nil
}

// GetFacultiesByIDs は指定したIDの教員一覧を取得する
func (r *FacultyRepository) GetFacultiesByIDs(ids []string) (map[string]domain.Faculty, error) {
	allFaculties, err := r.GetFaculties()
	if err != nil {
		return nil, err
	}

	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	result := make(map[string]domain.Faculty, len(ids))
	for _, f := range allFaculties {
		if _, ok := idSet[f.ID]; ok {
			result[f.ID] = f
		}
	}

	return result, nil
}
