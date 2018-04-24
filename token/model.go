package token

import jwt "github.com/dgrijalva/jwt-go"

type UserMap struct {
	Uid      string `json:"uid"`
	Key      string `json:"_key"`
	Id       string `json:"_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Status   string `json:"status"`
	jwt.StandardClaims
}
