package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"os"
	"storage/internal/availableservers"
	"storage/internal/db"
	"storage/internal/expressionstorage"
	"strconv"
	"time"
)

type API struct {
	db             *db.APIDb
	expressions    *expressionstorage.ExpressionStorage
	servers        *availableservers.AvailableServers
	execTimeConfig ExecTimeConfig
	checkAlive     int
}

func New(_db *db.APIDb) *API {
	num, err := strconv.Atoi(os.Getenv("CHECK_SERVER_DURATION"))
	if err != nil {
		zap.S().Fatal(err)
	}
	newAPI := &API{
		db:          _db,
		expressions: expressionstorage.New(_db, time.Duration(num)*time.Second),
		checkAlive:  num,
	}
	newAPI.servers = availableservers.New(newAPI.expressions)
	return newAPI
}

func (a *API) Start() *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/api/v1/ping", a.Ping)

	// for user
	router.POST("/api/v1/expression", a.PostExpression)
	router.GET("/api/v1/expression", a.GetAllExpressions)
	router.GET("/api/v1/expressionById", a.GetExpressionByID)
	router.POST("/api/v1/postOperationsAndTimes", a.PostOperationsAndTimes)
	router.GET("/api/v1/getOperationsAndTimes", a.GetOperationsAndTimes)
	router.GET("/api/v1/getExpressionsByServer", a.GetExpressionsByServer)
	router.GET("/api/v1/getComputingPowers", a.GetComputingPowers)

	// for calculation server
	router.GET("/api/v1/getUpdates", a.GetUpdates)
	router.POST("/api/v1/confirmStartCalculating", a.ConfirmStartCalculating)
	router.POST("/api/v1/postResult", a.PostResult)
	router.POST("/api/v1/keepAlive", a.KeepAlive)

	// docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
