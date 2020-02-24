package jwt

import (
	"errors"
	"go_gin_base/pkg/e"
	"go_gin_base/pkg/util"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		code = e.SUCCESS
		token, err := getHeaderAuth(c)
		if err != nil {
			code = e.ERROR_AUTH_NOT_FOUND
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}
		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func getHeaderAuth(c *gin.Context) (string, error) {
	headerStr := c.GetHeader("Authorization")
	if headerStr == "" {
		return "", errors.New("can not find token")
	}
	headerParts := strings.Split(headerStr, " ")
	if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}
	return headerParts[1], nil
}
