package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"storage/internal/db"
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
func (a *API) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, OutPing{Message: "pong"})
}

// for user

type InPostExpression struct {
	Expression string `json:"expression" binding:"required"`
}

type OutPostExpression struct {
	ID      int    `json:"id"`
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
func (a *API) PostExpression(c *gin.Context) {
	var in InPostExpression
	var out OutPostExpression
	if err := c.ShouldBindJSON(&in); err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// add expression to storage
	newExpression := db.Expression{
		ID:           0,
		Value:        in.Expression,
		Answer:       0,
		Logs:         "",
		Status:       db.ExpressionNotReady,
		CreationTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	newID, err := a.expressions.Add(newExpression)
	if err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusInternalServerError, out)
		return
	}

	out.ID = newID
	out.Message = "ok"
	c.JSON(http.StatusOK, out)
}

type OutGetAllExpressions struct {
	Expressions []db.Expression `json:"expressions"`
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
func (a *API) GetAllExpressions(c *gin.Context) {
	expressions := a.expressions.GetAll()
	c.JSON(http.StatusOK, OutGetAllExpressions{Expressions: expressions, Message: "ok"})
}

type InGetExpressionByID struct {
	ID int `json:"id" binding:"required"`
}

type OutGetExpressionByID struct {
	Expression db.Expression `json:"expression"`
	Message    string        `json:"message"`
}

// GetExpressionByID godoc
//
//	@Summary		Get expression by id
//	@Description	Get expression from storage by id
//	@Tags			expression
//	@Accept			json
//	@Produce		json
//	@Param			id	body		InGetExpressionByID	true	"Expression ID"
//	@Success		200	{object}	OutGetExpressionByID
//	@Failure		400	{object}	OutGetExpressionByID
//	@Failure		500	{object}	OutGetExpressionByID
//	@Router			/expressionById [get]
func (a *API) GetExpressionByID(c *gin.Context) {
	var in InGetExpressionByID
	var out OutGetExpressionByID
	if err := c.ShouldBindJSON(&in); err != nil {
		out.Expression = db.Expression{}
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// get expression from storage
	expression, err := a.expressions.GetByID(in.ID)
	if err != nil {
		out.Expression = db.Expression{}
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusBadRequest, out)
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
	Data    map[string]int `json:"data"` // executions times in milliseconds: {"+": 100,...}
	Message string         `json:"message"`
}

// GetOperationsAndTimes godoc
//
//	@Summary		Get operations and times
//	@Description	Get operations and times for calculation as a map of operation and time in milliseconds, {"+": 100,...}
//	@Tags			operations
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	OutGetOperationsAndTimes
//	@Router			/getOperationsAndTimes [get]
func (a *API) GetOperationsAndTimes(c *gin.Context) {
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
func (a *API) PostOperationsAndTimes(c *gin.Context) {
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
	Expressions []db.Expression `json:"tasks" binding:"required"`
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
func (a *API) GetUpdates(c *gin.Context) {
	expressions := a.expressions.GetNotWorkingExpressions()
	out := OutGetUpdates{Expressions: expressions, Message: "ok"}
	c.JSON(http.StatusOK, out)
}

type InConfirmStartOfCalculating struct {
	Expression db.Expression `json:"expression" binding:"required"`
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
func (a *API) ConfirmStartCalculating(c *gin.Context) {
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
	ok, err := a.expressions.IsExpressionNotReady(in.Expression.ID)
	if err != nil {
		out.Confirm = false
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusInternalServerError, out)
		return
	}
	if !ok {
		out.Confirm = false
		out.Message = "expression is not in pending"
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// change to working
	in.Expression.Status = db.ExpressionWorking
	in.Expression.AliveExpiresAt = int(time.Now().Add(time.Duration(a.checkAlive) * time.Second).Unix())
	if err = a.expressions.UpdateExpression(in.Expression); err != nil {
		out.Confirm = false
		c.JSON(http.StatusInternalServerError, out)
		return
	}

	// add server
	a.servers.Add(in.Expression.Servername)

	out.Confirm = true
	out.Message = "ok"
	c.JSON(http.StatusOK, out)
}

type InPostResult struct {
	Expression db.Expression `json:"expression" binding:"required"`
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
func (a *API) PostResult(c *gin.Context) {
	var in InPostResult
	var out OutPostResult
	if err := c.ShouldBindJSON(&in); err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// check if expression is in working
	ok, err := a.expressions.IsExpressionWorking(in.Expression.ID)
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

	in.Expression.EndCalculationTime = time.Now().Format("2006-01-02 15:04:05")
	if err = a.expressions.UpdateExpression(in.Expression); err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusInternalServerError, out)
		return
	}

	// add server
	a.servers.Add(in.Expression.Servername)
	a.statusWorkers.Store(in.Expression.Servername, fmt.Sprintf("%v -> server %v finished calculating %v",
		time.Now().Format("01-02-2006 15:04:05"), in.Expression.Servername, in.Expression.Value))

	out.Message = "ok"
	c.JSON(http.StatusOK, out)
}

type InKeepAlive struct {
	Expression    db.Expression `json:"expression" binding:"required"`
	StatusWorkers string        `json:"status_workers" binding:"required"`
}

type OutKeepAlive struct {
	Message string `json:"message"`
}

// KeepAlive godoc
//
//	@Summary		Keep alive
//	@Description	Keep alive for expression to coordinate work of calculation servers
//	@Tags			updates (used by calculation server)
//	@Accept			json
//	@Produce		json
//	@Param			expression	body		InKeepAlive	true	"Expression"
//	@Success		200			{object}	OutPing
//	@Failure		400			{object}	OutPing
//	@Failure		500			{object}	OutPing
//	@Router			/keepAlive [post]
func (a *API) KeepAlive(c *gin.Context) {
	var in InKeepAlive
	var out = OutKeepAlive{}
	if err := c.ShouldBindJSON(&in); err != nil {
		out.Message = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	expression, err := a.expressions.GetByID(in.Expression.ID)
	if err != nil {
		out.Message = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	expression.AliveExpiresAt = int(time.Now().Add(time.Duration(a.checkAlive) * time.Second).Unix())
	err = a.expressions.UpdateExpression(expression)
	if err != nil {
		out.Message = err.Error()
		c.JSON(http.StatusInternalServerError, out)
		return
	}

	a.statusWorkers.Store(expression.Servername, in.StatusWorkers)

	c.JSON(http.StatusOK, OutPing{Message: "ok"})
}

type InGetExpressionByServer struct {
	ServerName string `json:"server_name" binding:"required"`
}

type OutGetExpressionByServer struct {
	Expressions []db.Expression `json:"expressions"`
	Message     string          `json:"message"`
}

// GetExpressionsByServer godoc
//
//	@Summary		Get expression by server
//	@Description	Get expressions from storage by server name
//	@Tags			expression
//	@Accept			json
//	@Produce		json
//	@Param			server_name	body		InGetExpressionByServer	true	"Server name"
//	@Success		200			{object}	OutGetExpressionByServer
//	@Failure		400			{object}	OutGetExpressionByServer
//	@Router			/getExpressionByServer [get]
func (a *API) GetExpressionsByServer(c *gin.Context) {
	var in InGetExpressionByServer
	var out OutGetExpressionByServer
	if err := c.ShouldBindJSON(&in); err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// get expressions from storage
	expressions := a.expressions.GetByServer(in.ServerName)
	out.Expressions = expressions
	out.Message = "ok"
	c.JSON(http.StatusOK, out)
}

type OutGetComputingPowers struct {
	Servers []struct {
		ServerName            string `json:"server_name"`
		CalculatedExpressions []int  `json:"calculated_expressions"`
		ServerStatus          string `json:"server_status"`
	} `json:"servers"`
	Message string `json:"message"`
}

// GetComputingPowers godoc
//
//	@Summary		Get computing powers
//	@Description	Get computing powers from storage
//	@Tags			computing powers
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	OutGetComputingPowers
//	@Router			/getComputingPowers [get]
func (a *API) GetComputingPowers(c *gin.Context) {
	var out OutGetComputingPowers

	// get computing powers
	servers := a.servers.GetAll()
	for _, server := range servers {
		operations := a.servers.GetExpressions(server)
		ids := make([]int, 0)
		for _, expression := range operations {
			ids = append(ids, expression.ID)
		}

		val, ok := a.statusWorkers.Load(server)
		if !ok {
			val = "unknown"
		}
		out.Servers = append(out.Servers, struct {
			ServerName            string `json:"server_name"`
			CalculatedExpressions []int  `json:"calculated_expressions"`
			ServerStatus          string `json:"server_status"`
		}{ServerName: server, CalculatedExpressions: ids, ServerStatus: val.(string)})
	}

	out.Message = "ok"
	c.JSON(http.StatusOK, out)
}
