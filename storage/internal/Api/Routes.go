package Api

import (
	"github.com/gin-gonic/gin"
	"storage/internal/Db"
	"time"
)

// for user

type InPostExpression struct {
	Expression string `json:"expression" binding:"required"`
}

type OutPostExpression struct {
	Message string `json:"message"`
}

func (a *Api) PostExpression(c *gin.Context) {

}

type OutGetAllExpressions struct {
	Expressions []Db.Expression `json:"expressions"`
}

func (a *Api) GetAllExpressions(c *gin.Context) {

}

type InGetExpressionById struct {
	Id int `json:"id" binding:"required"`
}

type OutGetExpressionById struct {
	Expression Db.Expression `json:"expression"`
}

func (a *Api) GetExpressionById(c *gin.Context) {

}

type ExecTimeConfig struct {
	TimeAdd      time.Duration
	TimeSubtract time.Duration
	TimeDivide   time.Duration
	TimeMultiply time.Duration
}

type OutGetOperationsAndTimes struct {
	Data map[string]int `json:"data"`
}

func (a *Api) GetOperationsAndTimes(c *gin.Context) {

}

// for calculation server

type Task struct {
	Id         int    `json:"id"`
	Expression string `json:"expression"`
}

type OutGetUpdates struct {
	Tasks []Task `json:"tasks" binding:"required"`
}

func (a *Api) GetUpdates(c *gin.Context) {

}

type OutConfirmStartOfCalculating struct {
	Confirm bool `json:"confirm"`
}

func (a *Api) ConfirmStartOfCalculating(c *gin.Context) {

}

type InPostResult struct {
	Id     int    `json:"id" binding:"required"`
	Logs   string `json:"logs" binding:"required"`
	Answer int    `json:"answer" binding:"required"`
}

type OutPostResult struct {
	Message string `json:"message"`
}

func (a *Api) PostResult(c *gin.Context) {

}
