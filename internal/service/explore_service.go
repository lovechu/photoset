package service

import (
	"time"

	"photoset/internal/domain"
	"photoset/internal/repository"

	"gorm.io/gorm"
)

// ExploreService provides explore/discover feed functionality
type ExploreService struct {
	db            *gorm.DB
	photosetRepo  *repository.PhotoSetRepository
}

// NewExploreService creates a new ExploreService
func NewExploreService(db *gorm.DB, photosetRepo *repository.PhotoSetRepository) *ExploreService {
	return &ExploreService{
		db:           db,
		photosetRepo: photosetRepo,
	}
}

// ExploreItem represents an item in the explore feed
type ExploreItem struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	CoverImage  string  `json:"cover_image"`
	PhotoCount  int     `json:"photo_count"`
	ViewCount   int     `json:"view_count"`
	FavCount    int     `json:"fav_count"`
	Price       float64 `json:"price"`
	IsVip       bool    `json:"is_vip"`
	Category    string  `json:"category"`
	Tags        []string `json:"tags"`
	Source      string  `json:"source"` // hot, new, recommend
	CreatedAt   string  `json:"created_at"`
}

// ExploreFeedResponse represents the explore feed response
type ExploreFeedResponse struct {
	Items []ExploreItem `json:"items"`
	Page  int           `json:"page"`
	Total int           `json:"total"`
	HasMore bool        `json:"has_more"`
}

// GetExploreFeed 获取探索发现 feed（真实分页，混合热度与时间排序）
func (s *ExploreService) GetExploreFeed(userID uint, page, pageSize int) (*ExploreFeedResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 50 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	// 用子查询统计真实热度，按 (浏览+收藏) 与时间综合排序，真实分页
	type feedRow struct {
		domain.PhotoSet
		ViewCount    int `gorm:"column:view_count"`
		FavCount     int `gorm:"column:fav_count"`
		PhotoCountDB int `gorm:"column:photo_count_db"`
	}
	var rows []feedRow
	err := s.db.Model(&domain.PhotoSet{}).
		Select(`photosets.*,
			(SELECT COUNT(*) FROM view_histories vh WHERE vh.photoset_id = photosets.id AND vh.viewed_at >= ?) AS view_count,
			(SELECT COUNT(*) FROM favorites f WHERE f.photoset_id = photosets.id) AS fav_count,
			(SELECT COUNT(*) FROM photos p WHERE p.photoset_id = photosets.id) AS photo_count_db`, sevenDaysAgo).
		Where("photosets.status = ?", "published").
		Order("(view_count + fav_count * 3) DESC, photosets.created_at DESC").
		Offset(offset).Limit(pageSize).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]ExploreItem, len(rows))
	for i, r := range rows {
		source := "recommend"
		if r.CreatedAt.After(sevenDaysAgo) {
			source = "new"
		}
		if r.ViewCount+r.FavCount > 0 {
			source = "hot"
		}
		ps := r.PhotoSet
		ps.PhotoCount = r.PhotoCountDB
		items[i] = newExploreItem(ps, source, r.ViewCount, r.FavCount)
	}

	var total int64
	s.db.Model(&domain.PhotoSet{}).Where("status = ?", "published").Count(&total)

	return &ExploreFeedResponse{
		Items:   items,
		Page:    page,
		Total:   int(total),
		HasMore: int64(offset+len(items)) < total,
	}, nil
}

// getHotPhotosets gets hot photosets based on view count and favorite count
func (s *ExploreService) getHotPhotosets(limit int) ([]ExploreItem, error) {
	var photosets []domain.PhotoSet

	// 7天内的热门套图
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	
	err := s.db.Model(&domain.PhotoSet{}).
		Where("status = ? AND created_at >= ?", "published", sevenDaysAgo).
		Order("created_at DESC").
		Limit(limit).
		Find(&photosets).Error

	if err != nil {
		return nil, err
	}

	items := make([]ExploreItem, len(photosets))
	for i, ps := range photosets {
		items[i] = ExploreItem{
			ID:         ps.ID,
			Title:      ps.Title,
			CoverImage: ps.Cover,
			PhotoCount: ps.PhotoCount,
			ViewCount:  0,
			FavCount:   0,
			Price:      ps.Price,
			IsVip:      ps.IsFree == 0, // 0=付费=VIP
			Source:     "hot",
			CreatedAt:  ps.CreatedAt.Format(time.RFC3339),
		}
	}

	return items, nil
}

// getNewPhotosets gets newly created photosets
func (s *ExploreService) getNewPhotosets(limit int) ([]ExploreItem, error) {
	var photosets []domain.PhotoSet

	err := s.db.Model(&domain.PhotoSet{}).
		Where("status = ?", "published").
		Order("created_at DESC").
		Limit(limit).
		Find(&photosets).Error

	if err != nil {
		return nil, err
	}

	items := make([]ExploreItem, len(photosets))
	for i, ps := range photosets {
		items[i] = ExploreItem{
			ID:         ps.ID,
			Title:      ps.Title,
			CoverImage: ps.Cover,
			PhotoCount: ps.PhotoCount,
			ViewCount:  0,
			FavCount:   0,
			Price:      ps.Price,
			IsVip:      ps.IsFree == 0,
			Source:     "new",
			CreatedAt:  ps.CreatedAt.Format(time.RFC3339),
		}
	}

	return items, nil
}

