package Api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type authData struct {
	Secret string `json:"secret" binding:"required"`
}

func (a *Api) checkAuth(c *gin.Context) {
	var data authData
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if data.Secret != a.secret {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "wrong secret"})
		return
	}

	c.Next()
}
