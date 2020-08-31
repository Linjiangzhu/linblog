package model

import "github.com/dgrijalva/jwt-go"

type CustomClaim struct {
	jwt.StandardClaims
	UID    string
	RoleID uint
}
