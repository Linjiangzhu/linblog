package service

import (
	"github.com/Linjiangzhu/blog-v2/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}