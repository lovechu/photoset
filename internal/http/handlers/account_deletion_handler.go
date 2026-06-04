package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"photoset/internal/http/middleware"
	"photoset/internal/pkg/response"
	"photoset/internal/service"
)

type AccountDeletionHandler struct {
	deletionService *service.AccountDeletionService
}

func NewAccountDeletionHandler(deletionService *service.AccountDeletionService) *AccountDeletionHandler {
	return &AccountDeletionHandler{
		deletionService: deletionService,
	}
}

type RequestDeletionRequest struct {
	Password string `json:"password" binding:"required"`
	Reason   string `json:"reason"`
}

// @Summary      申请注销账号
// @Description  提交账号注销申请，进入30天冷静期
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body  RequestDeletionRequest  true  "注销请求 {password, reason}"
// @Success      200  {object}  object  "申请成功"
// @Failure      400  {object}  object  "参数错误"
// @Security     BearerAuth
// @Router       /api/auth/request-deletion [post]
// RequestDeletion 申请注销账号
func (h *AccountDeletionHandler) RequestDeletion(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req RequestDeletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请输入密码")
		return
	}

	if err := h.deletionService.RequestDeletion(userID, req.Password, req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注销申请已提交，30天冷静期后将自动处理",
	})
}

// @Summary      取消注销申请
// @Description  取消已提交的账号注销申请
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  object  "取消成功"
// @Failure      400  {object}  object  "操作失败"
// @Security     BearerAuth
// @Router       /api/auth/cancel-deletion [post]
// CancelDeletion 取消注销申请
func (h *AccountDeletionHandler) CancelDeletion(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	if err := h.deletionService.CancelDeletion(userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "已取消注销申请",
	})
}

// @Summary      查询注销状态
// @Description  查询当前账号的注销申请状态
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  object  "注销状态"
// @Failure      401  {object}  object  "未登录"
// @Security     BearerAuth
// @Router       /api/auth/deletion-status [get]
// GetDeletionStatus 获取注销状态
func (h *AccountDeletionHandler) GetDeletionStatus(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	status, err := h.deletionService.GetDeletionStatus(userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": status,
	})
}