package repository

import (
	"context"
	"fmt"

	"github.com/fun-dotto/app-bff-api/generated/external/subject_api"
	"github.com/fun-dotto/app-bff-api/internal/domain"
	"github.com/fun-dotto/app-bff-api/internal/external"
)

type SubjectRepository struct {
	client *subject_api.ClientWithResponses
}

func NewSubjectRepository(client *subject_api.ClientWithResponses) *SubjectRepository {
	return &SubjectRepository{client: client}
}

func (r *SubjectRepository) GetSubjects(query domain.SubjectQuery) ([]domain.Subject, error) {
	params := external.ToExternalSubjectQuery(query)

	response, err := r.client.SubjectsV1ListWithResponse(context.Background(), params)
	if err != nil {
		return nil, err
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get subjects: status %d", response.StatusCode())
	}

	subjects := response.JSON200.Subjects
	result := make([]domain.Subject, len(subjects))
	for i, s := range subjects {
		result[i] = external.ToDomainSubjectSummary(s)
	}

	return result, nil
}

func (r *SubjectRepository) GetSubject(id string) (*domain.Subject, error) {
	response, err := r.client.SubjectsV1DetailWithResponse(context.Background(), id)
	if err != nil {
		return nil, err
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get subject: status %d", response.StatusCode())
	}

	s := external.ToDomainSubject(response.JSON200.Subject)
	return &s, nil
}
