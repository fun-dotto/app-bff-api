package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/app-bff-api/generated"
)

func (h *Handler) CancelledClassesV1List(_ context.Context, _ api.CancelledClassesV1ListRequestObject) (api.CancelledClassesV1ListResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
