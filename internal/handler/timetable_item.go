package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/app-bff-api/generated"
)

// TimetableItemsV1List 時間割アイテム一覧を取得する
// TODO: 時間割APIを作成したら実装する
func (h *Handler) TimetableItemsV1List(_ context.Context, _ api.TimetableItemsV1ListRequestObject) (api.TimetableItemsV1ListResponseObject, error) {
	return nil, fmt.Errorf("TimetableItemsV1List is not implemented yet")
}
