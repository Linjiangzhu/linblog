package controller

import (
	"github.com/Linjiangzhu/linblog/linblog-backend/service"
)

type Controller struct {
	srv *service.Service
}

func NewController(srv *service.Service) *Controller {
	return &Controller{srv: srv}
}
