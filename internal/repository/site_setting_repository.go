package repository

import (
	"fmt"
	"photoset/internal/database"
	"photoset/internal/domain"
)

type SiteSettingRepository struct{}

func NewSiteSettingRepository() *SiteSettingRepository {
	return &SiteSettingRepository{}
}

// GetAll 返回所有配置，以 key->value 的 map 返回
func (r *SiteSettingRepository) GetAll() (map[string]string, error) {
	db := database.GetMySQL()
	var settings []domain.SiteSetting
	if err := db.Find(&settings).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}

// BatchUpsert 批量插入或更新（有则更新，无则插入）
// 使用原生 SQL 的 INSERT ... ON DUPLICATE KEY UPDATE 避免 key 保留字问题
func (r *SiteSettingRepository) BatchUpsert(data map[string]interface{}) error {
	db := database.GetMySQL()
	for k, v := range data {
		// 支持多种类型转换为字符串存储
		var val string
		switch vv := v.(type) {
		case string:
			val = vv
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			val = fmt.Sprintf("%d", vv)
		case float32, float64:
			val = fmt.Sprintf("%f", vv)
		case bool:
			if vv {
				val = "1"
			} else {
				val = "0"
			}
		default:
			val = fmt.Sprint(vv)
		}
		// 使用原生 SQL INSERT ... ON DUPLICATE KEY UPDATE
		// 避免 `key` 作为 MySQL 保留字导致的 SQL 语法错误
		sql := "INSERT INTO site_settings (`key`, `value`, created_at, updated_at) VALUES (?, ?, NOW(), NOW()) ON DUPLICATE KEY UPDATE `value` = VALUES(`value`), updated_at = NOW()"
		if err := db.Exec(sql, k, val).Error; err != nil {
			return err
		}
	}
	return nil
}
