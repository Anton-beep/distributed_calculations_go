package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *API) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}
