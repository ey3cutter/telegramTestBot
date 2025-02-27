package db

import (
	_ "fmt"
	_ "log"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) GetCategories() ([]Category, error) {
	return s.repo.GetCategories()
}

func (s *Service) GetSubCategories(categoryID int) ([]SubCategory, error) {
	return s.repo.GetSubCategories(categoryID)
}

func (s *Service) GetProducts(subCategoryID int) ([]Product, error) {
	return s.repo.GetProducts(subCategoryID)
}
