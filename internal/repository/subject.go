package repository

import (
	"context"
	"fmt"

	"github.com/fun-dotto/app-bff-api/generated/external/academic_api"
	"github.com/fun-dotto/app-bff-api/internal/domain"
	"github.com/fun-dotto/app-bff-api/internal/external"
)

type SubjectRepository struct {
	client *academic_api.ClientWithResponses
}

func NewSubjectRepository(client *academic_api.ClientWithResponses) *SubjectRepository {
	return &SubjectRepository{client: client}
}

// GetSubjects は外部APIから科目一覧を取得する
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

// GetSubject は外部APIから科目詳細を取得する
func (r *SubjectRepository) GetSubject(id string) (*domain.Subject, error) {
	ctx := context.Background()

	// 科目詳細を取得
	subjectResponse, err := r.client.SubjectsV1DetailWithResponse(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subject: %w", err)
	}
	if subjectResponse.JSON200 == nil {
		return nil, fmt.Errorf("failed to get subject: status %d", subjectResponse.StatusCode())
	}

	// シラバスを取得
	syllabusResponse, err := r.client.SyllabusV1DetailWithResponse(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get syllabus: %w", err)
	}
	if syllabusResponse.JSON200 == nil {
		return nil, fmt.Errorf("failed to get syllabus: status %d", syllabusResponse.StatusCode())
	}

	s := external.ToDomainSubject(subjectResponse.JSON200.Subject)

	syllabus := external.ToDomainSyllabus(syllabusResponse.JSON200.Syllabus)
	s.Syllabus = &syllabus

	return &s, nil
}
