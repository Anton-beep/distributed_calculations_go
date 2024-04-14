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

func makeToken(login string, secretSignature []byte) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": login,
		"nbf":  now.Unix(),
		"exp":  now.Add(5 * time.Minute).Unix(),
		"iat":  now.Unix(),
	})

	return token.SignedString(secretSignature)
}

type InRegister struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type OutRegister struct {
	Access  string `json:"access"`
	Message string `json:"message"`
}

// Register godoc
//
//	@Summary		Register
//	@Description	Register new user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			login		body		InRegister	true	"Login"
//	@Param			password	body		InRegister	true	"Password"
//	@Success		200			{object}	OutRegister
//	@Failure		400			{object}	OutRegister
//	@Failure		409			{object}	OutRegister
//	@Router			/register [post]
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

	tokenString, err := makeToken(in.Login, a.secretSignature)
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
	id, err := a.db.AddUser(user)
	if err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusInternalServerError, out)
		return
	}

	// add operations
	operation := db.Operation{
		User: id,
	}
	if _, err = a.db.AddOperation(operation); err != nil {
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
		c.JSON(http.StatusBadRequest, out)
		return
	}

	err = cryptPasswords.ComparePasswordWithHash(user.Password, in.Password)
	if err != nil {
		out.Message = "wrong password"
		c.JSON(http.StatusUnauthorized, out)
		return
	}

	tokenString, err := makeToken(in.Login, a.secretSignature)
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

// GetUser godoc
//
//	@Summary		Get user
//	@Description	Get user info
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	OutGetUser
//	@Failure		400	{object}	OutGetUser
//	@Failure		500	{object}	OutGetUser
//	@Router			/getUser [get]
func (a *API) GetUser(c *gin.Context) {
	user := c.MustGet("user")
	c.JSON(http.StatusOK, OutGetUser{
		Login: user.(db.User).Login,
	})
}

type InUpdateUser struct {
	NewPassword string `json:"password"`
	OldPassword string `json:"old_password"`
	Login       string `json:"login"`
}

type OutUpdateUser struct {
	Access  string `json:"access"`
	Message string `json:"message"`
}

// UpdateUser godoc
//
//	@Summary		Update user
//	@Description	Update user info
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			login			body		InUpdateUser	false	"New login"
//	@Param			password		body		InUpdateUser	false	"New password"
//	@Param			old_password	body		InUpdateUser	false	"Old password"
//	@Success		200				{object}	OutRegister
//	@Failure		400				{object}	OutRegister
//	@Failure		401				{object}	OutRegister
//	@Failure		500				{object}	OutRegister
//	@Router			/updateUser [post]
func (a *API) UpdateUser(c *gin.Context) {
	var in InUpdateUser
	var out OutUpdateUser
	if err := c.ShouldBindBodyWith(&in, binding.JSON); err != nil {
		out.Message = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}
	updated := false
	var tokenString string

	// check if user exists
	u := c.MustGet("user")
	user := u.(db.User)

	if in.OldPassword == "" {
		out.Message = "old password is empty"
		c.JSON(http.StatusBadRequest, out)
		return
	}

	err := cryptPasswords.ComparePasswordWithHash(user.Password, in.OldPassword)
	if err != nil {
		out.Message = "wrong password"
		c.JSON(http.StatusBadRequest, out)
		return
	}

	if in.NewPassword != "" {
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

		// make new token
		tokenString, err = makeToken(in.Login, a.secretSignature)
		if err != nil {
			out.Message = err.Error()
			zap.S().Error(out)
			c.JSON(http.StatusInternalServerError, out)
		}

		updated = true
	}
	if in.Login != "" {
		user.Login = in.Login
		if err := a.db.UpdateUser(user); err != nil {
			out.Message = err.Error()
			zap.S().Error(out)
			c.JSON(http.StatusInternalServerError, out)
			return
		}

		// make new token
		tokenString, err = makeToken(in.Login, a.secretSignature)
		if err != nil {
			out.Message = err.Error()
			zap.S().Error(out)
			c.JSON(http.StatusInternalServerError, out)
		}

		updated = true
	}
	if updated {
		out.Access = tokenString
		out.Message = "ok"
		c.JSON(http.StatusOK, out)
	} else {
		out.Message = "nothing to update"
		c.JSON(http.StatusBadRequest, out)
	}
}
