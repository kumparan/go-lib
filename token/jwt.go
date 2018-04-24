package token

import (
	"strings"

	"encoding/base64"
	"encoding/json"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context, jwtsecret string) (*UserMap, error) {
	var user *UserMap

	Authorization := c.Request.Header["Authorization"]
	if len(Authorization) != 0 {
		jwt := strings.Split(Authorization[0], " ")
		if jwt != nil && jwt[0] == "Bearer" && len(jwt) == 2 {
			token := jwt[1]
			user, _ = JwtParse(token,jwtsecret)
		}
	}

	Cookie, err := c.Request.Cookie("token")
	if err == nil {
		tokenBase64 := Cookie.Value
		tokenJson, _ := base64.StdEncoding.DecodeString(tokenBase64)
		var token map[string]string
		json.Unmarshal([]byte(tokenJson), &token)

		user, _ = JwtParse(token["token"],jwtsecret)
	}

	return user, err
}

func JwtParse(myToken string, jwtsecret string) (*UserMap, error) {
	token, err := jwt.ParseWithClaims(myToken, &UserMap{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtsecret), nil
	})

	claims := token.Claims.(*UserMap)
	claims.Role = getRole(claims.Level)

	return claims, err
}

