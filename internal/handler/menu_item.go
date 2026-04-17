package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/app-bff-api/generated"
)

func (h *Handler) MenuItemsV1List(_ context.Context, _ api.MenuItemsV1ListRequestObject) (api.MenuItemsV1ListResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
