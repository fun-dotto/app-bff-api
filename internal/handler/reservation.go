package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/app-bff-api/generated"
)

func (h *Handler) ReservationsV1List(_ context.Context, _ api.ReservationsV1ListRequestObject) (api.ReservationsV1ListResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
