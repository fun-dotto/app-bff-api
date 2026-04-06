package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/app-bff-api/generated"
)

func (h *Handler) RoomChangesV1List(_ context.Context, _ api.RoomChangesV1ListRequestObject) (api.RoomChangesV1ListResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
