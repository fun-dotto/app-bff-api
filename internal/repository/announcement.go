package repository

import (
	"context"
	"fmt"

	"github.com/fun-dotto/app-bff-api/generated/external/announcement_api"
	"github.com/fun-dotto/app-bff-api/internal/domain"
	"github.com/fun-dotto/app-bff-api/internal/external"
)

// AnnouncementRepository は外部APIからお知らせを取得する
type AnnouncementRepository struct {
	client *announcement_api.ClientWithResponses
}

func NewAnnouncementRepository(client *announcement_api.ClientWithResponses) *AnnouncementRepository {
	return &AnnouncementRepository{client: client}
}

func (r *AnnouncementRepository) GetAnnouncements(query domain.AnnouncementQuery) ([]domain.Announcement, error) {
	params := announcement_api.AnnouncementsListParams{
		SortByDate: func() *announcement_api.SortDirection {
			if query.SortByDate == nil {
				return nil
			}
			sortDirection := announcement_api.SortDirection(*query.SortByDate)
			return &sortDirection
		}(),
		FilterIsActive: query.FilterIsActive,
	}
	// 外部APIからデータ取得
	response, err := r.client.AnnouncementsListWithResponse(context.Background(), &params)
	if err != nil {
		return nil, err
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get announcements: status %d", response.StatusCode())
	}

	// 外部API形式 → ドメイン形式に変換
	result := make([]domain.Announcement, len(*response.JSON200))
	for i, a := range *response.JSON200 {
		result[i] = external.ToDomainAnnouncement(a)
	}

	return result, nil
}
