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

// @Summary      生成图形验证码
// @Description  生成登录/注册/找回密码用的图形验证码
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        action  query  string  false  "场景(login/register/forgot)"  default(login)
// @Success      200  {object}  response.Response  "验证码ID和Base64图片"
// @Failure      500  {object}  response.Response  "生成失败"
// @Router       /api/auth/captcha [get]
// Generate 生成图形验证码
// GET /api/captcha?action=login
func (h *CaptchaHandler) Generate(c *gin.Context) {
	action := c.Query("action")
	if action == "" {
		action = "login"
	}
	// 支持 login, register, forgot 三种场景
	if action != "login" && action != "register" && action != "forgot" {
		action = "login"
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
