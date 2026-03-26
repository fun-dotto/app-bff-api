package handler

import (
	"context"
	"fmt"
	"net/http"

	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/internal/domain"
	"github.com/fun-dotto/app-bff-api/internal/middleware"
)

// CourseRegistrationsV1List 履修登録一覧を取得する
func (h *Handler) CourseRegistrationsV1List(ctx context.Context, request api.CourseRegistrationsV1ListRequestObject) (api.CourseRegistrationsV1ListResponseObject, error) {
	if h.academicService == nil {
		return nil, errAcademicServiceNotConfigured
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context: %w", fmt.Errorf("%d", http.StatusUnauthorized))
	}

	semesters := make([]domain.CourseSemester, len(request.Params.Semesters))
	for i, semester := range request.Params.Semesters {
		semesters[i] = domain.CourseSemester(semester)
	}

	registrations, err := h.academicService.GetCourseRegistrations(userID, semesters, request.Params.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to get course registrations: %w", err)
	}

	apiRegistrations := make([]api.CourseRegistration, len(registrations))
	for i, r := range registrations {
		apiRegistrations[i] = toApiCourseRegistration(r)
	}

	return api.CourseRegistrationsV1List200JSONResponse{
		CourseRegistrations: apiRegistrations,
	}, nil
}

// CourseRegistrationsV1Create 履修登録を作成する
func (h *Handler) CourseRegistrationsV1Create(ctx context.Context, request api.CourseRegistrationsV1CreateRequestObject) (api.CourseRegistrationsV1CreateResponseObject, error) {
	if h.academicService == nil {
		return nil, errAcademicServiceNotConfigured
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context: %w", fmt.Errorf("%d", http.StatusUnauthorized))
	}

	registration, err := h.academicService.CreateCourseRegistration(userID, request.Body.SubjectId)
	if err != nil {
		return nil, fmt.Errorf("failed to create course registration: %w", err)
	}

	return api.CourseRegistrationsV1Create201JSONResponse{
		CourseRegistration: toApiCourseRegistration(*registration),
	}, nil
}

// CourseRegistrationsV1Delete 履修登録を削除する
func (h *Handler) CourseRegistrationsV1Delete(ctx context.Context, request api.CourseRegistrationsV1DeleteRequestObject) (api.CourseRegistrationsV1DeleteResponseObject, error) {
	if h.academicService == nil {
		return nil, errAcademicServiceNotConfigured
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context: %w", fmt.Errorf("%d", http.StatusUnauthorized))
	}
	_ = userID

	err := h.academicService.DeleteCourseRegistration(request.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete course registration: %w", err)
	}

	return api.CourseRegistrationsV1Delete204Response{}, nil
}

func toApiCourseRegistration(r domain.CourseRegistration) api.CourseRegistration {
	return api.CourseRegistration{
		Id:      r.ID,
		Subject: toApiSubjectSummary(r.Subject),
	}
}
