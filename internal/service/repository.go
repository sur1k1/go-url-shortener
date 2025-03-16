package service

import "github.com/sur1k1/go-url-shortener/internal/models"

type ServiceRepository interface {
	GetURL(shortURL string) (*models.URLData, error)
	SaveURL(urlData *models.URLData) error
}

type Service struct {
	repo ServiceRepository
}

func New(repo ServiceRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetURL(shortURL string) (*models.URLData, error) {
	return s.repo.GetURL(shortURL)
}

func (s *Service) SaveURL(urlData *models.URLData) error {
	return s.repo.SaveURL(urlData)
}