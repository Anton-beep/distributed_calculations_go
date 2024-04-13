package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type OutAuthData struct {
	Message string `json:"message"`
}

func (a *API) Auth(c *gin.Context) {
	var out OutAuthData
	access := c.GetHeader("Authorization")
	access = strings.Replace(access, "Bearer ", "", 1)

	tokenFromString, err := jwt.Parse(access, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			panic(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}

		return a.secretSignature, nil
	})

	if err != nil {
		out.Message = err.Error()
		c.JSON(http.StatusUnauthorized, out)
		c.Abort()
		return
	}

	claims, ok := tokenFromString.Claims.(jwt.MapClaims)

	if !ok {
		out.Message = "looks like wrong token"
		zap.S().Error(out)
		c.JSON(http.StatusUnauthorized, out)
		c.Abort()
		return
	}

	// find user in db
	user, err := a.db.GetUserByUsername(claims["name"].(string))
	if err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusUnauthorized, out)
		c.Abort()
		return
	}
	c.Set("user", user)

	c.Next()
}
