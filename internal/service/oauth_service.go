package service

import (
	"encoding/json"
	"errors"
	"strings"

	"photoset/internal/domain"
	"photoset/internal/repository"
)

type OAuthService interface {
	// 应用管理
	CreateClient(name, description string, redirectURIs []string, scopes []string, logoURL string, createdBy uint) (*domain.OAuthClient, error)
	UpdateClient(id uint, name, description string, redirectURIs []string, scopes []string, logoURL string, status int) error
	DeleteClient(id uint) error
	GetClientByID(id uint) (*domain.OAuthClient, error)
	GetClientByClientID(clientID string) (*domain.OAuthClient, error)
	ListClients() ([]domain.OAuthClient, error)

	// 授权流程
	ValidateClient(clientID, clientSecret, redirectURI string) (*domain.OAuthClient, error)
	CreateAuthorization(userID uint, clientID, scopes, redirectURI string) (*domain.OAuthAuthorization, error)
	ExchangeCode(code, clientID, clientSecret, redirectURI string) (*domain.OAuthToken, error)
	RefreshToken(refreshToken, clientID, clientSecret string) (*domain.OAuthToken, error)
	RevokeToken(accessToken string) error

	// 用户信息
	GetUserInfo(accessToken string) (map[string]interface{}, error)
	GetUserScopes(accessToken string) ([]string, error)
}

type oauthService struct {
	clientRepo      *repository.OAuthClientRepository
	authorizationRepo *repository.OAuthAuthorizationRepository
	tokenRepo       *repository.OAuthTokenRepository
	userRepo        repository.UserRepository
}

func NewOAuthService(
	clientRepo *repository.OAuthClientRepository,
	authorizationRepo *repository.OAuthAuthorizationRepository,
	tokenRepo *repository.OAuthTokenRepository,
	userRepo repository.UserRepository,
) OAuthService {
	return &oauthService{
		clientRepo:      clientRepo,
		authorizationRepo: authorizationRepo,
		tokenRepo:       tokenRepo,
		userRepo:        userRepo,
	}
}

// CreateClient 创建 OAuth 应用
func (s *oauthService) CreateClient(name, description string, redirectURIs []string, scopes []string, logoURL string, createdBy uint) (*domain.OAuthClient, error) {
	// 验证重定向 URI
	if len(redirectURIs) == 0 {
		return nil, errors.New("redirect_uris is required")
	}

	// 验证权限范围
	validScopes := []string{"userinfo", "userinfo:email", "photosets:read", "favorites:read", "favorites:write", "community:read", "community:write"}
	for _, scope := range scopes {
		if !contains(validScopes, scope) {
			return nil, errors.New("invalid scope: " + scope)
		}
	}

	return s.clientRepo.Create(name, description, redirectURIs, scopes, logoURL, createdBy)
}

// UpdateClient 更新 OAuth 应用
func (s *oauthService) UpdateClient(id uint, name, description string, redirectURIs []string, scopes []string, logoURL string, status int) error {
	// 验证重定向 URI
	if len(redirectURIs) == 0 {
		return errors.New("redirect_uris is required")
	}

	// 验证权限范围
	validScopes := []string{"userinfo", "userinfo:email", "photosets:read", "favorites:read", "favorites:write", "community:read", "community:write"}
	for _, scope := range scopes {
		if !contains(validScopes, scope) {
			return errors.New("invalid scope: " + scope)
		}
	}

	return s.clientRepo.Update(id, name, description, redirectURIs, scopes, logoURL, status)
}

// DeleteClient 删除 OAuth 应用
func (s *oauthService) DeleteClient(id uint) error {
	return s.clientRepo.Delete(id)
}

// GetClientByID 根据 ID 获取应用
func (s *oauthService) GetClientByID(id uint) (*domain.OAuthClient, error) {
	return s.clientRepo.GetByID(id)
}

// GetClientByClientID 根据 ClientID 获取应用
func (s *oauthService) GetClientByClientID(clientID string) (*domain.OAuthClient, error) {
	return s.clientRepo.GetByClientID(clientID)
}

// ListClients 获取所有应用
func (s *oauthService) ListClients() ([]domain.OAuthClient, error) {
	return s.clientRepo.List()
}

// ValidateClient 验证应用
func (s *oauthService) ValidateClient(clientID, clientSecret, redirectURI string) (*domain.OAuthClient, error) {
	client, err := s.clientRepo.GetByClientID(clientID)
	if err != nil {
		return nil, errors.New("invalid client_id")
	}

	// 验证 ClientSecret
	if !s.clientRepo.ValidateSecret(client, clientSecret) {
		return nil, errors.New("invalid client_secret")
	}

	// 验证重定向 URI
	var redirectURIs []string
	if err := json.Unmarshal([]byte(client.RedirectURIs), &redirectURIs); err != nil {
		return nil, errors.New("invalid redirect_uris format")
	}
	if !contains(redirectURIs, redirectURI) {
		return nil, errors.New("invalid redirect_uri")
	}

	return client, nil
}

