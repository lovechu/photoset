package service

import (
	"photoset/internal/domain"
	"photoset/internal/repository"
)

type PageService struct {
	repo *repository.PageRepository
}

func NewPageService(repo *repository.PageRepository) *PageService {
	return &PageService{repo: repo}
}

func (s *PageService) GetPublishedPage(slug string) (*domain.Page, error) {
	page, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (s *PageService) ListPublishedPages() ([]domain.Page, error) {
	return s.repo.GetAllPublished()
}

// Admin methods
func (s *PageService) AdminListPages(page, pageSize int, keyword string) ([]domain.Page, int64, error) {
	return s.repo.List(page, pageSize, keyword)
}

func (s *PageService) GetPageByID(id uint) (*domain.Page, error) {
	return s.repo.FindByID(id)
}

func (s *PageService) CreatePage(page *domain.Page) error {
	return s.repo.Create(page)
}

func (s *PageService) UpdatePage(page *domain.Page) error {
	return s.repo.Update(page)
}

func (s *PageService) DeletePage(id uint) error {
	return s.repo.Delete(id)
}