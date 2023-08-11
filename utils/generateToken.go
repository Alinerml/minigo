package utils

import (
	"github.com/RaymondCode/simple-demo/conf"
	"github.com/golang-jwt/jwt"
	"time"
)

func GenerateJWTToken(userid int64) (string, error) {

	claims := jwt.MapClaims{
		"sub": userid,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(conf.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
