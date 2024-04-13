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

type OutGetUser struct {
	Login string `json:"login"`
}

func (a *API) GetUser(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, OutGetUser{
		Login: user.(db.User).Login,
	})
}

type InUpdateUser struct {
	NewPassword string `json:"password"`
	OldPassword string `json:"old_password"`
	Login       string `json:"login"`
}

func (a *API) UpdateUser(c *gin.Context) {
	var in InUpdateUser
	var out OutRegister
	if err := c.ShouldBindBodyWith(&in, binding.JSON); err != nil {
		out.Message = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// check if user exists
	u, _ := c.Get("user")
	user := u.(db.User)

	if in.OldPassword != "" || in.NewPassword != "" {
		if in.OldPassword == "" {
			out.Message = "old password is empty"
			c.JSON(http.StatusBadRequest, out)
			return
		}

		if in.NewPassword == "" {
			out.Message = "new password is empty"
			c.JSON(http.StatusBadRequest, out)
			return
		}

		err := cryptPasswords.ComparePasswordWithHash(user.Password, in.OldPassword)
		if err != nil {
			out.Message = "wrong password"
			c.JSON(http.StatusUnauthorized, out)
			return
		}

		hash, err := cryptPasswords.GeneratePasswordHash(in.NewPassword)
		if err != nil {
			out.Message = err.Error()
			zap.S().Error(out)
			c.JSON(http.StatusInternalServerError, out)
			return
		}

		user.Password = hash
		if err := a.db.UpdateUser(user); err != nil {
			out.Message = err.Error()
			zap.S().Error(out)
			c.JSON(http.StatusInternalServerError, out)
			return
		}

		out.Message = "ok"
		c.JSON(http.StatusOK, out)
	} else if in.Login != "" {
		user.Login = in.Login
		if err := a.db.UpdateUser(user); err != nil {
			out.Message = err.Error()
			zap.S().Error(out)
			c.JSON(http.StatusInternalServerError, out)
			return
		}

		out.Message = "ok"
		c.JSON(http.StatusOK, out)
	} else {
		out.Message = "nothing to update"
		c.JSON(http.StatusBadRequest, out)
	}
}
