package handler

import "github.com/fun-dotto/app-bff-api/internal/service"

type Handler struct {
	announcementService *service.AnnouncementService
	//今後ここに追加
}

func NewHandler(announcementService *service.AnnouncementService /*今後ここに追加*/) *Handler {
	return &Handler{announcementService: announcementService}
}
