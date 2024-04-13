package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	if err := c.ShouldBindBodyWith(&in, binding.JSON); err != nil {
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
		User:         c.MustGet("user").(db.User).ID,
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
	user := c.MustGet("user").(db.User)
	expressions := a.expressions.GetAll(user.ID)
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
	if err := c.ShouldBindBodyWith(&in, binding.JSON); err != nil {
		out.Expression = db.Expression{}
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// get expression from storage
	user := c.MustGet("user").(db.User)

	expression, err := a.expressions.GetByUserAndID(user.ID, in.ID)
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
	operations, err := a.db.GetUserOperations(c.MustGet("user").(db.User).ID)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusInternalServerError, OutGetOperationsAndTimes{Message: err.Error()})
		return
	}
	outMap := make(map[string]int)
	outMap["+"] = operations.TimeAdd
	outMap["-"] = operations.TimeSubtract
	outMap["/"] = operations.TimeDivide
	outMap["*"] = operations.TimeMultiply
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

	if err := c.ShouldBindBodyWith(&in, binding.JSON); err != nil {
		out := OutSetOperationsAndTimes{Message: err.Error()}
		zap.S().Error(out)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	operations, err := a.db.GetUserOperations(c.MustGet("user").(db.User).ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, OutSetOperationsAndTimes{Message: err.Error()})
		return
	}

	msg := ""
	for key, value := range in {
		switch key {
		case "+":
			operations.TimeAdd = value
			msg += "changed for +;"
		case "-":
			operations.TimeSubtract = value
			msg += "changed for -;"
		case "/":
			operations.TimeDivide = value
			msg += "changed for /;"
		case "*":
			operations.TimeMultiply = value
			msg += "changed for *;"
		}
	}

	err = a.db.UpdateOperation(operations)

	out := OutSetOperationsAndTimes{Message: msg}
	c.JSON(http.StatusOK, out)
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
	if err := c.ShouldBindBodyWith(&in, binding.JSON); err != nil {
		out.Message = err.Error()
		zap.S().Error(out)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// get expressions from storage
	user := c.MustGet("user").(db.User)
	expressions := a.expressions.GetByServer(user.ID, in.ServerName)
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
	user := c.MustGet("user").(db.User)
	servers := a.servers.GetAll()
	for _, server := range servers {
		operations := a.servers.GetExpressions(user.ID, server)
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
