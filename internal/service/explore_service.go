package service

import (
	"math/rand"
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

// GetExploreFeed gets mixed explore feed
func (s *ExploreService) GetExploreFeed(userID uint, page, pageSize int) (*ExploreFeedResponse, error) {
	// 计算每种类型的数量
	hotCount := pageSize / 3
	newCount := pageSize / 3
	recommendCount := pageSize - hotCount - newCount

	// 并发获取各类内容
	type result struct {
		items []ExploreItem
		err   error
	}

	hotCh := make(chan result, 1)
	newCh := make(chan result, 1)
	recommendCh := make(chan result, 1)

	// 获取热门套图
	go func() {
		items, err := s.getHotPhotosets(hotCount)
		hotCh <- result{items, err}
	}()

	// 获取新上架套图
	go func() {
		items, err := s.getNewPhotosets(newCount)
		newCh <- result{items, err}
	}()

	// 获取推荐套图
	go func() {
		items, err := s.getRecommendedPhotosets(userID, recommendCount)
		recommendCh <- result{items, err}
	}()

	// 收集结果
	var allItems []ExploreItem

	hotResult := <-hotCh
	if hotResult.err == nil {
		allItems = append(allItems, hotResult.items...)
	}

	newResult := <-newCh
	if newResult.err == nil {
		allItems = append(allItems, newResult.items...)
	}

	recommendResult := <-recommendCh
	if recommendResult.err == nil {
		allItems = append(allItems, recommendResult.items...)
	}

	// 打乱顺序，避免每次都是一样的排列
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allItems), func(i, j int) {
		allItems[i], allItems[j] = allItems[j], allItems[i]
	})

	// 计算总数（简化处理，使用估算）
	total := len(allItems) * 10 // 假设还有更多内容

	return &ExploreFeedResponse{
		Items:   allItems,
		Page:    page,
		Total:   total,
		HasMore: page < 10, // 限制最多10页
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

			if err == nil && len(photosets) >= limit {
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