package controller

import (
	"encoding/json"
	"github.com/Linjiangzhu/blog-v2/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (ctrl *Controller) Login(c *gin.Context) {
	postBody, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	// scan to struct
	var loginReq model.LoginRequestEntity
	err = json.Unmarshal(postBody, &loginReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	// verify username and password
	token, err := ctrl.srv.VerifyUserPassword(&loginReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.LoginResponseEntity{AccessToken: token})

}

func (ctrl *Controller) Logout(c *gin.Context) {
	var token string
	var unixExpire int64
	if accessToken, exist := c.Get("access_token"); exist {
		token, _ = accessToken.(string)
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "missing access token in context"})
		return
	}
	if expireTime, exist := c.Get("expiration"); exist {
		unixExpire, _ = expireTime.(int64)
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "missing access token in context"})
		return
	}
	expire := time.Unix(unixExpire, 0)
	err := ctrl.srv.BlockJWT(token, expire)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"msg": "logout successful"})

}
