package service

import "github.com/Linjiangzhu/blog-v2/model"

func (s *Service) GetUserByUsername(username string) (*model.User, error) {
	return s.repo.GetUserByUsername(username)
}
