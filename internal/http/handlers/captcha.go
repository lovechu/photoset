package handlers

import (
	"photoset/internal/pkg/response"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

type CaptchaHandler struct {
	captchaService service.CaptchaService
}

func NewCaptchaHandler(cs service.CaptchaService) *CaptchaHandler {
	return &CaptchaHandler{captchaService: cs}
}

// Generate 生成图形验证码
// GET /api/captcha?action=login
func (h *CaptchaHandler) Generate(c *gin.Context) {
	action := c.Query("action")
	if action == "" {
		action = "login"
	}
	if action != "login" && action != "register" {
		response.BadRequest(c, "action must be login or register")
		return
	}

	id, b64s, err := h.captchaService.Generate("digit", action)
	if err != nil {
		response.ServerError(c, "failed to generate captcha")
		return
	}

	response.Success(c, gin.H{
		"captcha_id":   id,
		"captcha_image": b64s, // base64 编码的图片，前端直接 <img :src="data:image/png;base64,..." />
	})
}
