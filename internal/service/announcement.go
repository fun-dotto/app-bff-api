package service

import "github.com/fun-dotto/app-bff-api/internal/domain"

type AnouncementRepository interface {
	GetAnnouncements(query domain.AnnouncementQuery) ([]domain.Announcement, error)
}

type AnnouncementService struct {
	announcementRepository AnouncementRepository
}

func NewAnnouncementService(announcementRepository AnouncementRepository) *AnnouncementService {
	return &AnnouncementService{announcementRepository: announcementRepository}
}

func (s *AnnouncementService) GetAnnouncements() ([]domain.Announcement, error) {
	sortByDateAsc := false
	filterIsActive := true
	query := domain.AnnouncementQuery{
		SortByDateAsc:  &sortByDateAsc,
		FilterIsActive: &filterIsActive,
	}
	return s.announcementRepository.GetAnnouncements(query)
}
