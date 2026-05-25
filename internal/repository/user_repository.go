package repository

import (
	"errors"
	"photoset/internal/database"
	"photoset/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	// Admin APIs
	List(page, pageSize int, role, keyword string, status int) ([]domain.User, int64, error)
	FindByIDWithStats(id uint) (*domain.User, int64, int64, float64, int64, error)
	UpdateStatus(id uint, status int) error
	UpdateRole(id uint, role string) error
	CountAll() (int64, error)
	// Follow APIs
	UpdateField(id uint, field string, value interface{}) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.MySQL,
	}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// ============ Admin APIs ============

// List 分页查询用户列表（管理员用）
func (r *userRepository) List(page, pageSize int, role, keyword string, status int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	query := r.db.Model(&domain.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != -1 { // status 为 -1 时不筛选
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("nickname LIKE ? OR email LIKE ?", like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Select("id, nickname, email, role, status, created_at, last_login_at, membership_expires").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&users).Error

	return users, total, err
}

// FindByIDWithStats 获取用户详情及统计（含收藏数）
func (r *userRepository) FindByIDWithStats(id uint) (*domain.User, int64, int64, float64, int64, error) {
	var user domain.User
	err := r.db.Select("id, nickname, email, role, status, created_at, last_login_at, membership_expires").
		Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	var photoSetCount int64
	r.db.Model(&domain.PhotoSet{}).Where("user_id = ?", id).Count(&photoSetCount)

	var orderCount int64
	var totalSpent float64
	r.db.Model(&domain.Order{}).Where("user_id = ? AND status = ?", id, "paid").Count(&orderCount)
	r.db.Model(&domain.Order{}).Where("user_id = ? AND status = ?", id, "paid").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalSpent)

	var favoriteCount int64
	r.db.Model(&domain.Favorite{}).Where("user_id = ?", id).Count(&favoriteCount)

	return &user, photoSetCount, orderCount, totalSpent, favoriteCount, nil
}

// UpdateStatus 更新用户状态（封号/解封）
func (r *userRepository) UpdateStatus(id uint, status int) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateRole 更新用户角色
func (r *userRepository) UpdateRole(id uint, role string) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update("role", role).Error
}

// CountAll 统计全部用户数量
func (r *userRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Count(&count).Error
	return count, err
}

// UpdateField 更新用户指定字段
func (r *userRepository) UpdateField(id uint, field string, value interface{}) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update(field, value).Error
}

