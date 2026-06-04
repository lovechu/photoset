package domain

import "time"

// CreatorStats 创作者数据统计
type CreatorStats struct {
	TotalPhotosets  int64   `json:"total_photosets"`
	TotalViews      int64   `json:"total_views"`
	TotalDownloads  int64   `json:"total_downloads"`
	TotalFavorites  int64   `json:"total_favorites"`
	TotalRevenue    float64 `json:"total_revenue"`
	TotalOrders     int64   `json:"total_orders"`
	TodayViews      int64   `json:"today_views"`
	TodayDownloads  int64   `json:"today_downloads"`
	TodayFavorites  int64   `json:"today_favorites"`
	TodayRevenue    float64 `json:"today_revenue"`
	WeekViews       int64   `json:"week_views"`
	WeekDownloads   int64   `json:"week_downloads"`
	WeekFavorites   int64   `json:"week_favorites"`
	WeekRevenue     float64 `json:"week_revenue"`
}

// DailyStats 每日统计
type DailyStats struct {
	Date       string  `json:"date"`
	Views      int     `json:"views"`
	Downloads  int     `json:"downloads"`
	Favorites  int     `json:"favorites"`
	Revenue    float64 `json:"revenue"`
}

// PhotoSetStats 单个套图统计
type PhotoSetStats struct {
	PhotoSetID   uint    `json:"photoset_id"`
	Title        string  `json:"title"`
	CoverImage   string  `json:"cover_image"`
	ViewCount    int     `json:"view_count"`
	DownloadCount int    `json:"download_count"`
	FavCount     int     `json:"fav_count"`
	Revenue      float64 `json:"revenue"`
	OrderCount   int     `json:"order_count"`
	CreatedAt    time.Time `json:"created_at"`
}
