package repository

import (
	"gorm.io/gorm"
	"photoset/internal/domain"
)

type PageRepository struct {
	db *gorm.DB
}

func NewPageRepository(db *gorm.DB) *PageRepository {
	return &PageRepository{db: db}
}

func (r *PageRepository) FindBySlug(slug string) (*domain.Page, error) {
	var page domain.Page
	if err := r.db.Where("slug = ? AND status = ?", slug, domain.PageStatusPublished).First(&page).Error; err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepository) GetAllPublished() ([]domain.Page, error) {
	var pages []domain.Page
	if err := r.db.Where("status = ?", domain.PageStatusPublished).Order("created_at desc").Find(&pages).Error; err != nil {
		return nil, err
	}
	return pages, nil
}

// Admin methods
func (r *PageRepository) List(page, pageSize int, keyword string) ([]domain.Page, int64, error) {
	var pages []domain.Page
	var total int64

	query := r.db.Model(&domain.Page{})
	if keyword != "" {
		query = query.Where("title LIKE ? OR slug LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// Count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginate
	offset := (page - 1) * pageSize
	if err := query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&pages).Error; err != nil {
		return nil, 0, err
	}

	return pages, total, nil
}

func (r *PageRepository) FindByID(id uint) (*domain.Page, error) {
	var page domain.Page
	if err := r.db.First(&page, id).Error; err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepository) Create(page *domain.Page) error {
	return r.db.Create(page).Error
}

func (r *PageRepository) Update(page *domain.Page) error {
	return r.db.Save(page).Error
}

func (r *PageRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Page{}, id).Error
}

// Markdown to HTML conversion (optional - currently just returns raw markdown)
// We can add goldmark later if needed.
func (r *PageRepository) RenderMarkdown(md string) (string, error) {
	// For now, just return raw markdown. Frontend will handle rendering.
	return md, nil
}

// AutoMigrate creates the table if not exists
func (r *PageRepository) AutoMigrate() error {
	return r.db.AutoMigrate(&domain.Page{})
}

// EnsureDefaultPages creates default pages if not exist
func (r *PageRepository) EnsureDefaultPages(userID uint) error {
	defaultPages := []domain.Page{
		{Slug: "about", Title: "关于我们", ContentMD: "# 关于我们\n\n我们是一个摄影社区...", UserID: userID},
		{Slug: "terms", Title: "使用协议", ContentMD: "# 使用协议\n\n请遵守平台规则...", UserID: userID},
		{Slug: "privacy", Title: "隐私政策", ContentMD: "# 隐私政策\n\n我们非常重视您的隐私...", UserID: userID},
		{Slug: "faq", Title: "常见问题", ContentMD: "# 常见问题\n\n### 如何上传作品？\n...", UserID: userID},
		{Slug: "contact", Title: "联系我们", ContentMD: "# 联系我们\n\n邮箱：support@photoset.io", UserID: userID},
	}

	for _, page := range defaultPages {
		var exists domain.Page
		if err := r.db.Where("slug = ?", page.Slug).First(&exists).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := r.db.Create(&page).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}