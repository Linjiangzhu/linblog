package service

import (
	"github.com/Linjiangzhu/linblog/linblog-backend/model"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (s *Service) VerifyUserPassword(logReq *model.LoginRequestEntity) (token string, err error) {
	u, err := s.GetUserByUsername(logReq.Username)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(logReq.Password))
	if err != nil {
		return "", err
	}
	jwt, err := s.GenerateJWT(u)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func (s *Service) GenerateJWT(u *model.User) (string, error) {
	claims := model.CustomClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
		UID:    u.ID,
		RoleID: u.RoleID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := token.SignedString([]byte("jwt"))
	if err != nil {
		return "", err
	}
	return signedStr, nil
}

func (s *Service) BlockJWT(token string, expire time.Time) error {
	return s.repo.CountDownKey(token, "none", expire)
}
