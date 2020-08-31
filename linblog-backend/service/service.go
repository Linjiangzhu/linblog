package service

import (
	"github.com/Linjiangzhu/linblog/linblog-backend/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}
