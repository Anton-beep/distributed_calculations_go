package storageclient

import (
	"calculationServer/pkg/expressionparser"
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"strconv"
	"time"
)

const (
	ExpressionNotReady = 0
	ExpressionWorking  = 1
	ExpressionReady    = 2
	ExpressionError    = 3
)

type Client struct {
	storageServer    string
	expressionParser *expressionparser.ExpressionParser
	keepAlive        time.Duration
	serverName       string
	connection       *grpc.ClientConn
	gRPCClient       ExpressionsServiceClient
}

/*
	type Expression struct {
		ID                 int     `json:"id"`
		Value              string  `json:"value"`
		Answer             float64 `json:"answer"`
		Logs               string  `json:"logs"`
		Status             int     `json:"ready"` // 0 - not ready, 1 - working, 2 - ready, 3 - error
		AliveExpiresAt     int     `json:"alive_expires_at"`
		CreationTime       string  `json:"creation_time"`
		EndCalculationTime string  `json:"end_calculation_time"`
		Servername         string  `json:"server_name"`
		User               int     `db:"user_id" json:"user_id"`
	}
*/
func New() (*Client, error) {
	c := &Client{}
	c.storageServer = os.Getenv("STORAGE_URL")

	conn, err := grpc.Dial(c.storageServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c.connection = conn

	c.gRPCClient = NewExpressionsServiceClient(conn)

	c.expressionParser = expressionparser.New()

	num, err := strconv.Atoi(os.Getenv("NUMBER_OF_CALCULATORS"))
	if err != nil {
		return nil, err
	}
	err = c.expressionParser.SetNumberOfWorkers(num)
	if err != nil {
		return nil, err
	}

	num, err = strconv.Atoi(os.Getenv("SEND_ALIVE_DURATION"))
	if err != nil {
		return nil, err
	}
	c.keepAlive = time.Duration(num) * time.Second

	c.serverName = os.Getenv("CALCULATION_SERVER_NAME")
	if c.serverName == "" {
		c.serverName = "noname"
	}
	return c, nil
}

func (c *Client) CloseConn() error {
	err := c.connection.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) tryGetUpdates() (*Expression, bool) {
	zap.S().Info("try to get updates")
	expressions, err := c.GetUpdates()
	if err != nil {
		zap.S().Error(err)
	}
	// try to take first expression for calculation
	if len(expressions) == 0 {
		zap.S().Info("no expressions")
		time.Sleep(2000 * time.Millisecond)
		return nil, false
	}
	exp := expressions[0]
	zap.S().Infof("got expression: %v", exp)
	exp.ServerName = c.serverName
	ok, err := c.TryToConfirm(exp)
	if err != nil {
		zap.S().Error(err)
	}
	if ok {
		// server can calculate this expression
		zap.S().Info("confirmed")
		return exp, true
	}
	time.Sleep(2000 * time.Millisecond)
	zap.S().Info("can't confirm, try to get updates again")
	return nil, false
}

func (c *Client) tryUpdateTimeConfig() {
	config, err := c.GetOperationsAndTimes()
	if err != nil {
		zap.S().Error(err)
	}
	err = c.expressionParser.SetExecTimes(config)
	if err != nil {
		zap.S().Error(err)
	}
	zap.S().Info("exec time config updated")
}

func (c *Client) trySendResult(exp *Expression) {
	// try 10 times
	for i := 0; i < 10; i++ {
		zap.S().Info("try to send result")
		ok, err := c.SendResult(exp)
		if err != nil {
			zap.S().Error(err)
		}
		if ok {
			zap.S().Info("result sent successfully")
			break
		}
		time.Sleep(2000 * time.Millisecond)
		zap.S().Info("can't send result, try to send again")
	}
}

func (c *Client) keepAliveExpression(exp *Expression, done <-chan bool, ticker *time.Ticker) {
	for {
		select {
		case <-done:
			zap.S().Info("calculation done")
			return
		case <-ticker.C:
			zap.S().Info("send alive")
			err := c.KeepAlive(exp)
			if err != nil {
				zap.S().Error(err)
			}
		}
	}
}

func (c *Client) Run() {
	for {
		exp := &Expression{}
		var ok bool
		for {
			exp, ok = c.tryGetUpdates()
			if ok {
				break
			}
		}

		// update exec time config
		c.tryUpdateTimeConfig()

		ticker := time.NewTicker(c.keepAlive)
		done := make(chan bool)
		// keep this client alive for the server
		go c.keepAliveExpression(exp, done, ticker)
		res, logs, err := c.expressionParser.CalculateExpression(exp.Value)
		ticker.Stop()
		done <- true
		if err != nil {
			zap.S().Error(err)
			exp.Status = ExpressionError
			exp.Logs = err.Error()
		} else {
			exp.Status = ExpressionReady
			exp.Logs = logs
		}
		exp.Answer = res
		zap.S().Infof("result: %v", exp)

		// send result
		c.trySendResult(exp)
	}
}

type AnsGetUpdates struct {
	Tasks   []Expression `json:"tasks" binding:"required"`
	Message string       `json:"message"`
}

// GetUpdates returns all expressions for calculation.
func (c *Client) GetUpdates() ([]*Expression, error) {
	updates, err := c.gRPCClient.GetUpdates(
		context.Background(),
		&Empty{},
	)

	if err != nil {
		return nil, err
	}

	res := make([]*Expression, 0)
	for {
		expression, err := updates.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		res = append(res, expression)
	}

	return res, nil
}

type SendConfirmStartOfCalculating struct {
	Expression Expression `json:"expression"`
}

type AnsConfirmStartOfCalculating struct {
	Confirm bool   `json:"confirm"`
	Message string `json:"message"`
}

// TryToConfirm returns true if the expression is not being calculated by another server.
func (c *Client) TryToConfirm(expression *Expression) (bool, error) {
	confirm, err := c.gRPCClient.ConfirmStartCalculating(
		context.Background(),
		expression,
	)
	if confirm == nil {
		return false, fmt.Errorf("confirm is nil")
	}

	return confirm.Confirm, err
}

type sendResult struct {
	Expression Expression `json:"expression"`
}

type AnsSendResult struct {
	Message string `json:"message"`
}

// SendResult sends the result of the expression to the storage.
func (c *Client) SendResult(expression *Expression) (bool, error) {
	msg, err := c.gRPCClient.PostResult(
		context.Background(),
		expression,
	)
	if msg == nil {
		return false, fmt.Errorf("msg is nil")
	}

	return msg.Message == "ok", err
}

type AnsGetOperationsAndTimes struct {
	Data    map[string]int `json:"data"`
	Message string         `json:"message"`
}

// GetOperationsAndTimes returns the time for each operation from the storage.
func (c *Client) GetOperationsAndTimes() (expressionparser.ExecTimeConfig, error) {
	ans, err := c.gRPCClient.GetOperationsAndTimes(
		context.Background(),
		&Empty{},
	)

	if err != nil {
		return expressionparser.ExecTimeConfig{}, err
	}

	var config expressionparser.ExecTimeConfig
	for key, val := range ans.Operations {
		switch key {
		case "+":
			config.TimeAdd = time.Duration(val) * time.Millisecond
		case "-":
			config.TimeSubtract = time.Duration(val) * time.Millisecond
		case "/":
			config.TimeDivide = time.Duration(val) * time.Millisecond
		case "*":
			config.TimeMultiply = time.Duration(val) * time.Millisecond
		default:
			return expressionparser.ExecTimeConfig{}, fmt.Errorf("unknown operator: %s", key)
		}
	}

	return config, nil
}

func (c *Client) KeepAlive(expression *Expression) error {
	var send KeepAliveMsg
	send.Expression = expression
	send.StatusWorkers = fmt.Sprintf("%v -> %v from %v workers are runninng to calcualte %v",
		time.Now().Format("01-02-2006 15:04:05"), c.expressionParser.GetWorkingWorkers(),
		c.expressionParser.GetTotalNumberOfWorkers(), expression.Value)
	_, err := c.gRPCClient.KeepAlive(
		context.Background(),
		&send,
	)
	return err
}
