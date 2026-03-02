package service

import (
	"github.com/fun-dotto/app-bff-api/internal/domain"
)

type FacultyRepository interface {
	GetFaculties() ([]domain.Faculty, error)
	GetFaculty(id string) (*domain.Faculty, error)
	GetFacultiesByIDs(ids []string) (map[string]domain.Faculty, error)
}

type FacultyService struct {
	facultyRepository FacultyRepository
}

func NewFacultyService(facultyRepository FacultyRepository) *FacultyService {
	return &FacultyService{facultyRepository: facultyRepository}
}

func (s *FacultyService) GetFaculties() ([]domain.Faculty, error) {
	return s.facultyRepository.GetFaculties()
}

func (s *FacultyService) GetFaculty(id string) (*domain.Faculty, error) {
	return s.facultyRepository.GetFaculty(id)
}

func (s *FacultyService) GetFacultiesByIDs(ids []string) (map[string]domain.Faculty, error) {
	return s.facultyRepository.GetFacultiesByIDs(ids)
}
