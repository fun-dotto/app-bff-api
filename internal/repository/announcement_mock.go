package repository

import (
	"fmt"
	"time"

	"github.com/fun-dotto/app-bff-api/internal/domain"
)

// MockAnnouncementRepository はテスト用のモック
type MockAnnouncementRepository struct {
	announcements []domain.Announcement
}

func NewMockAnnouncementRepository() *MockAnnouncementRepository {
	return &MockAnnouncementRepository{
		announcements: []domain.Announcement{
			{
				ID:             "1",
				Title:          "Announcement 1",
				AvailableFrom:  time.Now(),
				AvailableUntil: nil,
				URL:            "https://example.com",
			},
		},
	}
}

func (m *MockAnnouncementRepository) GetAnnouncements(query domain.AnnouncementQuery) ([]domain.Announcement, error) {
	return m.announcements, nil
}

func (m *MockAnnouncementRepository) GetAnnouncement(id string) (*domain.Announcement, error) {
	for _, a := range m.announcements {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, fmt.Errorf("announcement not found: %s", id)
}

func (m *MockAnnouncementRepository) CreateAnnouncement(req domain.AnnouncementRequest) (*domain.Announcement, error) {
	a := domain.Announcement{
		ID:             fmt.Sprintf("%d", len(m.announcements)+1),
		Title:          req.Title,
		AvailableFrom:  req.AvailableFrom,
		AvailableUntil: req.AvailableUntil,
		URL:            req.URL,
	}
	m.announcements = append(m.announcements, a)
	return &a, nil
}

func (m *MockAnnouncementRepository) UpdateAnnouncement(id string, req domain.AnnouncementRequest) (*domain.Announcement, error) {
	for i, a := range m.announcements {
		if a.ID == id {
			m.announcements[i].Title = req.Title
			m.announcements[i].AvailableFrom = req.AvailableFrom
			m.announcements[i].AvailableUntil = req.AvailableUntil
			m.announcements[i].URL = req.URL
			return &m.announcements[i], nil
		}
	}
	return nil, fmt.Errorf("announcement not found: %s", id)
}

func (m *MockAnnouncementRepository) DeleteAnnouncement(id string) error {
	for i, a := range m.announcements {
		if a.ID == id {
			m.announcements = append(m.announcements[:i], m.announcements[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("announcement not found: %s", id)
}
