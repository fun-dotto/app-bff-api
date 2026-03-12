package repository

import (
	"context"
	"fmt"

	"github.com/fun-dotto/app-bff-api/generated/external/academic_api"
	"github.com/fun-dotto/app-bff-api/internal/domain"
	"github.com/fun-dotto/app-bff-api/internal/external"
)

type FacultyRepository struct {
	client *academic_api.ClientWithResponses
}

func NewFacultyRepository(client *academic_api.ClientWithResponses) *FacultyRepository {
	return &FacultyRepository{client: client}
}

// GetFaculties は全教員一覧を取得する
func (r *FacultyRepository) GetFaculties() ([]domain.Faculty, error) {
	response, err := r.client.FacultiesV1ListWithResponse(context.Background(), nil)
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
	if len(ids) == 0 {
		return make(map[string]domain.Faculty), nil
	}

	params := &academic_api.FacultiesV1ListParams{
		Ids: &ids,
	}
	response, err := r.client.FacultiesV1ListWithResponse(context.Background(), params)
	if err != nil {
		return nil, fmt.Errorf("failed to call faculty API: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get faculties by IDs: status %d", response.StatusCode())
	}

	result := make(map[string]domain.Faculty, len(response.JSON200.Faculties))
	for _, f := range response.JSON200.Faculties {
		result[f.Id] = external.ToDomainFaculty(f)
	}

	return result, nil
}
