package service

import (
	"errors"
	"photoset/internal/domain"
	"photoset/internal/repository"

	"gorm.io/gorm"
)

type CollectionService struct {
	repo *repository.UserCollectionRepository
}

func NewCollectionService(repo *repository.UserCollectionRepository) *CollectionService {
	return &CollectionService{repo: repo}
}

// CreateCollection 创建合集
func (s *CollectionService) CreateCollection(userID uint, name, description string, isPublic bool) (*domain.UserCollection, error) {
	if name == "" {
		return nil, errors.New("合集名称不能为空")
	}
	if len(name) > 100 {
		return nil, errors.New("合集名称不能超过100个字符")
	}

	collection := &domain.UserCollection{
		UserID:      userID,
		Name:        name,
		Description: description,
		IsPublic:    isPublic,
	}

	if err := s.repo.Create(collection); err != nil {
		return nil, err
	}
	return collection, nil
}

// UpdateCollection 更新合集
func (s *CollectionService) UpdateCollection(userID, collectionID uint, name, description string, isPublic *bool) (*domain.UserCollection, error) {
	collection, err := s.repo.GetByID(collectionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("合集不存在")
		}
		return nil, err
	}
	if collection.UserID != userID {
		return nil, errors.New("无权修改此合集")
	}

	if name != "" {
		if len(name) > 100 {
			return nil, errors.New("合集名称不能超过100个字符")
		}
		collection.Name = name
	}
	if description != "" {
		collection.Description = description
	}
	if isPublic != nil {
		collection.IsPublic = *isPublic
	}

	if err := s.repo.Update(collection); err != nil {
		return nil, err
	}
	return collection, nil
}

// DeleteCollection 删除合集
func (s *CollectionService) DeleteCollection(userID, collectionID uint) error {
	collection, err := s.repo.GetByID(collectionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("合集不存在")
		}
		return err
	}
	if collection.UserID != userID {
		return errors.New("无权删除此合集")
	}
	return s.repo.Delete(userID, collectionID)
}

// GetCollection 获取合集详情
func (s *CollectionService) GetCollection(collectionID uint) (*domain.UserCollection, error) {
	collection, err := s.repo.GetByID(collectionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("合集不存在")
		}
		return nil, err
	}
	return collection, nil
}

// ListCollections 获取用户合集列表
func (s *CollectionService) ListCollections(userID uint, page, pageSize int) ([]domain.UserCollection, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.repo.ListByUser(userID, page, pageSize)
}

// AddItem 添加套图到合集
func (s *CollectionService) AddItem(userID, collectionID, photosetID uint) error {
	collection, err := s.repo.GetByID(collectionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("合集不存在")
		}
		return err
	}
	if collection.UserID != userID {
		return errors.New("无权修改此合集")
	}
	return s.repo.AddItem(collectionID, photosetID)
}

// RemoveItem 从合集移除套图
func (s *CollectionService) RemoveItem(userID, collectionID, photosetID uint) error {
	collection, err := s.repo.GetByID(collectionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("合集不存在")
		}
		return err
	}
	if collection.UserID != userID {
		return errors.New("无权修改此合集")
	}
	return s.repo.RemoveItem(collectionID, photosetID)
}

// BatchAddItems 批量添加套图到合集
func (s *CollectionService) BatchAddItems(userID, collectionID uint, photosetIDs []uint) error {
	if len(photosetIDs) > 50 {
		return errors.New("单次最多添加50个套图")
	}
	collection, err := s.repo.GetByID(collectionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("合集不存在")
		}
		return err
	}
	if collection.UserID != userID {
		return errors.New("无权修改此合集")
	}
	return s.repo.BatchAddItems(collectionID, photosetIDs)
}

// GetCollectionsContaining 获取包含指定套图的用户合集
func (s *CollectionService) GetCollectionsContaining(userID, photosetID uint) ([]domain.UserCollection, error) {
	return s.repo.GetUserCollectionsContaining(userID, photosetID)
}
