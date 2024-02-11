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

// Ping godoc
//
//	@Summary		Ping
//	@Description	Check connection with server
//	@Tags			ping
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	OutPing
//	@Router			/ping [get]
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

// PostExpression godoc
//
//	@Summary		Add expression
//	@Description	Add expression to storage
//	@Tags			expression
//	@Accept			json
//	@Produce		json
//	@Param			expression	body		InPostExpression	true	"Expression"
//	@Success		200			{object}	OutPostExpression
//	@Failure		400			{object}	OutPostExpression
//	@Router			/expression [post]
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

// GetAllExpressions godoc
//
//	@Summary		Get all expressions
//	@Description	Get all expressions from storage
//	@Tags			expression
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	OutGetAllExpressions
//	@Router			/expression [get]
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

// GetExpressionById godoc
//
//	@Summary		Get expression by id
//	@Description	Get expression from storage by id
//	@Tags			expression
//	@Accept			json
//	@Produce		json
//	@Param			id	body		InGetExpressionById	true	"Expression ID"
//	@Success		200	{object}	OutGetExpressionById
//	@Failure		400	{object}	OutGetExpressionById
//	@Failure		500	{object}	OutGetExpressionById
//	@Router			/expressionById [get]
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

	// get expression from storage
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

// GetOperationsAndTimes godoc
//
//	@Summary		Get operations and times
//	@Description	Get operations and times for calculation as a map of operation and time in milliseconds, {"+": 100,...}
//	@Tags			operations
//	@Accept			json
//	@Produce		json
//	@Param			map[string]int	true		"Operations and times"
//	@Success		200				{object}	OutGetOperationsAndTimes
//	@Router			/getOperationsAndTimes [get]
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

// PostOperationsAndTimes godoc
//
//	@Summary		Set operations and times
//	@Description	Set operations and times for calculation as a map of operation and time in milliseconds, {"+": 100,...}
//	@Tags			operations
//	@Accept			json
//	@Produce		json
//	@Param			data	body		map[string]int	true	"Operations and times"
//	@Success		200		{object}	OutSetOperationsAndTimes
//	@Failure		400		{object}	OutSetOperationsAndTimes
//	@Router			/postOperationsAndTimes [post]
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

// GetUpdates godoc
//
//	@Summary		Get updates
//	@Description	Get not working expressions for calculation server
//	@Tags			updates (used by calculation server)
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	OutGetUpdates
//	@Router			/getUpdates [get]
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

// ConfirmStartCalculating godoc
//
//	@Summary		Confirm start calculating
//	@Description	Confirm start calculating for expression to coordinate work of calculation servers
//	@Tags			updates (used by calculation server)
//	@Accept			json
//	@Produce		json
//	@Param			expression	body		InConfirmStartOfCalculating	true	"Expression"
//	@Success		200			{object}	OutConfirmStartOfCalculating
//	@Failure		400			{object}	OutConfirmStartOfCalculating
//	@Failure		500			{object}	OutConfirmStartOfCalculating
//	@Router			/confirmStartCalculating [post]
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

	// change to working
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

// PostResult godoc
//
//	@Summary		Post result
//	@Description	Post result of the calculation
//	@Tags			updates (used by calculation server)
//	@Accept			json
//	@Produce		json
//	@Param			expression	body		InPostResult	true	"Expression"
//	@Success		200			{object}	OutPostResult
//	@Failure		400			{object}	OutPostResult
//	@Failure		500			{object}	OutPostResult
//	@Router			/postResult [post]
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

	// change to ready
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
