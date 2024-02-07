package Api

import (
	_ "calculationServer/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Api struct {
	storageUrl string
	secret     string
}

func NewApi(storageUrl, secret string) *Api {
	return &Api{
		storageUrl: storageUrl,
		secret:     secret,
	}
}

func (a *Api) SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/api/v1/ping", a.pong)

	// authorized := r.Group("/auth", a.checkAuth)
	// swagger (documentation)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}

func (a *Api) Start(router *gin.Engine) {
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
