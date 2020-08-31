package service

import (
	"errors"
	"fmt"
	"github.com/Linjiangzhu/blog-v2/model"
)

func (s *Service) GetPost(pid uint) (*model.Post, error) {
	return s.repo.GetPost(pid)
}

func (s *Service) GetVisiblePost(pid uint) (*model.Post, error) {
	return s.repo.GetVisiblePost(pid)
}

func (s *Service) GetPosts(page, pageSize uint) ([]model.Post, error) {
	return s.repo.GetPosts((page-1)*pageSize, pageSize)
}

func (s *Service) GetVisiblePosts(page, pageSize uint) ([]model.Post, error) {
	return s.repo.GetVisiblePosts((page-1)*pageSize, pageSize)
}

func (s *Service) CreatePost(p *model.Post) (*model.Post, error) {
	p, errs := s.repo.CreatePost(p)
	errStr := ""
	if len(errs) > 0 {
		for _, err := range errs {
			errStr += fmt.Sprintf("%s\n", err.Error())
		}
		return nil, errors.New(errStr)
	}
	return p, nil
}

func (s *Service) DeletePost(pid uint) error {
	errs := s.repo.DeletePost(pid)
	errStr := ""
	if len(errs) > 0 {
		for _, err := range errs {
			errStr += fmt.Sprintf("%s\n", err.Error())
		}
		return errors.New(errStr)
	}
	return nil
}

func (s *Service) UpdatePost(p *model.Post) (*model.Post, error) {
	p, errs := s.repo.UpdatePost(p)
	errStr := ""
	if len(errs) > 0 {
		for _, err := range errs {
			errStr += fmt.Sprintf("%s\n", err.Error())
		}
		return nil, errors.New(errStr)
	}
	return p, nil
}
