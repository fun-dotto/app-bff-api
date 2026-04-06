package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/app-bff-api/generated"
)

func (h *Handler) MakeupClassesV1List(_ context.Context, _ api.MakeupClassesV1ListRequestObject) (api.MakeupClassesV1ListResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
