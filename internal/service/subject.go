package service

import (
	"fmt"

	"github.com/fun-dotto/app-bff-api/internal/domain"
)

type SubjectRepository interface {
	GetSubjects(query domain.SubjectQuery) ([]domain.Subject, error)
	GetSubject(id string) (*domain.Subject, error)
}

type FacultyGetter interface {
	GetFacultiesByIDs(ids []string) (map[string]domain.Faculty, error)
}

type SubjectService struct {
	subjectRepository SubjectRepository
	facultyService    FacultyGetter
}

func NewSubjectService(subjectRepository SubjectRepository, facultyService FacultyGetter) *SubjectService {
	return &SubjectService{
		subjectRepository: subjectRepository,
		facultyService:    facultyService,
	}
}

func (s *SubjectService) GetSubjects(query domain.SubjectQuery) ([]domain.Subject, error) {
	subjects, err := s.subjectRepository.GetSubjects(query)
	if err != nil {
		return nil, err
	}

	if err := s.enrichWithFaculties(subjects); err != nil {
		return nil, fmt.Errorf("failed to enrich with faculties: %w", err)
	}

	return subjects, nil
}

func (s *SubjectService) GetSubject(id string) (*domain.Subject, error) {
	subject, err := s.subjectRepository.GetSubject(id)
	if err != nil {
		return nil, err
	}

	if err := s.enrichWithFaculties([]domain.Subject{*subject}); err != nil {
		return nil, fmt.Errorf("failed to enrich with faculties: %w", err)
	}

	return subject, nil
}

// enrichWithFaculties は科目一覧にFaculty情報を補完する
func (s *SubjectService) enrichWithFaculties(subjects []domain.Subject) error {
	facultyIDs := collectFacultyIDs(subjects)
	if len(facultyIDs) == 0 {
		return nil
	}

	facultyMap, err := s.facultyService.GetFacultiesByIDs(facultyIDs)
	if err != nil {
		return err
	}

	for i := range subjects {
		for j := range subjects[i].Faculties {
			if faculty, ok := facultyMap[subjects[i].Faculties[j].Faculty.ID]; ok {
				subjects[i].Faculties[j].Faculty = faculty
			}
		}
	}

	return nil
}

// collectFacultyIDs は科目一覧からFaculty IDを収集する
func collectFacultyIDs(subjects []domain.Subject) []string {
	idSet := make(map[string]struct{})
	for _, s := range subjects {
		for _, f := range s.Faculties {
			if f.Faculty.ID != "" {
				idSet[f.Faculty.ID] = struct{}{}
			}
		}
	}

	ids := make([]string, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	return ids
}