// CreateAuthorization 创建授权记录
func (s *oauthService) CreateAuthorization(userID uint, clientID, scopes, redirectURI string) (*domain.OAuthAuthorization, error) {
	// 验证应用
	client, err := s.clientRepo.GetByClientID(clientID)
	if err != nil {
		return nil, errors.New("invalid client_id")
	}

	// 验证应用请求的权限范围是否在其注册范围内
	clientScopes := strings.Split(client.Scopes, ",")
	requestedScopes := strings.Split(scopes, ",")
	for _, requestedScope := range requestedScopes {
		if !contains(clientScopes, requestedScope) {
			return nil, errors.New("scope not allowed: " + requestedScope)
		}
	}

	return s.authorizationRepo.Create(userID, clientID, scopes)
}

// ExchangeCode 用授权码换取访问令牌
func (s *oauthService) ExchangeCode(code, clientID, clientSecret, redirectURI string) (*domain.OAuthToken, error) {
	// 验证应用
	_, err := s.ValidateClient(clientID, clientSecret, redirectURI)
	if err != nil {
		return nil, err
	}

	// 获取授权记录
	auth, err := s.authorizationRepo.GetByCode(code)
	if err != nil {
		return nil, errors.New("invalid or expired code")
	}

	// 验证授权码对应的客户端
	if auth.ClientID != clientID {
		return nil, errors.New("code was not issued to this client")
	}

	// 创建访问令牌
	token, err := s.tokenRepo.Create(auth.UserID, auth.ClientID, auth.Scopes)
	if err != nil {
		return nil, err
	}

	// 删除已使用的授权码
	s.authorizationRepo.Delete(auth.ID)

	return token, nil
}

// RefreshToken 刷新访问令牌
func (s *oauthService) RefreshToken(refreshToken, clientID, clientSecret string) (*domain.OAuthToken, error) {
	// 获取旧令牌
	oldToken, err := s.tokenRepo.GetByRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh_token")
	}

	// 验证客户端
	if oldToken.ClientID != clientID {
		return nil, errors.New("refresh_token was not issued to this client")
	}

	client, err := s.clientRepo.GetByClientID(clientID)
	if err != nil {
		return nil, errors.New("invalid client_id")
	}

	// 验证 ClientSecret
	if !s.clientRepo.ValidateSecret(client, clientSecret) {
		return nil, errors.New("invalid client_secret")
	}

	// 验证刷新令牌
	if !s.tokenRepo.ValidateRefreshToken(oldToken, refreshToken) {
		return nil, errors.New("invalid refresh_token")
	}

	// 刷新令牌
	return s.tokenRepo.Refresh(oldToken)
}

// RevokeToken 撤销访问令牌
func (s *oauthService) RevokeToken(accessToken string) error {
	token, err := s.tokenRepo.GetByAccessToken(accessToken)
	if err != nil {
		return errors.New("invalid access_token")
	}

	return s.tokenRepo.Revoke(token.ID)
}

// GetUserInfo 获取用户信息
func (s *oauthService) GetUserInfo(accessToken string) (map[string]interface{}, error) {
	// 获取令牌
	token, err := s.tokenRepo.GetByAccessToken(accessToken)
	if err != nil {
		return nil, errors.New("invalid access_token")
	}

	// 验证访问令牌
	if !s.tokenRepo.ValidateAccessToken(token, accessToken) {
		return nil, errors.New("invalid access_token")
	}

	// 获取用户信息
	user, err := s.userRepo.FindByID(token.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 根据权限范围返回用户信息
	userInfo := make(map[string]interface{})
	scopes := strings.Split(token.Scopes, ",")

	if contains(scopes, "userinfo") {
		userInfo["id"] = user.ID
		userInfo["nickname"] = user.Nickname
		userInfo["avatar"] = user.Avatar
		userInfo["bio"] = user.Bio
	}

	if contains(scopes, "userinfo:email") {
		userInfo["email"] = user.Email
	}

	return userInfo, nil
}

// GetUserScopes 获取用户授权的权限范围
func (s *oauthService) GetUserScopes(accessToken string) ([]string, error) {
	token, err := s.tokenRepo.GetByAccessToken(accessToken)
	if err != nil {
		return nil, errors.New("invalid access_token")
	}

	// 验证访问令牌
	if !s.tokenRepo.ValidateAccessToken(token, accessToken) {
		return nil, errors.New("invalid access_token")
	}

	return strings.Split(token.Scopes, ","), nil
}

// contains 检查字符串是否在切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}