package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"storage/internal/db"
	"storage/internal/expressionstorage"
)

type API struct {
	db             *db.APIDb
	expressions    *expressionstorage.ExpressionStorage
	execTimeConfig ExecTimeConfig
}

func New(_db *db.APIDb) *API {
	newAPI := &API{
		db:          _db,
		expressions: expressionstorage.New(_db),
	}
	return newAPI
}

func (a *API) Start() *gin.Engine {
	router := gin.Default()

	router.GET("/api/v1/ping", a.Ping)

	// for user
	router.POST("/api/v1/expression", a.PostExpression)
	router.GET("/api/v1/expression", a.GetAllExpressions)
	router.GET("/api/v1/expressionById", a.GetExpressionByID)
	router.POST("/api/v1/postOperationsAndTimes", a.PostOperationsAndTimes)
	router.GET("/api/v1/getOperationsAndTimes", a.GetOperationsAndTimes)

	// for calculation server
	router.GET("/api/v1/getUpdates", a.GetUpdates)
	router.POST("/api/v1/confirmStartCalculating", a.ConfirmStartCalculating)
	router.POST("/api/v1/postResult", a.PostResult)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
