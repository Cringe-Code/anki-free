package anki

import (
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Name  string `json:"name"`
	Login string `json:"login"`
	Id    int64  `json:"id"`
}
type UserReq struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserRes struct {
	Name string `json:"name"`
}

type Pack struct {
	Name string `json:"name"`
	Rank int64  `json:"rank"`
}

type TokenClaims struct {
	jwt.StandardClaims
	Login string `json:"login"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
