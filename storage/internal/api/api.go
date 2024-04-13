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
	"sync"
)

type API struct {
	db              *db.APIDb
	expressions     *expressionstorage.ExpressionStorage
	servers         *availableservers.AvailableServers
	statusWorkers   *sync.Map
	execTimeConfig  *ExecTimeConfig
	checkAlive      int
	secretSignature []byte
}

func New(_db *db.APIDb, expressions *expressionstorage.ExpressionStorage, statusWorkers *sync.Map, servers *availableservers.AvailableServers, execTimeConfig *ExecTimeConfig) *API {
	num, err := strconv.Atoi(os.Getenv("CHECK_SERVER_DURATION"))
	if err != nil {
		zap.S().Fatal(err)
	}
	newAPI := &API{
		db:              _db,
		statusWorkers:   statusWorkers,
		checkAlive:      num,
		secretSignature: []byte(os.Getenv("SECRET_SIGNATURE")),
		execTimeConfig:  execTimeConfig,
	}
	newAPI.expressions = expressions
	newAPI.servers = servers
	return newAPI
}

func (a *API) Start() *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/api/v1/ping", a.Ping)

	authorized := router.Group("/api/v1")
	authorized.Use(a.Auth)

	// for users
	router.POST("/api/v1/register", a.Register)
	router.POST("/api/v1/login", a.Login)
	router.GET("/api/v1/login", a.Login)

	authorized.GET("/getUser", a.GetUser)
	authorized.POST("/updateUser", a.UpdateUser)
	authorized.POST("/expression", a.PostExpression)
	authorized.GET("/expression", a.GetAllExpressions)
	authorized.GET("/expressionById", a.GetExpressionByID)
	authorized.POST("/postOperationsAndTimes", a.PostOperationsAndTimes)
	authorized.GET("/getOperationsAndTimes", a.GetOperationsAndTimes)
	authorized.GET("/getExpressionsByServer", a.GetExpressionsByServer)
	authorized.GET("/getComputingPowers", a.GetComputingPowers)

	// docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
