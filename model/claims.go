package model

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UID      string `json:"uID"`
	Username string `json:"username"`
	jwt.StandardClaims
}
