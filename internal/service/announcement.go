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
	sortByDate := "desc"
	filterIsActive := true
	query := domain.AnnouncementQuery{
		SortByDate:     &sortByDate,
		FilterIsActive: &filterIsActive,
	}
	return s.announcementRepository.GetAnnouncements(query)
}
