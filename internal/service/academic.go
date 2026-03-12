package service

import (
	"fmt"

	"github.com/fun-dotto/app-bff-api/internal/domain"
)

type AcademicRepository interface {
	GetFaculties() ([]domain.Faculty, error)
	GetFaculty(id string) (*domain.Faculty, error)
	GetFacultiesByIDs(ids []string) (map[string]domain.Faculty, error)
	GetSubjects(query domain.SubjectQuery) ([]domain.Subject, error)
	GetSubject(id string) (*domain.Subject, error)
}

type AcademicService struct {
	repository AcademicRepository
}

func NewAcademicService(repository AcademicRepository) *AcademicService {
	return &AcademicService{repository: repository}
}

func (s *AcademicService) GetFaculties() ([]domain.Faculty, error) {
	return s.repository.GetFaculties()
}

func (s *AcademicService) GetFaculty(id string) (*domain.Faculty, error) {
	return s.repository.GetFaculty(id)
}

func (s *AcademicService) GetFacultiesByIDs(ids []string) (map[string]domain.Faculty, error) {
	return s.repository.GetFacultiesByIDs(ids)
}

func (s *AcademicService) GetSubjects(query domain.SubjectQuery) ([]domain.Subject, error) {
	subjects, err := s.repository.GetSubjects(query)
	if err != nil {
		return nil, err
	}

	if err := s.enrichWithFaculties(subjects); err != nil {
		return nil, fmt.Errorf("failed to enrich with faculties: %w", err)
	}

	return subjects, nil
}

func (s *AcademicService) GetSubject(id string) (*domain.Subject, error) {
	subject, err := s.repository.GetSubject(id)
	if err != nil {
		return nil, err
	}

	if err := s.enrichSubjectWithFaculties(subject); err != nil {
		return nil, fmt.Errorf("failed to enrich with faculties: %w", err)
	}

	return subject, nil
}

// enrichSubjectWithFaculties は単一の科目にFaculty情報を補完する
func (s *AcademicService) enrichSubjectWithFaculties(subject *domain.Subject) error {
	if subject == nil {
		return nil
	}

	// collectFacultyIDs は読み取り専用なので、値コピーで包んでも問題ない
	facultyIDs := collectFacultyIDs([]domain.Subject{*subject})
	if len(facultyIDs) == 0 {
		return nil
	}

	facultyMap, err := s.repository.GetFacultiesByIDs(facultyIDs)
	if err != nil {
		return err
	}

	for i := range subject.Faculties {
		if faculty, ok := facultyMap[subject.Faculties[i].Faculty.ID]; ok {
			subject.Faculties[i].Faculty = faculty
		}
	}

	return nil
}

// enrichWithFaculties は科目一覧にFaculty情報を補完する
func (s *AcademicService) enrichWithFaculties(subjects []domain.Subject) error {
	facultyIDs := collectFacultyIDs(subjects)
	if len(facultyIDs) == 0 {
		return nil
	}

	facultyMap, err := s.repository.GetFacultiesByIDs(facultyIDs)
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
