package service

import (
	"time"

	"photoset/internal/domain"

	"gorm.io/gorm"
)

// CreatorStatsService 创作者数据统计服务
type CreatorStatsService struct {
	db *gorm.DB
}

// NewCreatorStatsService 创建创作者数据统计服务
func NewCreatorStatsService(db *gorm.DB) *CreatorStatsService {
	return &CreatorStatsService{db: db}
}

// GetCreatorStats 获取创作者统计数据
func (s *CreatorStatsService) GetCreatorStats(userID uint) (*domain.CreatorStats, error) {
	stats := &domain.CreatorStats{}

	// 获取总套图数
	s.db.Model(&domain.PhotoSet{}).
		Where("user_id = ? AND status = ?", userID, "published").
		Count(&stats.TotalPhotosets)

	// 获取总浏览量、下载量、收藏量
	s.db.Model(&domain.PhotoSet{}).
		Where("user_id = ? AND status = ?", userID, "published").
		Select("COALESCE(SUM(view_count), 0), COALESCE(SUM(download_count), 0), COALESCE(SUM(fav_count), 0)").
		Row().
		Scan(&stats.TotalViews, &stats.TotalDownloads, &stats.TotalFavorites)

	// 获取总收益（已支付的订单）
	s.db.Model(&domain.Order{}).
		Joins("JOIN photosets ON orders.photoset_id = photosets.id").
		Where("photosets.user_id = ? AND orders.status = ?", userID, "paid").
		Select("COALESCE(SUM(orders.amount), 0)").
		Row().
		Scan(&stats.TotalRevenue)

	// 获取总订单数
	s.db.Model(&domain.Order{}).
		Joins("JOIN photosets ON orders.photoset_id = photosets.id").
		Where("photosets.user_id = ? AND orders.status = ?", userID, "paid").
		Count(&stats.TotalOrders)

	// 获取今日数据
	today := time.Now().Format("2006-01-02")
	todayStart, _ := time.Parse("2006-01-02", today)

	s.db.Model(&domain.ViewHistory{}).
		Where("user_id IN (?) AND viewed_at >= ?",
			s.db.Model(&domain.PhotoSet{}).Select("id").Where("user_id = ?", userID),
			todayStart,
		).
		Count(&stats.TodayViews)

	// 获取本周数据
	weekStart := todayStart.AddDate(0, 0, -7)

	s.db.Model(&domain.ViewHistory{}).
		Where("user_id IN (?) AND viewed_at >= ?",
			s.db.Model(&domain.PhotoSet{}).Select("id").Where("user_id = ?", userID),
			weekStart,
		).
		Count(&stats.WeekViews)

	return stats, nil
}

// GetDailyStats 获取每日统计数据
func (s *CreatorStatsService) GetDailyStats(userID uint, days int) ([]domain.DailyStats, error) {
	var stats []domain.DailyStats

	// 获取最近N天的日期列表
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// 查询套图ID列表
	var photosetIDs []uint
	s.db.Model(&domain.PhotoSet{}).
		Where("user_id = ? AND status = ?", userID, "published").
		Pluck("id", &photosetIDs)

	if len(photosetIDs) == 0 {
		return stats, nil
	}

	// 按日统计浏览量
	rows, err := s.db.Raw(`
		SELECT DATE(viewed_at) as date, COUNT(*) as views
		FROM view_histories
		WHERE photoset_id IN ? AND viewed_at >= ? AND viewed_at < ?
		GROUP BY DATE(viewed_at)
		ORDER BY date ASC
	`, photosetIDs, startDate, endDate).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statsMap := make(map[string]*domain.DailyStats)
	for rows.Next() {
		var ds domain.DailyStats
		rows.Scan(&ds.Date, &ds.Views)
		statsMap[ds.Date] = &ds
	}

	// 查询每日订单收益
	var orderStats []struct {
		Date    string
		Revenue float64
	}
	s.db.Model(&domain.Order{}).
		Where("photoset_id IN ? AND status = ? AND created_at >= ? AND created_at < ?",
			photosetIDs, "paid", startDate, endDate).
		Select("DATE(created_at) as date, SUM(amount) as revenue").
		Group("DATE(created_at)").
		Scan(&orderStats)

	for _, os := range orderStats {
		if existing, ok := statsMap[os.Date]; ok {
			existing.Revenue = os.Revenue
		} else {
			statsMap[os.Date] = &domain.DailyStats{
				Date:    os.Date,
				Revenue: os.Revenue,
			}
		}
	}

	// 按日期排序输出
	for i := days - 1; i >= 0; i-- {
		date := endDate.AddDate(0, 0, -i).Format("2006-01-02")
		if ds, ok := statsMap[date]; ok {
			stats = append(stats, *ds)
		} else {
			stats = append(stats, domain.DailyStats{Date: date})
		}
	}

	return stats, nil
}

// GetPhotoSetStats 获取各套图统计数据
func (s *CreatorStatsService) GetPhotoSetStats(userID uint, limit int) ([]domain.PhotoSetStats, error) {
	var stats []domain.PhotoSetStats

	s.db.Model(&domain.PhotoSet{}).
		Where("user_id = ? AND status = ?", userID, "published").
		Select("id as photoset_id, title, cover_image, view_count, download_count, fav_count, created_at").
		Order("view_count DESC").
		Limit(limit).
		Scan(&stats)

	// 获取每个套图的收益
	for i := range stats {
		s.db.Model(&domain.Order{}).
			Where("photoset_id = ? AND status = ?", stats[i].PhotoSetID, "paid").
			Select("COALESCE(SUM(amount), 0), COUNT(*)").
			Row().
			Scan(&stats[i].Revenue, &stats[i].OrderCount)
	}

	return stats, nil
}
