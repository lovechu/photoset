package signurl

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strconv"
	"time"
)

// SignURL 生成签名 URL
func SignURL(rawURL, secret string, expireSeconds int) string {
	expireAt := time.Now().Unix() + int64(expireSeconds)
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	query := u.Query()
	query.Set("expires", strconv.FormatInt(expireAt, 10))
	signStr := u.Path + "?expires=" + strconv.FormatInt(expireAt, 10)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signStr))
	signature := hex.EncodeToString(mac.Sum(nil))
	query.Set("sign", signature)
	u.RawQuery = query.Encode()
	return u.String()
}

// VerifyURL 验证签名 URL
func VerifyURL(rawURL, secret string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	query := u.Query()
	expiresStr := query.Get("expires")
	sign := query.Get("sign")
	if expiresStr == "" || sign == "" {
		return false
	}
	expires, err := strconv.ParseInt(expiresStr, 10, 64)
	if err != nil {
		return false
	}
	if time.Now().Unix() > expires {
		return false
	}
	signStr := u.Path + "?expires=" + expiresStr
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signStr))
	expectedSign := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(sign), []byte(expectedSign))
}

// SignPhotoURLs 为图片列表批量签名
func SignPhotoURLs(urls []string, secret string, expireSeconds int) []string {
	result := make([]string, len(urls))
	for i, u := range urls {
		result[i] = SignURL(u, secret, expireSeconds)
	}
	return result
}
