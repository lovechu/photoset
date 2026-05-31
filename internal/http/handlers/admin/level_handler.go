package admin

import (
	"photoset/internal/domain"
	"photoset/internal/pkg/response"
	"photoset/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminLevelHandler handles admin level/achievement/points-mall requests
type AdminLevelHandler struct {
	levelRepo      *repository.UserLevelRepository
	achievementRepo *repository.AchievementRepository
	pointsMallRepo  *repository.PointsMallRepository
	pointRepo       *repository.UserPointRepository
	db              *gorm.DB
}

// NewAdminLevelHandler creates a new AdminLevelHandler
func NewAdminLevelHandler(db *gorm.DB) *AdminLevelHandler {
	return &AdminLevelHandler{
		levelRepo:       repository.NewUserLevelRepository(db),
		achievementRepo: repository.NewAchievementRepository(db),
		pointsMallRepo:  repository.NewPointsMallRepository(db),
		pointRepo:       repository.NewUserPointRepository(db),
		db:              db,
	}
}

// ============ Level Config Management ============

// GetLevelConfigs returns all level configurations
func (h *AdminLevelHandler) GetLevelConfigs(c *gin.Context) {
	configs, err := h.levelRepo.GetAllLevelConfigs()
	if err != nil {
		response.ServerError(c, "failed to get level configs")
		return
	}
	response.Success(c, configs)
}

// UpdateLevelConfig updates a level configuration
func (h *AdminLevelHandler) UpdateLevelConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid level id")
		return
	}

	var req domain.UserLevelConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	existing, err := h.levelRepo.GetLevelConfigByID(uint(id))
	if err != nil {
		response.NotFound(c, "level config not found")
		return
	}

	// Update allowed fields
	existing.Name = req.Name
	existing.Icon = req.Icon
	existing.Color = req.Color
	existing.MinPoints = req.MinPoints
	existing.MaxPoints = req.MaxPoints
	existing.Description = req.Description
	existing.CanCreatePost = req.CanCreatePost
	existing.CanCreateReply = req.CanCreateReply
	existing.CanUploadImage = req.CanUploadImage
	existing.CanUploadVideo = req.CanUploadVideo
	existing.CanCreateTopic = req.CanCreateTopic
	existing.CanPinPost = req.CanPinPost
	existing.CanDeleteReply = req.CanDeleteReply
	existing.MaxPostPerDay = req.MaxPostPerDay
	existing.MaxReplyPerDay = req.MaxReplyPerDay
	existing.MaxImagePerPost = req.MaxImagePerPost
	existing.MaxVideoPerPost = req.MaxVideoPerPost
	existing.MaxPostLength = req.MaxPostLength
	existing.RewardPoints = req.RewardPoints
	existing.RewardBadge = req.RewardBadge
	existing.RewardTitle = req.RewardTitle

	if err := h.levelRepo.UpdateLevelConfig(existing); err != nil {
		response.ServerError(c, "failed to update level config")
		return
	}

	response.Success(c, existing)
}

// ============ Achievement Management ============

// GetAchievements returns all achievements (admin, including hidden)
func (h *AdminLevelHandler) GetAchievements(c *gin.Context) {
	achievements, err := h.achievementRepo.GetAllAchievementsIncludingHidden()
	if err != nil {
		response.ServerError(c, "failed to get achievements")
		return
	}
	response.Success(c, achievements)
}

// CreateAchievement creates a new achievement
func (h *AdminLevelHandler) CreateAchievement(c *gin.Context) {
	var req domain.Achievement
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.achievementRepo.CreateAchievement(&req); err != nil {
		response.ServerError(c, "failed to create achievement")
		return
	}

	response.Success(c, req)
}

// UpdateAchievement updates an achievement
func (h *AdminLevelHandler) UpdateAchievement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid achievement id")
		return
	}

	var req domain.Achievement
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	existing, err := h.achievementRepo.GetAchievementByID(uint(id))
	if err != nil {
		response.NotFound(c, "achievement not found")
		return
	}

	// Update fields
	existing.Name = req.Name
	existing.Title = req.Title
	existing.Description = req.Description
	existing.Icon = req.Icon
	existing.BadgeImage = req.BadgeImage
	existing.Type = req.Type
	existing.ConditionType = req.ConditionType
	existing.ConditionValue = req.ConditionValue
	existing.RewardPoints = req.RewardPoints
	existing.RewardTitle = req.RewardTitle
	existing.SortOrder = req.SortOrder
	existing.IsHidden = req.IsHidden

	if err := h.achievementRepo.UpdateAchievement(existing); err != nil {
		response.ServerError(c, "failed to update achievement")
		return
	}

	response.Success(c, existing)
}

// DeleteAchievement soft-deletes an achievement
func (h *AdminLevelHandler) DeleteAchievement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid achievement id")
		return
	}

	if err := h.achievementRepo.DeleteAchievement(uint(id)); err != nil {
		response.ServerError(c, "failed to delete achievement")
		return
	}

	response.Success(c, gin.H{"message": "achievement deleted"})
}

// ============ Points Mall Management ============

// GetPointsMallItems returns all points mall items (admin)
func (h *AdminLevelHandler) GetPointsMallItems(c *gin.Context) {
	items, err := h.pointsMallRepo.GetAllItemsIncludingInactive()
	if err != nil {
		response.ServerError(c, "failed to get items")
		return
	}
	response.Success(c, items)
}

// CreatePointsMallItem creates a new points mall item
func (h *AdminLevelHandler) CreatePointsMallItem(c *gin.Context) {
	var req domain.PointsMallItem
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.pointsMallRepo.CreateItem(&req); err != nil {
		response.ServerError(c, "failed to create item")
		return
	}

	response.Success(c, req)
}

// UpdatePointsMallItem updates a points mall item
func (h *AdminLevelHandler) UpdatePointsMallItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid item id")
		return
	}

	var req domain.PointsMallItem
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	existing, err := h.pointsMallRepo.GetItemByID(uint(id))
	if err != nil {
		response.NotFound(c, "item not found")
		return
	}

	// Update fields
	existing.Name = req.Name
	existing.Description = req.Description
	existing.Image = req.Image
	existing.Category = req.Category
	existing.PointsCost = req.PointsCost
	existing.ItemType = req.ItemType
	existing.ItemValue = req.ItemValue
	existing.TotalStock = req.TotalStock
	existing.IsUnlimited = req.IsUnlimited
	existing.MinLevel = req.MinLevel
	existing.IsActive = req.IsActive
	existing.SortOrder = req.SortOrder

	if err := h.pointsMallRepo.UpdateItem(existing); err != nil {
		response.ServerError(c, "failed to update item")
		return
	}

	response.Success(c, existing)
}

// DeletePointsMallItem soft-deletes a points mall item
func (h *AdminLevelHandler) DeletePointsMallItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid item id")
		return
	}

	if err := h.pointsMallRepo.DeleteItem(uint(id)); err != nil {
		response.ServerError(c, "failed to delete item")
		return
	}

	response.Success(c, gin.H{"message": "item deleted"})
}

// ============ Exchange History ============

// GetAllExchanges returns all exchange records (admin)
func (h *AdminLevelHandler) GetAllExchanges(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	exchanges, total, err := h.pointsMallRepo.GetAllExchanges(page, pageSize)
	if err != nil {
		response.ServerError(c, "failed to get exchanges")
		return
	}

	response.Success(c, gin.H{
		"exchanges": exchanges,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}
