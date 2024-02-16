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

func (a *API) InputExpression(c *gin.Context) {
}

func (a *API) AllExpressions(c *gin.Context) {
}

func (a *API) Operations(c *gin.Context) {
}

func (a *API) ComputingPower(c *gin.Context) {
}
