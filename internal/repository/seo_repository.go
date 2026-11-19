package repository

import (
	"time"

	"photoset/internal/database"
	"photoset/internal/domain"
)

// SeoPhotoset sitemap 所需的套图精简信息
type SeoPhotoset struct {
	ID        uint
	UpdatedAt time.Time
	Cover     string
	ImageURLs []string
}

// SeoPost sitemap 所需的帖子精简信息
type SeoPost struct {
	ID        uint
	UpdatedAt time.Time
}

// SeoRepository sitemap 数据查询
type SeoRepository struct{}

func NewSeoRepository() *SeoRepository { return &SeoRepository{} }

// ListPublishedPhotosets 查询所有已发布、未删除的套图及其图片 URL
func (r *SeoRepository) ListPublishedPhotosets() ([]SeoPhotoset, error) {
	db := database.GetMySQL()
	var rows []struct {
		ID        uint
		UpdatedAt time.Time
		Cover     string
	}
	if err := db.Model(&domain.PhotoSet{}).
		Select("id", "updated_at", "cover").
		Where("status = ? AND deleted_at IS NULL", "published").
		Order("id ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]SeoPhotoset, 0, len(rows))
	for _, row := range rows {
		sp := SeoPhotoset{ID: row.ID, UpdatedAt: row.UpdatedAt, Cover: row.Cover}
		if row.Cover != "" {
			sp.ImageURLs = append(sp.ImageURLs, row.Cover)
		}
		result = append(result, sp)
	}

	// 批量查询每个套图的图片 URL（最多取每套图前 1000 张，符合图片 sitemap 限制）
	if len(result) > 0 {
		ids := make([]uint, 0, len(result))
		for _, sp := range result {
			ids = append(ids, sp.ID)
		}
		var photos []struct {
			PhotoSetID uint
			URL        string
		}
		// 使用窗口函数限制每套图图片数（MySQL 8.0 支持）
		err := db.Raw(`
			SELECT photoset_id, url FROM (
				SELECT photoset_id, url,
				       ROW_NUMBER() OVER (PARTITION BY photoset_id ORDER BY id ASC) AS rn
				FROM photos
				WHERE photoset_id IN ? AND deleted_at IS NULL
			) t WHERE rn <= 1000
		`, ids).Scan(&photos).Error
		if err == nil {
			idx := make(map[uint]*SeoPhotoset)
			for i := range result {
				idx[result[i].ID] = &result[i]
			}
			seen := make(map[uint]map[string]bool)
			for _, p := range photos {
				sp := idx[p.PhotoSetID]
				if sp == nil || p.URL == "" {
					continue
				}
				if seen[sp.ID] == nil {
					seen[sp.ID] = make(map[string]bool)
				}
				if seen[sp.ID][p.URL] {
					continue
				}
				seen[sp.ID][p.URL] = true
				sp.ImageURLs = append(sp.ImageURLs, p.URL)
				if len(sp.ImageURLs) > 1000 {
					break
				}
			}
		}
		// photos 表可能没有 deleted_at 列：若窗口查询失败，回退到简单查询
		if err != nil {
			var photos2 []struct {
				PhotoSetID uint
				URL        string
			}
			if err2 := db.Model(&domain.Photo{}).
				Select("photoset_id", "url").
				Where("photoset_id IN ?", ids).
				Order("id ASC").
				Limit(100000).
				Scan(&photos2).Error; err2 == nil {
				idx := make(map[uint]*SeoPhotoset)
				for i := range result {
					idx[result[i].ID] = &result[i]
				}
				counts := make(map[uint]int)
				for _, p := range photos2 {
					sp := idx[p.PhotoSetID]
					if sp == nil || p.URL == "" || counts[sp.ID] >= 1000 {
						continue
					}
					sp.ImageURLs = append(sp.ImageURLs, p.URL)
					counts[sp.ID]++
				}
			}
		}
	}
	return result, nil
}

// ListPublicPosts 查询所有已审核、公开的帖子
func (r *SeoRepository) ListPublicPosts() ([]SeoPost, error) {
	db := database.GetMySQL()
	var rows []SeoPost
	if err := db.Model(&domain.Post{}).
		Select("id", "updated_at").
		Where("status = ? AND visibility = ? AND deleted_at IS NULL", "approved", "public").
		Order("id ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// ListCategories 查询所有有已发布套图的分类 slug
func (r *SeoRepository) ListCategories() ([]string, error) {
	db := database.GetMySQL()
	var cats []string
	if err := db.Model(&domain.PhotoSet{}).
		Distinct("category").
		Where("status = ? AND deleted_at IS NULL AND category != ?", "published", "").
		Pluck("category", &cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}
