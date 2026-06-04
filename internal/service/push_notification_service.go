package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"

	"photoset/internal/domain"
	"photoset/internal/logger"
	"photoset/internal/repository"
)

type PushNotificationService struct {
	tokenRepo   *repository.PushTokenRepository
	fcmClient   *messaging.Client
	enabled     bool
}

func NewPushNotificationService(tokenRepo *repository.PushTokenRepository) *PushNotificationService {
	svc := &PushNotificationService{
		tokenRepo: tokenRepo,
		enabled:   false,
	}

	// 尝试初始化 Firebase
	svc.initFirebase()

	return svc
}

// initFirebase 初始化 Firebase Admin SDK
func (s *PushNotificationService) initFirebase() {
	// 从环境变量或文件获取 Firebase 配置
	credFile := os.Getenv("FIREBASE_CREDENTIALS_FILE")
	projectID := os.Getenv("FIREBASE_PROJECT_ID")

	if credFile == "" && projectID == "" {
		logger.Warn("Firebase 配置未设置，推送通知功能已禁用")
		return
	}

	ctx := context.Background()
	var app *firebase.App
	var err error

	if credFile != "" {
		opt := option.WithCredentialsFile(credFile)
		config := &firebase.Config{ProjectID: projectID}
		app, err = firebase.NewApp(ctx, config, opt)
	} else {
		// 使用默认凭证（在 GCP 环境中）
		config := &firebase.Config{ProjectID: projectID}
		app, err = firebase.NewApp(ctx, config)
	}

	if err != nil {
		logger.Error("Firebase 初始化失败", "error", err)
		return
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		logger.Error("Firebase Messaging 客户端创建失败", "error", err)
		return
	}

	s.fcmClient = client
	s.enabled = true
	logger.Info("Firebase 推送通知服务已启用")
}

// IsEnabled 检查推送服务是否可用
func (s *PushNotificationService) IsEnabled() bool {
	return s.enabled
}

// RegisterToken 注册推送令牌
func (s *PushNotificationService) RegisterToken(userID uint, token, platform, deviceID, deviceName string) error {
	pushToken := &domain.PushToken{
		UserID:     userID,
		Token:      token,
		Platform:   platform,
		DeviceID:   deviceID,
		DeviceName: deviceName,
		IsActive:   true,
	}
	return s.tokenRepo.Save(pushToken)
}

// UnregisterToken 注销推送令牌
func (s *PushNotificationService) UnregisterToken(token string) error {
	return s.tokenRepo.DeactivateToken(token)
}

// SendPushNotification 发送推送通知给单个用户
func (s *PushNotificationService) SendPushNotification(userID uint, title, body string, data map[string]string) error {
	if !s.enabled {
		return nil
	}

	tokens, err := s.tokenRepo.GetActiveTokensByUserID(userID)
	if err != nil {
		return fmt.Errorf("获取用户推送令牌失败: %w", err)
	}

	if len(tokens) == 0 {
		return nil
	}

	// 构建消息
	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	// 根据平台设置特定配置
	tokenStrings := make([]string, 0, len(tokens))
	for _, t := range tokens {
		tokenStrings = append(tokenStrings, t.Token)
	}
	message.Tokens = tokenStrings

	// 发送消息
	br, err := s.fcmClient.SendEachForMulticast(context.Background(), message)
	if err != nil {
		logger.Error("发送推送通知失败", "error", err, "user_id", userID)
		return fmt.Errorf("发送推送通知失败: %w", err)
	}

	// 处理失败的令牌
	for i, resp := range br.Responses {
		if !resp.Success {
			errCode := ""
			invalidToken := false
			if resp.Error != nil {
				errCode = resp.Error.Error()
				// 令牌无效则标记停用
				if strings.Contains(errCode, "not-registered") ||
				   strings.Contains(errCode, "invalid-argument") ||
				   strings.Contains(errCode, "UNREGISTERED") ||
				   strings.Contains(errCode, "INVALID_ARGUMENT") {
					invalidToken = true
				}
			}
			logger.Warn("推送通知发送失败",
				"token", tokenStrings[i][:20]+"...",
				"error", errCode,
			)
			
			if invalidToken {
				_ = s.tokenRepo.DeactivateToken(tokenStrings[i])
			}
		}
	}

	logger.Info("推送通知已发送", "user_id", userID, "success", br.SuccessCount, "failure", br.FailureCount)
	return nil
}

// SendPushToMultipleUsers 发送推送通知给多个用户
func (s *PushNotificationService) SendPushToMultipleUsers(userIDs []uint, title, body string, data map[string]string) error {
	if !s.enabled {
		return nil
	}

	for _, userID := range userIDs {
		if err := s.SendPushNotification(userID, title, body, data); err != nil {
			logger.Error("发送推送通知失败", "error", err, "user_id", userID)
			// 继续发送给其他用户
		}
	}
	return nil
}