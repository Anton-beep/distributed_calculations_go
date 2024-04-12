package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
)

type InAuthData struct {
	Access string `json:"access" binding:"required"`
}

type OutAuthData struct {
	Message string `json:"message"`
}

func (a *API) Auth(c *gin.Context) {
	var in InAuthData
	var out OutAuthData
	if err := c.ShouldBindBodyWith(&in, binding.JSON); err != nil {
		out.Message = err.Error()
		c.JSON(http.StatusUnauthorized, out)
		c.Abort()
		return
	}

	tokenFromString, err := jwt.Parse(in.Access, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			panic(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}

		return a.secretSignature, nil
	})

	if err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusInternalServerError, out)
		c.Abort()
		return
	}

	claims, ok := tokenFromString.Claims.(jwt.MapClaims)

	if !ok {
		out.Message = "looks like wrong token"
		zap.S().Error(out)
		c.JSON(http.StatusBadRequest, out)
		c.Abort()
		return
	}

	// find user in db
	user, err := a.db.GetUserByUsername(claims["name"].(string))
	if err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusInternalServerError, out)
		c.Abort()
		return
	}
	c.Set("user", user)

	c.Next()
}
