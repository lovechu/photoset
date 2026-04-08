package service

import (
	"errors"
	"photoset/internal/domain"
	"photoset/internal/repository"

	"gorm.io/gorm"
)

type PhotoSetService struct {
	repo *repository.PhotoSetRepository
}

func NewPhotoSetService(repo *repository.PhotoSetRepository) *PhotoSetService {
	return &PhotoSetService{repo: repo}
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

	return nil
}

// GetPhotoSetList 获取套图列表
func (s *PhotoSetService) GetPhotoSetList(page, pageSize int, tag string, keyword string, userID uint, onlyMine bool) ([]domain.PhotoSet, int64, error) {
	return s.repo.List(page, pageSize, tag, keyword, userID, onlyMine)
}

// GetPhotoSetDetail 获取套图详情（带完整图片列表）
func (s *PhotoSetService) GetPhotoSetDetail(id uint) (*domain.PhotoSet, error) {
	return s.repo.FindByID(id)
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

	return false
}

// GetAllTags 获取所有标签
func (s *PhotoSetService) GetAllTags() ([]domain.Tag, error) {
	return s.repo.ListTags()
}