// getRecommendedPhotosets gets personalized recommendations
func (s *ExploreService) getRecommendedPhotosets(userID uint, limit int) ([]ExploreItem, error) {
	// 简化推荐：获取用户浏览历史相关的套图
	var photosets []domain.PhotoSet

	if userID > 0 {
		// 基于用户浏览历史的简单推荐
		// 获取用户最近浏览的套图分类
		var recentCategories []string
		s.db.Model(&domain.ViewHistory{}).
			Joins("JOIN photosets ON view_histories.photoset_id = photosets.id").
			Where("view_histories.user_id = ?", userID).
			Distinct("photosets.category").
			Order("view_histories.viewed_at DESC").
			Limit(3).
			Pluck("photosets.category", &recentCategories)

		if len(recentCategories) > 0 {
			// 推荐同类别的其他套图
			err := s.db.Model(&domain.PhotoSet{}).
				Where("status = ? AND category IN ?", "published", recentCategories).
				Order("created_at DESC").
				Limit(limit).
				Find(&photosets).Error

			if err == nil && len(photosets) > 0 {
				items := make([]ExploreItem, len(photosets))
				for i, ps := range photosets {
					items[i] = ExploreItem{
						ID:         ps.ID,
						Title:      ps.Title,
						CoverImage: ps.Cover,
						PhotoCount: ps.PhotoCount,
						ViewCount:  0,
						FavCount:   0,
						Price:      ps.Price,
						IsVip:      ps.IsFree == 0,
						Source:     "recommend",
						CreatedAt:  ps.CreatedAt.Format(time.RFC3339),
					}
				}
				return items, nil
			}
		}
	}

	// 降级：返回随机套图
	err := s.db.Model(&domain.PhotoSet{}).
		Where("status = ?", "published").
		Order("RAND()").
		Limit(limit).
		Find(&photosets).Error

	if err != nil {
		return nil, err
	}

	items := make([]ExploreItem, len(photosets))
	for i, ps := range photosets {
		items[i] = ExploreItem{
			ID:         ps.ID,
			Title:      ps.Title,
			CoverImage: ps.Cover,
			PhotoCount: ps.PhotoCount,
			ViewCount:  0,
			FavCount:   0,
			Price:      ps.Price,
			IsVip:      ps.IsFree == 0,
			Source:     "recommend",
			CreatedAt:  ps.CreatedAt.Format(time.RFC3339),
		}
	}

	return items, nil
}

// GetHotPhotosets 获取热门套图排行榜（带时间衰减的热度算法，支持分页）
// 热度分 = (近7天浏览数×1 + 收藏数×5 + 1) / (发布至今小时数+2)^1.5
func (s *ExploreService) GetHotPhotosets(page, pageSize int) (*ExploreFeedResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 50 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	// 用子查询统计真实浏览/收藏数，按热度公式排序
	type hotRow struct {
		domain.PhotoSet
		ViewCount    int `gorm:"column:view_count"`
		FavCount     int `gorm:"column:fav_count"`
		PhotoCountDB int `gorm:"column:photo_count_db"`
	}
	var rows []hotRow
	err := s.db.Model(&domain.PhotoSet{}).
		Select(`photosets.*,
			(SELECT COUNT(*) FROM view_histories vh WHERE vh.photoset_id = photosets.id AND vh.viewed_at >= ?) AS view_count,
			(SELECT COUNT(*) FROM favorites f WHERE f.photoset_id = photosets.id) AS fav_count,
			(SELECT COUNT(*) FROM photos p WHERE p.photoset_id = photosets.id) AS photo_count_db`, sevenDaysAgo).
		Where("photosets.status = ?", "published").
		Order("((view_count * 1.0 + fav_count * 5.0 + 1) / POWER(TIMESTAMPDIFF(HOUR, photosets.created_at, NOW()) + 2, 1.5)) DESC").
		Offset(offset).Limit(pageSize).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]ExploreItem, len(rows))
	for i, r := range rows {
		ps := r.PhotoSet
		ps.PhotoCount = r.PhotoCountDB
		items[i] = newExploreItem(ps, "hot", r.ViewCount, r.FavCount)
	}

	var total int64
	s.db.Model(&domain.PhotoSet{}).Where("status = ?", "published").Count(&total)

	return &ExploreFeedResponse{
		Items:   items,
		Page:    page,
		Total:   int(total),
		HasMore: int64(offset+len(items)) < total,
	}, nil
}

// newExploreItem 统一构造 ExploreItem，保证各来源统计字段一致
func newExploreItem(ps domain.PhotoSet, source string, viewCount, favCount int) ExploreItem {
	return ExploreItem{
		ID:         ps.ID,
		Title:      ps.Title,
		CoverImage: ps.Cover,
		PhotoCount: ps.PhotoCount,
		ViewCount:  viewCount,
		FavCount:   favCount,
		Price:      ps.Price,
		IsVip:      ps.IsFree == 0,
		Category:   ps.Category,
		Source:     source,
		CreatedAt:  ps.CreatedAt.Format(time.RFC3339),
	}
}