package Api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"storage/internal/Db"
	"storage/internal/ExpressionStorage"
)

type Api struct {
	db             *Db.ApiDb
	expressions    *ExpressionStorage.ExpressionStorage
	execTimeConfig ExecTimeConfig
}

func New(_db *Db.ApiDb) *Api {
	newApi := &Api{
		db:          _db,
		expressions: ExpressionStorage.New(_db),
	}
	return newApi
}

func (a *Api) Start() *gin.Engine {
	router := gin.Default()

	router.GET("/api/v1/ping", a.Ping)

	// for user
	router.POST("/api/v1/expression", a.PostExpression)
	router.GET("/api/v1/expression", a.GetAllExpressions)
	router.GET("/api/v1/expressionById", a.GetExpressionById)
	router.POST("/api/v1/postOperationsAndTimes", a.PostOperationsAndTimes)
	router.GET("/api/v1/getOperationsAndTimes", a.GetOperationsAndTimes)

	// for calculation server
	router.GET("/api/v1/getUpdates", a.GetUpdates)
	router.POST("/api/v1/confirmStartCalculating", a.ConfirmStartCalculating)
	router.POST("/api/v1/postResult", a.PostResult)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
