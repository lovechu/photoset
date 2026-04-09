package handlers

import (
	"net/http"

	"photoset/internal/pkg/response"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	service *service.PhotoSetService
}

func NewTagHandler(service *service.PhotoSetService) *TagHandler {
	return &TagHandler{service: service}
}

// List 标签列表
func (h *TagHandler) List(c *gin.Context) {
	tags, err := h.service.GetAllTags()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取标签列表失败")
		return
	}

	response.Success(c, tags)
}
