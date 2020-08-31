package service

import "github.com/Linjiangzhu/linblog/linblog-backend/model"

func (s *Service) GetUserByUsername(username string) (*model.User, error) {
	return s.repo.GetUserByUsername(username)
}
