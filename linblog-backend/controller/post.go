package controller

import (
	"encoding/json"
	"github.com/Linjiangzhu/blog-v2/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
)

func (ctrl *Controller) GetPost(c *gin.Context) {
	pidStr := strings.TrimSpace(c.Param("pid"))
	pid64, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"err": "resource does not exist"})
		return
	}
	pid := uint(pid64)

	role := "visitor"
	if val, exist := c.Get("role"); exist {
		role = val.(string)
	}
	var post *model.Post
	switch role {
	case "admin":
		post, err = ctrl.srv.GetPost(pid)
	default:
		post, err = ctrl.srv.GetVisiblePost(pid)
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"err": "resource does not exist"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post.ViewClass(role))
}

func (ctrl *Controller) GetPosts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "25")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	pageSize, err := strconv.ParseUint(pageSizeStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	role := "visitor"
	if ctxRole, exist := c.Get("role"); exist {
		if ctxRole.(string) == "admin" {
			role = "admin"
		}
	}
	var posts []model.Post
	switch role {
	case "admin":
		posts, err = ctrl.srv.GetPosts(uint(page), uint(pageSize))
	default:
		posts, err = ctrl.srv.GetVisiblePosts(uint(page), uint(pageSize))
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	var viewJSON []interface{}
	for _, p := range posts {
		viewP := p.ViewClass(role)
		viewJSON = append(viewJSON, viewP)
	}
	c.JSON(http.StatusOK, viewJSON)
}

func (ctrl *Controller) CreatePost(c *gin.Context) {
	byteStr, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	var reqPost model.Post
	err = json.Unmarshal(byteStr, &reqPost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	p, err := ctrl.srv.CreatePost(&reqPost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (ctrl *Controller) DeletePost(c *gin.Context) {
	pidStr := strings.TrimSpace(c.Param("pid"))
	pid64, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"err": "resource does not exist"})
		return
	}
	pid := uint(pid64)
	err = ctrl.srv.DeletePost(pid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"err": "resource does not exist"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
}

func (ctrl *Controller) UpdatePost(c *gin.Context) {
	pidStr := strings.TrimSpace(c.Param("pid"))
	pid64, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"err": "resource does not exist"})
		return
	}
	pid := uint(pid64)
	byteStr, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	var reqPost model.Post
	err = json.Unmarshal(byteStr, &reqPost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	reqPost.ID = pid
	p, err := ctrl.srv.UpdatePost(&reqPost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}
