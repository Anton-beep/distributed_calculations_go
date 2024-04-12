package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"storage/internal/cryptPasswords"
	"storage/internal/db"
	"time"
)

type InRegister struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type OutRegister struct {
	Access  string `json:"access"`
	Message string `json:"message"`
}

func (a *API) Register(c *gin.Context) {
	var in InRegister
	var out OutRegister
	if err := c.ShouldBindBodyWith(&in, binding.JSON); err != nil {
		out.Message = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// check if user already exists
	_, err := a.db.GetUserByUsername(in.Login)
	if err == nil {
		out.Message = "user already exists"
		c.JSON(http.StatusConflict, out)
		return
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": in.Login,
		"nbf":  now.Unix(),
		"exp":  now.Add(5 * time.Minute).Unix(),
		"iat":  now.Unix(),
	})

	tokenString, err := token.SignedString(a.secretSignature)
	if err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusInternalServerError, out)
	}

	// add user to db
	hash, err := cryptPasswords.GeneratePasswordHash(in.Password)
	if err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusInternalServerError, out)
		return
	}
	user := db.User{
		Login:    in.Login,
		Password: hash,
	}
	if err := a.db.AddUser(user); err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusInternalServerError, out)
		return
	}

	out.Access = tokenString
	out.Message = "ok"
	c.JSON(http.StatusOK, out)
}

type InLogin struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type OutLogin struct {
	Access  string `json:"access"`
	Message string `json:"message"`
}

func (a *API) Login(c *gin.Context) {
	var in InLogin
	var out OutLogin
	if err := c.ShouldBindBodyWith(&in, binding.JSON); err != nil {
		out.Message = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// check if user exists
	user, err := a.db.GetUserByUsername(in.Login)
	if err != nil {
		out.Message = "user not found"
		c.JSON(http.StatusNotFound, out)
		return
	}

	err = cryptPasswords.ComparePasswordWithHash(user.Password, in.Password)
	if err != nil {
		out.Message = "wrong password"
		c.JSON(http.StatusUnauthorized, out)
		return
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": in.Login,
		"nbf":  now.Unix(),
		"exp":  now.Add(5 * time.Minute).Unix(),
		"iat":  now.Unix(),
	})

	tokenString, err := token.SignedString(a.secretSignature)
	if err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusInternalServerError, out)
	}

	out.Access = tokenString
	out.Message = "ok"
	c.JSON(http.StatusOK, out)
}
