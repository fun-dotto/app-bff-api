package handler

import (
	"errors"

	"github.com/fun-dotto/app-bff-api/internal/service"
)

var (
	errAnnouncementServiceNotConfigured = errors.New("announcementService is not configured")
	errSubjectServiceNotConfigured      = errors.New("subjectService is not configured")
)

type Handler struct {
	announcementService *service.AnnouncementService
	subjectService      *service.SubjectService
}

type Option func(*Handler)

func WithAnnouncementService(s *service.AnnouncementService) Option {
	return func(h *Handler) {
		h.announcementService = s
	}
}

func WithSubjectService(s *service.SubjectService) Option {
	return func(h *Handler) {
		h.subjectService = s
	}
}

func NewHandler(opts ...Option) *Handler {
	h := &Handler{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}
