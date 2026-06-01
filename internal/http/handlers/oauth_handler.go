package handlers

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"photoset/internal/http/middleware"
	"photoset/internal/pkg/response"
	"photoset/internal/service"
)

type OAuthHandler struct {
	oauthService service.OAuthService
}

func NewOAuthHandler(oauthService service.OAuthService) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
	}
}

// Authorize 获取授权页面信息
func (h *OAuthHandler) Authorize(c *gin.Context) {
	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")
	scope := c.Query("scope")
	state := c.Query("state")
	responseType := c.Query("response_type")

	// 验证参数
	if clientID == "" || redirectURI == "" {
		response.BadRequest(c, "client_id and redirect_uri are required")
		return
	}

	if responseType != "code" {
		response.BadRequest(c, "response_type must be 'code'")
		return
	}

	// 验证应用
	client, err := h.oauthService.GetClientByClientID(clientID)
	if err != nil {
		response.BadRequest(c, "invalid client_id")
		return
	}

	// 验证重定向 URI
	var redirectURIs []string
	if err := json.Unmarshal([]byte(client.RedirectURIs), &redirectURIs); err != nil {
		response.ServerError(c, "invalid redirect_uris format")
		return
	}
	if !containsString(redirectURIs, redirectURI) {
		response.BadRequest(c, "invalid redirect_uri")
		return
	}

	// 返回授权页面所需的信息
	response.Success(c, gin.H{
		"client": gin.H{
			"name":        client.Name,
			"description": client.Description,
			"logo_url":    client.LogoURL,
		},
		"scopes": strings.Split(scope, ","),
		"state":  state,
	})
}

// AuthorizeConfirm 用户确认授权
func (h *OAuthHandler) AuthorizeConfirm(c *gin.Context) {
	var req struct {
		ClientID    string `json:"client_id" binding:"required"`
		RedirectURI string `json:"redirect_uri" binding:"required"`
		Scope       string `json:"scope"`
		State       string `json:"state"`
		Approved    bool   `json:"approved" binding:"required"`
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

	// 如果用户拒绝授权
	if !req.Approved {
		// 重定向回客户端，带错误参数
		redirectURL := buildRedirectURL(req.RedirectURI, map[string]string{
			"error": "access_denied",
			"state": req.State,
		})
		response.Success(c, gin.H{
			"redirect_url": redirectURL,
		})
		return
	}

	// 创建授权记录
	auth, err := h.oauthService.CreateAuthorization(userID, req.ClientID, req.Scope, req.RedirectURI)
	if err != nil {
		response.Error(c, -1, err.Error())
		return
	}

	// 重定向回客户端，带授权码
	redirectURL := buildRedirectURL(req.RedirectURI, map[string]string{
		"code":  auth.Code,
		"state": req.State,
	})
	response.Success(c, gin.H{
		"redirect_url": redirectURL,
	})
}

// Token 用授权码换取访问令牌
func (h *OAuthHandler) Token(c *gin.Context) {
	var req struct {
		GrantType    string `json:"grant_type" binding:"required"`
		Code         string `json:"code"`
		RedirectURI  string `json:"redirect_uri"`
		ClientID     string `json:"client_id" binding:"required"`
		ClientSecret string `json:"client_secret" binding:"required"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	switch req.GrantType {
	case "authorization_code":
		// 用授权码换取访问令牌
		token, err := h.oauthService.ExchangeCode(req.Code, req.ClientID, req.ClientSecret, req.RedirectURI)
		if err != nil {
			response.Error(c, -1, err.Error())
			return
		}
		response.Success(c, gin.H{
			"access_token":  token.AccessToken,
			"token_type":    "Bearer",
			"expires_in":    7200, // 2小时
			"refresh_token": token.RefreshToken,
			"scope":         token.Scopes,
		})

	case "refresh_token":
		// 刷新访问令牌
		token, err := h.oauthService.RefreshToken(req.RefreshToken, req.ClientID, req.ClientSecret)
		if err != nil {
			response.Error(c, -1, err.Error())
			return
		}
		response.Success(c, gin.H{
			"access_token":  token.AccessToken,
			"token_type":    "Bearer",
			"expires_in":    7200, // 2小时
			"refresh_token": token.RefreshToken,
			"scope":         token.Scopes,
		})

	default:
		response.BadRequest(c, "unsupported grant_type")
	}
}

// Revoke 撤销令牌
func (h *OAuthHandler) Revoke(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.oauthService.RevokeToken(req.Token); err != nil {
		response.Error(c, -1, err.Error())
		return
	}

	response.Success(c, nil)
}

// UserInfo 获取用户信息
func (h *OAuthHandler) UserInfo(c *gin.Context) {
	token := extractBearerToken(c)
	if token == "" {
		response.Unauthorized(c, "missing access_token")
		return
	}

	userInfo, err := h.oauthService.GetUserInfo(token)
	if err != nil {
		response.Error(c, -1, err.Error())
		return
	}

	response.Success(c, userInfo)
}

// buildRedirectURL 构建重定向 URL
func buildRedirectURL(baseURL string, params map[string]string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

// extractBearerToken 从 Authorization header 提取 Bearer token
func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}
	return c.Query("access_token")
}

// containsString 检查字符串是否在切片中
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}