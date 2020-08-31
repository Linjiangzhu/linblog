package controller

import (
	"github.com/Linjiangzhu/blog-v2/service"
)

type Controller struct {
	srv *service.Service
}

func NewController(srv *service.Service) *Controller {
	return &Controller{srv: srv}
}
