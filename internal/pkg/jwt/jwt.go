package jwt

import (
	"errors"
	"photoset/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var cfg *config.Config

func Init(config *config.Config) {
	cfg = config
}

func GenerateToken(userID uint, role string) (string, error) {
	if cfg == nil {
		return "", errors.New("JWT config not initialized")
	}

	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(cfg.JWT.ExpireHour) * time.Hour)

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	if cfg == nil {
		return nil, errors.New("JWT config not initialized")
	}

	// 密钥轮换过渡：依次尝试当前密钥和旧密钥，让存量 token 在窗口期内继续有效。
	// 签发始终只用 cfg.JWT.Secret（见 GenerateToken），旧密钥只用于验签。
	secrets := make([]string, 0, 2)
	if cfg.JWT.Secret != "" {
		secrets = append(secrets, cfg.JWT.Secret)
	}
	if cfg.JWT.PreviousSecret != "" && cfg.JWT.PreviousSecret != cfg.JWT.Secret {
		secrets = append(secrets, cfg.JWT.PreviousSecret)
	}
	if len(secrets) == 0 {
		return nil, errors.New("no JWT secret configured")
	}

	var lastErr error
	for _, secret := range secrets {
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			lastErr = err
			continue
		}
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, errors.New("invalid token")
}
