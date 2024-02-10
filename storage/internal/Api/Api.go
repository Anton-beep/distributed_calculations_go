package Api

import (
	"github.com/gin-gonic/gin"
	"storage/internal/Db"
)

type Api struct {
	db          *Db.ApiDb
	expressions map[int]Db.Expression
}

func New(_db *Db.ApiDb) *Api {
	newApi := &Api{db: _db}
	return newApi
}

func (a *Api) Start() *gin.Engine {
	router := gin.Default()

	//router.GET("/api/ping", a.Pong)
	//router.POST("/api/register", a.Register)
	//router.POST("/api/login", a.Login)
	//
	//authGroup := router.Group("/api/private")
	//authGroup.Use(a.AuthMiddleware())
	//
	//authGroup.POST("getChats", a.GetUsersChats)
	//authGroup.POST("getMessagesByChatID", a.GetMessagesByChatID)
	//authGroup.POST("getInfoUser", a.GetInfoAboutUser)
	//authGroup.POST("editMessage", a.EditMessage)
	//authGroup.POST("editStatus", a.EditStatus)
	//authGroup.POST("createMessage", a.CreateNewMessage)
	//authGroup.POST("createChat", a.CreateNewChat)
	//authGroup.POST("getUpdatesMessage", a.GetMessageUpdates)
	//authGroup.POST("isUserExists", a.IsUserExists)
	//authGroup.POST("createChatByUsernames", a.CreateChatByUsernames)
	//authGroup.POST("getInfoChat", a.GetInfoChat)

	return router
}
