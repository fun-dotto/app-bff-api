package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/app-bff-api/generated"
)

var errNotImplemented = errors.New("not implemented")

// UsersV1Detail ユーザーを取得する
func (h *Handler) UsersV1Detail(ctx context.Context, request api.UsersV1DetailRequestObject) (api.UsersV1DetailResponseObject, error) {
	return nil, errNotImplemented
}

// UsersV1Upsert ユーザーを作成または更新する
func (h *Handler) UsersV1Upsert(ctx context.Context, request api.UsersV1UpsertRequestObject) (api.UsersV1UpsertResponseObject, error) {
	return nil, errNotImplemented
}
