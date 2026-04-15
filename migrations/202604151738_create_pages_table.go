package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	migrations = append(migrations, &gormigrate.Migration{
		ID: "202604151738_create_pages_table",
		Migrate: func(tx *gorm.DB) error {
			type Page struct {
				ID        uint   `gorm:"primarykey"`
				Slug      string `gorm:"size:100;uniqueIndex;not null;comment:路径标识"`
				Title     string `gorm:"size:200;not null;comment:页面标题"`
				ContentMD string `gorm:"type:text;comment:Markdown内容"`
				UserID    uint   `gorm:"not null;index;comment:创建者"`
				Status    string `gorm:"size:20;default:'published';index;comment:状态: published/draft"`
				CreatedAt int64  `gorm:"autoCreateTime"`
				UpdatedAt int64  `gorm:"autoUpdateTime"`
			}

			return tx.AutoMigrate(&Page{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("pages")
		},
	})
}