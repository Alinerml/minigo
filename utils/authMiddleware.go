package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"minigo/conf"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token") //获取body中的数据
		if tokenString == "" {
			tokenString = c.PostForm("token") //获取表单中的数据
		}
		if tokenString == "" {
			c.JSON(http.StatusOK, gin.H{
				"status_code": -1,
				"status_msg":  "Unauthorized",
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { //看解析对不对
			return conf.SecretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusOK, gin.H{
				"status_code": -1,
				"status_msg":  "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
