package middleware

import (
	"github.com/Linjiangzhu/linblog/linblog-backend/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func NewJWTAuth(db *gorm.DB, rc *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get authorization header from context
		authHeader := c.GetHeader("Authorization")
		splited := strings.Split(authHeader, "Bearer ")
		// check if have bearer fragment and token, if not, return with 401
		if len(splited) <= 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "missing authorization header"})
			return
		}
		tokenStr := splited[1]
		// check if token string exist in block list, if so, return with 401
		_, err := rc.Get(tokenStr).Result()
		if err == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "token expired"})
			return
		}
		// parse token, if not valid, return with 401
		token, err := jwt.ParseWithClaims(tokenStr, &model.CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("jwt"), nil
		})
		if err != nil {
			valError, _ := err.(*jwt.ValidationError)
			if valError.Errors == jwt.ValidationErrorExpired {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "token expired"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}
		if claims, ok := token.Claims.(*model.CustomClaim); ok && token.Valid {
			c.Set("access_token", tokenStr)
			c.Set("uid", claims.UID)
			c.Set("expiration", claims.ExpiresAt)
			role := model.Role{ID: claims.RoleID}
			err := db.First(&role).Error
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			}
			c.Set("role", role.Name)
			c.Next()
		}
	}
}
