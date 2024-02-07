package Api

import "github.com/gin-gonic/gin"

type Pong struct {
	Message string `json:"message"`
}

// PingPong godoc
//
//	@Summary	ping to check a server
//	@Schemes
//	@Description	do ping
//	@Tags			ping
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	pong
//	@Router			/ping [get]
func (a *Api) pong(c *gin.Context) {
	c.JSON(200, Pong{Message: "pong"})
}

type newTask struct {
}

func (a *Api) postTask(c *gin.Context) {

}

func (a *Api) getStatus(c *gin.Context) {

}
