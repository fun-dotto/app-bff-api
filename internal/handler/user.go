package handler

import (
	"context"
	"fmt"
	"net/http"

	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/internal/domain"
	"github.com/fun-dotto/app-bff-api/internal/middleware"
)

// UsersV1Detail ユーザーを取得する
func (h *Handler) UsersV1Detail(ctx context.Context, request api.UsersV1DetailRequestObject) (api.UsersV1DetailResponseObject, error) {
	if h.userService == nil {
		return nil, errUserServiceNotConfigured
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context: %w", fmt.Errorf("%d", http.StatusUnauthorized))
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return api.UsersV1Detail200JSONResponse{
		User: toApiUserInfo(*user),
	}, nil
}

// UsersV1Upsert ユーザーを作成または更新する
func (h *Handler) UsersV1Upsert(ctx context.Context, request api.UsersV1UpsertRequestObject) (api.UsersV1UpsertResponseObject, error) {
	if h.userService == nil {
		return nil, errUserServiceNotConfigured
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context: %w", fmt.Errorf("%d", http.StatusUnauthorized))
	}

	email, _ := middleware.UserEmailFromContext(ctx)

	req := toDomainUserRequest(*request.Body)
	req.Email = email
	user, err := h.userService.UpsertUser(userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}

	return api.UsersV1Upsert200JSONResponse{
		User: toApiUserInfo(*user),
	}, nil
}

func toApiUserInfo(user domain.User) api.UserInfo {
	info := api.UserInfo{}
	if user.Grade != nil {
		g := api.DottoFoundationV1Grade(*user.Grade)
		info.Grade = &g
	}
	if user.Course != nil {
		c := api.DottoFoundationV1Course(*user.Course)
		info.Course = &c
	}
	if user.Class != nil {
		cl := api.DottoFoundationV1Class(*user.Class)
		info.Class = &cl
	}
	return info
}

func toDomainUserRequest(body api.UsersV1UpsertJSONRequestBody) domain.UserRequest {
	req := domain.UserRequest{}
	if body.Grade != nil {
		g := domain.Grade(*body.Grade)
		req.Grade = &g
	}
	if body.Course != nil {
		c := domain.Course(*body.Course)
		req.Course = &c
	}
	if body.Class != nil {
		cl := domain.Class(*body.Class)
		req.Class = &cl
	}
	return req
}
