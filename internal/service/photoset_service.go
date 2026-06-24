package service

import (
	"context"
	"errors"
	"photoset/internal/config"
	"photoset/internal/domain"
	"photoset/internal/logger"
	"photoset/internal/pkg/signurl"
	"photoset/internal/repository"
	"photoset/internal/storage"
	"time"

	"gorm.io/gorm"
)

type PhotoSetService struct {
	repo         *repository.PhotoSetRepository
	orderRepo    *repository.OrderRepository
	cacheService *CacheService
	cfg          *config.Config
	storage      storage.Storage
}

func NewPhotoSetService(repo *repository.PhotoSetRepository, orderRepo *repository.OrderRepository, cfg *config.Config, stor storage.Storage) *PhotoSetService {
	return &PhotoSetService{
		repo:         repo,
		orderRepo:    orderRepo,
		cacheService: NewCacheService(),
		cfg:          cfg,
		storage:      stor,
	}
}

// CreatePhotoSet 创建套图（事务保护）
func (s *PhotoSetService) CreatePhotoSet(photoset *domain.PhotoSet, tagNames []string, photos []domain.Photo) error {
	// 处理价格
	if photoset.IsFree == 1 {
		photoset.Price = 0
	}

	// 使用事务确保数据一致性
	err := s.repo.Transaction(func(tx *gorm.DB) error {
		// 创建套图
		if err := tx.Create(photoset).Error; err != nil {
			return err
		}

		// 处理标签
		var tagIDs []uint
		for _, tagName := range tagNames {
			var tag domain.Tag
			if err := tx.Where("name = ?", tagName).First(&tag).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 标签不存在，创建新标签
					newTag := &domain.Tag{Name: tagName}
					if err := tx.Create(newTag).Error; err != nil {
						return err
					}
					tagIDs = append(tagIDs, newTag.ID)
				} else {
					return err
				}
			} else {
				tagIDs = append(tagIDs, tag.ID)
			}
		}

		// 创建套图标签关联
		if len(tagIDs) > 0 {
			for _, tagID := range tagIDs {
				photosetTag := map[string]interface{}{
					"photoset_id": photoset.ID,
					"tag_id":      tagID,
				}
				if err := tx.Table("photoset_tags").Create(&photosetTag).Error; err != nil {
					return err
				}
			}
		}

		// 创建图片
		if len(photos) > 0 {
			for i := range photos {
				photos[i].PhotoSetID = photoset.ID
			}
			if err := tx.Create(&photos).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	s.InvalidateAllPhotosetListCache()
	// 同步写入付费状态缓存
	SetPaidStatus(photoset.ID, photoset.IsFree == 0)

	return nil
}

// GetPhotoSetList 获取套图列表（带 Redis 缓存）- 兼容旧接口
func (s *PhotoSetService) GetPhotoSetList(page, pageSize int, tag string, keyword string, userID uint, onlyMine bool) ([]domain.PhotoSet, int64, error) {
	return s.GetPhotoSetListAdvanced(page, pageSize, tag, keyword, userID, onlyMine, "", 0, 0, nil, "", "", 0, "")
}

// GetPhotoSetListAdvanced 高级套图列表查询（带 Redis 缓存）
func (s *PhotoSetService) GetPhotoSetListAdvanced(
	page, pageSize int,
	tag string, keyword string,
	userID uint, onlyMine bool,
	category string,
	priceMin, priceMax float64,
	isFree *bool,
	sortBy, timeRange string,
	filterUserID uint,
	status string,
) ([]domain.PhotoSet, int64, error) {
	// 生成高级缓存键
	cacheKey := PhotosetAdvancedListCacheKey(
		page, pageSize, tag, keyword, userID, onlyMine,
		category, priceMin, priceMax, isFree, sortBy, timeRange, filterUserID, status,
	)
	
	// 只缓存非"我的"列表（mine=true 或指定了status都不缓存，数据个性化）
	if !onlyMine && status == "" {
		ctx := context.Background()
		var cached struct {
			List  []domain.PhotoSet `json:"list"`
			Total int64             `json:"total"`
		}
		if err := s.cacheService.Get(ctx, cacheKey, &cached); err == nil {
			return cached.List, cached.Total, nil
		}
	}

	// 调用增强的 Repository 方法
	photosets, total, err := s.repo.ListAdvanced(
		page, pageSize, tag, keyword, userID, onlyMine,
		category, priceMin, priceMax, isFree, sortBy, timeRange, filterUserID, status,
	)
	if err != nil {
		return nil, 0, err
	}

	// 写入缓存（5 分钟）
	if !onlyMine && status == "" {
		ctx := context.Background()
		s.cacheService.Set(ctx, cacheKey, map[string]interface{}{
			"list":  photosets,
			"total": total,
		}, 5*time.Minute)
	}

	return photosets, total, nil
}

// GetPhotoSetDetail 获取套图详情（带 Redis 缓存 + URL 签名）
func (s *PhotoSetService) GetPhotoSetDetail(id uint) (*domain.PhotoSet, error) {
	ctx := context.Background()
	cacheKey := PhotosetDetailCacheKey(id)

	var cached domain.PhotoSet
	if err := s.cacheService.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	photoset, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 对付费套图的图片 URL 进行签名
	if photoset.IsFree == 0 && s.cfg.Storage.SignSecret != "" {
		expire := s.cfg.Storage.SignExpire
		if expire <= 0 {
			expire = 7200
		}
		for i := range photoset.Photos {
			photoset.Photos[i].URL = signurl.SignURL(photoset.Photos[i].URL, s.cfg.Storage.SignSecret, expire)
		}
		photoset.Cover = signurl.SignURL(photoset.Cover, s.cfg.Storage.SignSecret, expire)
	}

	// 写入缓存（10 分钟）
	s.cacheService.Set(ctx, cacheKey, photoset, 10*time.Minute)
	return photoset, nil
}

// GetPhotoSetDetailBasic 获取套图基础信息（不含完整图片列表）
func (s *PhotoSetService) GetPhotoSetDetailBasic(id uint) (*domain.PhotoSet, error) {
	return s.repo.FindByIDWithoutPhotos(id)
}

// GetPhotoSetDownload 获取套图下载信息（验证权限后返回图片 URL 列表）
func (s *PhotoSetService) GetPhotoSetDownload(id uint, userRole string, userID uint, isLoggedIn bool) ([]string, error) {
	// 1. 获取套图基础信息
	photoset, err := s.repo.FindByIDWithoutPhotos(id)
	if err != nil {
		return nil, errors.New("套图不存在")
	}

	// 2. 验证查看权限
	canViewFull := s.CanViewFullPhotos(photoset, userRole, userID, isLoggedIn)
	if !canViewFull {
		return nil, errors.New("无权下载此套图，请先购买")
	}

	// 3. 获取完整套图信息（含图片列表）
	photoset, err = s.GetPhotoSetDetail(id)
	if err != nil {
		return nil, errors.New("获取套图详情失败")
	}

	// 4. 收集图片 URL
	var photoURLs []string
	for _, photo := range photoset.Photos {
		photoURLs = append(photoURLs, photo.URL)
	}

	return photoURLs, nil
}

// CanViewFullPhotos 判断用户是否可以查看完整图片列表
func (s *PhotoSetService) CanViewFullPhotos(photoset *domain.PhotoSet, userRole string, userID uint, isLoggedIn bool) bool {
	// 如果是免费套图，任何人都可以查看
	if photoset.IsFree == 1 {
		return true
	}

	// 如果未登录，不能查看付费套图
	if !isLoggedIn {
		return false
	}

	// 如果是作者本人，可以查看
	if userID == photoset.UserID {
		return true
	}

	// 如果是管理员或会员，可以查看
	if userRole == "admin" || userRole == "member" {
		return true
	}

	// 如果用户已购买该套图，可以查看（单套图购买）
	if s.orderRepo != nil {
		hasPaid, err := s.orderRepo.HasPaidOrder(userID, photoset.ID)
		if err == nil && hasPaid {
			return true
		}
	}

	return false
}

// GetAllTags 获取所有标签（带 Redis 缓存）
func (s *PhotoSetService) GetAllTags() ([]domain.Tag, error) {
	ctx := context.Background()

	var cached []domain.Tag
	if err := s.cacheService.Get(ctx, CachePrefixTags, &cached); err == nil {
		return cached, nil
	}

	tags, err := s.repo.ListTags()
	if err != nil {
		return nil, err
	}

	// 写入缓存（30 分钟）
	s.cacheService.Set(ctx, CachePrefixTags, tags, 30*time.Minute)
	return tags, nil
}

// UpdatePhotoSet 更新套图（含标签、图片替换）
func (s *PhotoSetService) UpdatePhotoSet(id uint, updates map[string]interface{}, tags []string, photos []domain.Photo) error {
	if err := s.repo.Update(id, updates); err != nil {
		return err
	}
	if err := s.repo.ReplaceTags(id, tags); err != nil {
		return err
	}
	if photos != nil {
		if err := s.repo.ReplacePhotos(id, photos); err != nil {
			return err
		}
	}

	s.InvalidatePhotosetCache(id)
	s.InvalidateAllPhotosetListCache()
	s.InvalidateTagsCache()
	s.InvalidateCategoriesCache()

	// 如果 is_free 发生变化，同步刷新付费状态缓存
	if isFreeVal, ok := updates["is_free"]; ok {
		isFree := 0
		switch v := isFreeVal.(type) {
		case int8:
			isFree = int(v)
		case int:
			isFree = v
		case float64:
			isFree = int(v)
		}
		SetPaidStatus(id, isFree == 0)
	}

	return nil
}

// DeletePhotoSet 删除套图（软删除）
func (s *PhotoSetService) DeletePhotoSet(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	s.InvalidatePhotosetCache(id)
	s.InvalidateAllPhotosetListCache()
	s.InvalidateTagsCache()
	s.InvalidateCategoriesCache()
	InvalidatePaidStatus(id)
	return nil
}

// InvalidatePhotosetCache 清除套图相关缓存（创建/更新/删除时调用）
func (s *PhotoSetService) InvalidatePhotosetCache(id uint) {
	ctx := context.Background()
	// 清除详情缓存
	s.cacheService.Delete(ctx, PhotosetDetailCacheKey(id))
	// 清除所有列表缓存（因为列表数据包含此套图）
	s.cacheService.DeleteByPattern(ctx, CachePrefixPhotosetList+"*")
}

// InvalidateAllPhotosetListCache 清除所有套图列表缓存
func (s *PhotoSetService) InvalidateAllPhotosetListCache() {
	ctx := context.Background()
	s.cacheService.DeleteByPattern(ctx, CachePrefixPhotosetList+"*")
}

// InvalidateTagsCache 清除标签缓存
func (s *PhotoSetService) InvalidateTagsCache() {
	ctx := context.Background()
	s.cacheService.Delete(ctx, CachePrefixTags)
}

// ============ Category 服务层 ============

// GetAllCategories 获取所有分类（带 Redis 缓存 30 分钟）
func (s *PhotoSetService) GetAllCategories() ([]domain.Category, error) {
	ctx := context.Background()
	cacheKey := "categories:all"

	var cached []domain.Category
	if err := s.cacheService.Get(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	categories, err := s.repo.ListCategories()
	if err != nil {
		return nil, err
	}

	s.cacheService.Set(ctx, cacheKey, categories, 30*time.Minute)
	return categories, nil
}

// CreateCategory 创建分类 + 清缓存
func (s *PhotoSetService) CreateCategory(cat *domain.Category) error {
	if err := s.repo.CreateCategory(cat); err != nil {
		return err
	}
	s.InvalidateCategoriesCache()
	return nil
}

// UpdateCategory 更新分类 + 清缓存
func (s *PhotoSetService) UpdateCategory(id uint, updates map[string]interface{}) error {
	if err := s.repo.UpdateCategory(id, updates); err != nil {
		return err
	}
	s.InvalidateCategoriesCache()
	s.InvalidateAllPhotosetListCache()
	return nil
}

// DeleteCategory 删除分类 + 清缓存
func (s *PhotoSetService) DeleteCategory(id uint) error {
	if err := s.repo.DeleteCategory(id); err != nil {
		return err
	}
	s.InvalidateCategoriesCache()
	s.InvalidateAllPhotosetListCache()
	return nil
}

// InvalidateCategoriesCache 清除分类列表缓存
func (s *PhotoSetService) InvalidateCategoriesCache() {
	ctx := context.Background()
	s.cacheService.Delete(ctx, "categories:all")
}

// ============ 回收站功能 ============

// GetTrashList 获取回收站列表
func (s *PhotoSetService) GetTrashList(userID uint) ([]domain.PhotoSet, error) {
	return s.repo.GetTrash(userID)
}

// AdminGetTrashList 管理员获取回收站列表（全站，支持分页）
func (s *PhotoSetService) AdminGetTrashList(page, pageSize int) ([]domain.PhotoSet, int64, error) {
	return s.repo.AdminGetTrash(page, pageSize)
}

// RestorePhotoSet 恢复已删除的套图
func (s *PhotoSetService) RestorePhotoSet(id uint, userID uint) error {
	// 验证所有权
	photoset, err := s.repo.FindByIDWithoutPhotos(id)
	if err != nil {
		return errors.New("套图不存在")
	}
	if photoset.UserID != userID {
		return errors.New("无权恢复此套图")
	}

	if err := s.repo.Restore(id); err != nil {
		return err
	}

	s.InvalidateAllPhotosetListCache()
	return nil
}

// AdminRestorePhotoSet 管理员恢复套图（无需验证所有权）
func (s *PhotoSetService) AdminRestorePhotoSet(id uint) error {
	if err := s.repo.Restore(id); err != nil {
		return err
	}

	s.InvalidateAllPhotosetListCache()
	return nil
}

// PermanentDeletePhotoSet 永久删除套图（含物理文件清理）
func (s *PhotoSetService) PermanentDeletePhotoSet(id uint, userID uint) error {
	// 验证所有权
	photoset, err := s.repo.FindByIDWithoutPhotos(id)
	if err != nil {
		return errors.New("套图不存在")
	}
	if photoset.UserID != userID {
		return errors.New("无权删除此套图")
	}

	// 先查询完整套图信息（含图片列表），用于收集文件 URL
	fullPhotoset, _ := s.repo.FindByID(id)

	// 永久删除数据库记录
	if err := s.repo.PermanentDelete(id); err != nil {
		return err
	}

	// 异步清理物理文件（不阻塞主流程，失败仅记录日志）
	if fullPhotoset != nil {
		go s.cleanupPhotosetFiles(fullPhotoset)
	}

	s.InvalidatePhotosetCache(id)
	s.InvalidateAllPhotosetListCache()
	return nil
}

// AdminPermanentDeletePhotoSet 管理员永久删除套图（含物理文件清理，无需验证所有权）
func (s *PhotoSetService) AdminPermanentDeletePhotoSet(id uint) error {
	// 先查询完整套图信息（含图片列表），用于收集文件 URL
	fullPhotoset, _ := s.repo.FindByID(id)

	// 永久删除数据库记录
	if err := s.repo.PermanentDelete(id); err != nil {
		return err
	}

	// 异步清理物理文件
	if fullPhotoset != nil {
		go s.cleanupPhotosetFiles(fullPhotoset)
	}

	s.InvalidatePhotosetCache(id)
	s.InvalidateAllPhotosetListCache()
	return nil
}

// cleanupPhotosetFiles 清理套图关联的所有物理文件
func (s *PhotoSetService) cleanupPhotosetFiles(photoset *domain.PhotoSet) {
	if s.storage == nil {
		return
	}

	var urls []string
	// 封面图
	if photoset.Cover != "" {
		urls = append(urls, photoset.Cover)
	}
	// 所有图片
	for _, photo := range photoset.Photos {
		if photo.URL != "" {
			urls = append(urls, photo.URL)
		}
	}
	// 评论中的图片
	commentImageURLs := s.repo.GetCommentImageURLs(photoset.ID)
	urls = append(urls, commentImageURLs...)

	for _, url := range urls {
		if err := s.storage.Delete(url); err != nil {
			logger.Warn("清理套图文件失败", "url", url, "error", err)
		}
	}
}
