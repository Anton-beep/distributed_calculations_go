package Api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"storage/internal/Db"
	"time"
)

type OutPing struct {
	Message string `json:"message"`
}

func (a *Api) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, OutPing{Message: "pong"})
}

// for user

type InPostExpression struct {
	Expression string `json:"expression" binding:"required"`
}

type OutPostExpression struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

func (a *Api) PostExpression(c *gin.Context) {
	var in InPostExpression
	var out OutPostExpression
	if err := c.ShouldBindJSON(&in); err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// add expression to storage
	newExpression := Db.Expression{
		Id:     0,
		Value:  in.Expression,
		Answer: 0,
		Logs:   "",
		Status: Db.ExpressionNotReady,
	}
	newId, err := a.expressions.Add(newExpression)
	if err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusInternalServerError, out)
		return
	}

	out.Id = newId
	out.Message = "ok"
	c.JSON(http.StatusOK, out)
}

type OutGetAllExpressions struct {
	Expressions []Db.Expression `json:"expressions"`
	Message     string          `json:"message"`
}

func (a *Api) GetAllExpressions(c *gin.Context) {
	expressions := a.expressions.GetAll()
	c.JSON(http.StatusOK, OutGetAllExpressions{Expressions: expressions, Message: "ok"})
}

type InGetExpressionById struct {
	Id int `json:"id" binding:"required"`
}

type OutGetExpressionById struct {
	Expression Db.Expression `json:"expression"`
	Message    string        `json:"message"`
}

func (a *Api) GetExpressionById(c *gin.Context) {
	var in InGetExpressionById
	var out OutGetExpressionById
	if err := c.ShouldBindJSON(&in); err != nil {
		out.Expression = Db.Expression{}
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	expression, err := a.expressions.GetById(in.Id)
	if err != nil {
		out.Expression = Db.Expression{}
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusInternalServerError, out)
		return
	}

	out.Expression = expression
	out.Message = "ok"
	c.JSON(http.StatusOK, out)
}

type ExecTimeConfig struct {
	TimeAdd      time.Duration
	TimeSubtract time.Duration
	TimeDivide   time.Duration
	TimeMultiply time.Duration
}

type OutGetOperationsAndTimes struct {
	Data    map[string]int `json:"data"`
	Message string         `json:"message"`
}

func (a *Api) GetOperationsAndTimes(c *gin.Context) {
	outMap := make(map[string]int)
	outMap["+"] = int(a.execTimeConfig.TimeAdd.Milliseconds())
	outMap["-"] = int(a.execTimeConfig.TimeSubtract.Milliseconds())
	outMap["/"] = int(a.execTimeConfig.TimeDivide.Milliseconds())
	outMap["*"] = int(a.execTimeConfig.TimeMultiply.Milliseconds())
	c.JSON(http.StatusOK, OutGetOperationsAndTimes{Data: outMap, Message: "ok"})
}

type OutSetOperationsAndTimes struct {
	Message string `json:"message"`
}

func (a *Api) PostOperationsAndTimes(c *gin.Context) {
	var in map[string]int

	if err := c.ShouldBindJSON(&in); err != nil {
		out := OutSetOperationsAndTimes{Message: err.Error()}
		zap.S().Error(out)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	msg := ""
	for key, value := range in {
		switch key {
		case "+":
			a.execTimeConfig.TimeAdd = time.Duration(value) * time.Millisecond
			msg += "changed for +;"
		case "-":
			a.execTimeConfig.TimeSubtract = time.Duration(value) * time.Millisecond
			msg += "changed for -;"
		case "/":
			a.execTimeConfig.TimeDivide = time.Duration(value) * time.Millisecond
			msg += "changed for /;"
		case "*":
			a.execTimeConfig.TimeMultiply = time.Duration(value) * time.Millisecond
			msg += "changed for *;"
		}
	}

	out := OutSetOperationsAndTimes{Message: msg}
	c.JSON(http.StatusOK, out)
}

// for calculation server

type OutGetUpdates struct {
	Expressions []Db.Expression `json:"tasks" binding:"required"`
	Message     string          `json:"message"`
}

func (a *Api) GetUpdates(c *gin.Context) {
	expressions := a.expressions.GetNotWorkingExpressions()
	out := OutGetUpdates{Expressions: expressions, Message: "ok"}
	c.JSON(http.StatusOK, out)
}

type InConfirmStartOfCalculating struct {
	Expression Db.Expression `json:"expression" binding:"required"`
}

type OutConfirmStartOfCalculating struct {
	Confirm bool   `json:"confirm"`
	Message string `json:"message"`
}

func (a *Api) ConfirmStartCalculating(c *gin.Context) {
	var in InConfirmStartOfCalculating
	var out OutConfirmStartOfCalculating
	if err := c.ShouldBindJSON(&in); err != nil {
		out.Confirm = false
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// check if expression is not ready
	ok, err := a.expressions.IsExpressionNotReady(in.Expression.Id)
	if err != nil {
		out.Confirm = false
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusInternalServerError, out)
	}
	if !ok {
		out.Confirm = false
		out.Message = "expression is not in pending"
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// move to working
	in.Expression.Status = Db.ExpressionWorking
	if err := a.expressions.UpdatePendingExpression(in.Expression); err != nil {
		out.Confirm = false
		c.JSON(http.StatusInternalServerError, out)
		return
	}

	out.Confirm = true
	out.Message = "ok"
	c.JSON(http.StatusOK, out)
}

type InPostResult struct {
	Expression Db.Expression `json:"expression" binding:"required"`
}

type OutPostResult struct {
	Message string `json:"message"`
}

func (a *Api) PostResult(c *gin.Context) {
	var in InPostResult
	var out OutPostResult
	if err := c.ShouldBindJSON(&in); err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// check if expression is in working
	ok, err := a.expressions.IsExpressionWorking(in.Expression.Id)
	if err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusInternalServerError, out)
		return
	}
	if !ok {
		out.Message = "expression is not working"
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// move to ready
	in.Expression.Status = Db.ExpressionReady
	if err := a.expressions.PendingToReady(in.Expression); err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusInternalServerError, out)
		return
	}

	out.Message = "ok"
	c.JSON(http.StatusOK, out)
}
