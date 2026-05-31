package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"photoset/internal/http/middleware"
	"photoset/internal/pkg/response"
	"photoset/internal/service"
)

type AdminOAuthHandler struct {
	oauthService service.OAuthService
}

func NewAdminOAuthHandler(oauthService service.OAuthService) *AdminOAuthHandler {
	return &AdminOAuthHandler{
		oauthService: oauthService,
	}
}

// GetClients 获取所有 OAuth 应用
func (h *AdminOAuthHandler) GetClients(c *gin.Context) {
	clients, err := h.oauthService.ListClients()
	if err != nil {
		response.ServerError(c, "failed to get OAuth clients")
		return
	}
	response.Success(c, clients)
}

// GetClient 获取单个 OAuth 应用
func (h *AdminOAuthHandler) GetClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	client, err := h.oauthService.GetClientByID(uint(id))
	if err != nil {
		response.NotFound(c, "OAuth client not found")
		return
	}
	response.Success(c, client)
}

// CreateClient 创建 OAuth 应用
func (h *AdminOAuthHandler) CreateClient(c *gin.Context) {
	var req struct {
		Name         string   `json:"name" binding:"required"`
		Description  string   `json:"description"`
		RedirectURIs []string `json:"redirect_uris" binding:"required"`
		Scopes       []string `json:"scopes" binding:"required"`
		LogoURL      string   `json:"logo_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID, exists := middleware.GetOAuthUserID(c)
	if !exists {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	client, err := h.oauthService.CreateClient(req.Name, req.Description, req.RedirectURIs, req.Scopes, req.LogoURL, userID)
	if err != nil {
		response.Error(c, -1, err.Error())
		return
	}
	response.Success(c, client)
}

// UpdateClient 更新 OAuth 应用
func (h *AdminOAuthHandler) UpdateClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	var req struct {
		Name         string   `json:"name" binding:"required"`
		Description  string   `json:"description"`
		RedirectURIs []string `json:"redirect_uris" binding:"required"`
		Scopes       []string `json:"scopes" binding:"required"`
		LogoURL      string   `json:"logo_url"`
		Status       int      `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	err = h.oauthService.UpdateClient(uint(id), req.Name, req.Description, req.RedirectURIs, req.Scopes, req.LogoURL, req.Status)
	if err != nil {
		response.Error(c, -1, err.Error())
		return
	}
	response.Success(c, nil)
}

// DeleteClient 删除 OAuth 应用
func (h *AdminOAuthHandler) DeleteClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	err = h.oauthService.DeleteClient(uint(id))
	if err != nil {
		response.Error(c, -1, err.Error())
		return
	}
	response.Success(c, nil)
}