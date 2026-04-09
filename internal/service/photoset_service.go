package service

import (
	"context"
	"errors"
	"photoset/internal/config"
	"photoset/internal/domain"
	"photoset/internal/pkg/signurl"
	"photoset/internal/repository"
	"time"

	"gorm.io/gorm"
)

type PhotoSetService struct {
	repo         *repository.PhotoSetRepository
	orderRepo    *repository.OrderRepository
	cacheService *CacheService
	cfg          *config.Config
}

func NewPhotoSetService(repo *repository.PhotoSetRepository, orderRepo *repository.OrderRepository) *PhotoSetService {
	return &PhotoSetService{
		repo:         repo,
		orderRepo:    orderRepo,
		cacheService: NewCacheService(),
		cfg:          config.Load(),
	}
}

// CreatePhotoSet 创建套图
func (s *PhotoSetService) CreatePhotoSet(photoset *domain.PhotoSet, tagNames []string, photos []domain.Photo) error {
	// 处理价格
	if photoset.IsFree == 1 {
		photoset.Price = 0
	}

	// 创建套图
	if err := s.repo.Create(photoset); err != nil {
		return err
	}

	// 处理标签
	var tagIDs []uint
	for _, tagName := range tagNames {
		tag, err := s.repo.FindTagByName(tagName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 标签不存在，创建新标签
				newTag := &domain.Tag{Name: tagName}
				if err := s.repo.CreateTag(newTag); err != nil {
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
		if err := s.repo.CreatePhotoSetTags(photoset.ID, tagIDs); err != nil {
			return err
		}
	}

	// 创建图片
	if len(photos) > 0 {
		// 设置 photoset_id
		for i := range photos {
			photos[i].PhotoSetID = photoset.ID
		}
		if err := s.repo.CreatePhotos(photos); err != nil {
			return err
		}
	}

	s.InvalidateAllPhotosetListCache()

	return nil
}

// GetPhotoSetList 获取套图列表（带 Redis 缓存）
func (s *PhotoSetService) GetPhotoSetList(page, pageSize int, tag string, keyword string, userID uint, onlyMine bool) ([]domain.PhotoSet, int64, error) {
	// 只缓存非"我的"列表（mine=true 不缓存，数据个性化）
	cacheKey := PhotosetListCacheKey(page, pageSize, tag, keyword)
	if !onlyMine {
		ctx := context.Background()
		var cached struct {
			List  []domain.PhotoSet `json:"list"`
			Total int64             `json:"total"`
		}
		if err := s.cacheService.Get(ctx, cacheKey, &cached); err == nil {
			return cached.List, cached.Total, nil
		}
	}

	photosets, total, err := s.repo.List(page, pageSize, tag, keyword, userID, onlyMine)
	if err != nil {
		return nil, 0, err
	}

	// 写入缓存（5 分钟）
	if !onlyMine {
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
