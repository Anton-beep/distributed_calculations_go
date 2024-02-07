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

func (a *Api) Start() {
	router := gin.Default()

	router.GET("/api/v1/ping", a.pong)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
