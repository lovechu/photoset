package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mojocn/base64Captcha"
	"photoset/internal/database"
)

var store = base64Captcha.DefaultMemStore

// CaptchaService 验证码服务
type CaptchaService interface {
	Generate(captchaType, action string) (string, string, error)
	Verify(captchaId, captchaValue, action string) bool
}

type captchaService struct{}

func NewCaptchaService() CaptchaService {
	return &captchaService{}
}

func (s *captchaService) Generate(captchaType, action string) (string, string, error) {
	// 生成图形验证码
	driver := base64Captcha.NewDriverDigit(
		80,  // height
		240, // width
		5,   // length
		0.7, // maxSkew
		80,  // dotCount
	)

	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := c.Generate()
	if err != nil {
		return "", "", err
	}

	// 将验证码 ID 与 action 绑定存入 Redis，5 分钟过期
	key := fmt.Sprintf("captcha:%s:%s", action, id)
	database.RedisClient.Set(context.TODO(), key, id, 5*time.Minute)

	return id, b64s, nil
}

func (s *captchaService) Verify(captchaId, captchaValue, action string) bool {
	// 1. 检查 Redis 中是否存在该验证码 key（防伪造 ID）
	key := fmt.Sprintf("captcha:%s:%s", action, captchaId)
	exists, err := database.RedisClient.Exists(context.TODO(), key).Result()
	if err != nil || exists == 0 {
		return false
	}

	// 2. 验证验证码
	if !store.Verify(captchaId, captchaValue, true) {
		return false
	}

	// 3. 验证成功后删除 Redis key（一次性使用）
	database.RedisClient.Del(context.TODO(), key)

	return true
}
