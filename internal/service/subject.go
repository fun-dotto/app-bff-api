package service

import (
	"github.com/fun-dotto/app-bff-api/internal/domain"
)

type SubjectRepository interface {
	GetSubjects(query domain.SubjectQuery) ([]domain.Subject, error)
	GetSubject(id string) (*domain.Subject, error)
}

type SubjectService struct {
	subjectRepository SubjectRepository
}

func NewSubjectService(subjectRepository SubjectRepository) *SubjectService {
	return &SubjectService{subjectRepository: subjectRepository}
}

func (s *SubjectService) GetSubjects(query domain.SubjectQuery) ([]domain.Subject, error) {
	return s.subjectRepository.GetSubjects(query)
}

func (s *SubjectService) GetSubject(id string) (*domain.Subject, error) {
	return s.subjectRepository.GetSubject(id)
}
