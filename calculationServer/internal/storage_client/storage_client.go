package storage_client

import (
	"calculationServer/pkg/expression_parser"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	expressionParser *expression_parser.ExpressionParser
	keepAlive        time.Duration
}

type Expression struct {
	ID             int     `json:"id"`
	Value          string  `json:"value"`
	Answer         float64 `json:"answer"`
	Logs           string  `json:"logs"`
	Status         int     `json:"ready"` // 0 - not ready, 1 - working, 2 - ready, 3 - error
	AliveExpiresAt int     `json:"alive_experise_at"`
}

func New() (*Client, error) {
	c := &Client{}
	c.storageServer = os.Getenv("STORAGE_URL")
	c.expressionParser = expression_parser.New()

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
	return c, nil
}

func (c *Client) tryGetUpdates() (Expression, bool) {
	zap.S().Info("try to get updates")
	expressions, err := c.GetUpdates()
	if err != nil {
		zap.S().Error(err)
	}
	// try to take first expression for calculation
	if len(expressions) == 0 {
		zap.S().Info("no expressions")
		time.Sleep(2000 * time.Millisecond)
		return Expression{}, false
	}
	exp := expressions[0]
	zap.S().Infof("got expression: %v", exp)
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
	return Expression{}, false
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

func (c *Client) trySendResult(exp Expression) {
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

func (c *Client) keepAliveExpression(exp Expression, done <-chan bool, ticker *time.Ticker) {
	for {
		select {
		case <-done:
			zap.S().Info("calculation done")
			return
		case <-ticker.C:
			zap.S().Info("send alive")
			err := c.SendAlive(exp)
			if err != nil {
				zap.S().Error(err)
			}
		}
	}
}

func (c *Client) Run() {
	for {
		exp := Expression{}
		ok := false
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
		}
		exp.Answer = res
		exp.Logs = logs
		exp.Status = ExpressionError
		zap.S().Infof("result: %v", exp)

		// send result
		c.trySendResult(exp)
	}
}

type AnsGetUpdates struct {
	Tasks   []Expression `json:"tasks" binding:"required"`
	Message string       `json:"message"`
}

// GetUpdates returns all expressions for calculation
func (c *Client) GetUpdates() ([]Expression, error) {
	resp, err := http.Get(c.storageServer + "/getUpdates")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ans AnsGetUpdates
	if err = json.Unmarshal(body, &ans); err != nil {
		return nil, err
	}

	return ans.Tasks, nil
}

type SendConfirmStartOfCalculating struct {
	Expression Expression `json:"expression"`
}

type AnsConfirmStartOfCalculating struct {
	Confirm bool   `json:"confirm"`
	Message string `json:"message"`
}

// TryToConfirm returns true if the expression is not being calculated by another server
func (c *Client) TryToConfirm(expression Expression) (bool, error) {
	data := SendConfirmStartOfCalculating{Expression: expression}
	body, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	resp, err := http.Post(c.storageServer+"/confirmStartCalculating", "application/json", strings.NewReader(string(body)))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var ans AnsConfirmStartOfCalculating
	if err = json.Unmarshal(body, &ans); err != nil {
		return false, err
	}
	return ans.Confirm, nil
}

type sendResult struct {
	Expression Expression `json:"expression"`
}

type AnsSendResult struct {
	Message string `json:"message"`
}

// SendResult sends the result of the expression to the storage
func (c *Client) SendResult(expression Expression) (bool, error) {
	data := sendResult{Expression: expression}
	body, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	post, err := http.Post(c.storageServer+"/postResult", "application/json", strings.NewReader(string(body)))
	if err != nil {
		return false, err
	}
	defer post.Body.Close()

	body, err = io.ReadAll(post.Body)
	if err != nil {
		return false, err
	}

	var ans AnsSendResult
	if err = json.Unmarshal(body, &ans); err != nil {
		return false, err
	}
	return ans.Message == "ok", nil
}

type AnsGetOperationsAndTimes struct {
	Data    map[string]int `json:"data"`
	Message string         `json:"message"`
}

// GetOperationsAndTimes returns the time for each operation from the storage
func (c *Client) GetOperationsAndTimes() (expression_parser.ExecTimeConfig, error) {
	var ans AnsGetOperationsAndTimes

	resp, err := http.Get(c.storageServer + "/getOperationsAndTimes")
	if err != nil {
		return expression_parser.ExecTimeConfig{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return expression_parser.ExecTimeConfig{}, err
	}

	if err = json.Unmarshal(body, &ans); err != nil {
		return expression_parser.ExecTimeConfig{}, err
	}

	var config expression_parser.ExecTimeConfig
	for key, val := range ans.Data {
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
			return expression_parser.ExecTimeConfig{}, fmt.Errorf("unknown operator: %s", key)
		}
	}

	return config, nil
}

type sendAlive struct {
	Expression Expression `json:"expression"`
}

func (c *Client) SendAlive(expression Expression) error {
	var send sendAlive
	send.Expression = expression
	body, err := json.Marshal(send)
	if err != nil {
		return err
	}

	_, err = http.Post(c.storageServer+"/keepAlive", "application/json", strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	return nil
}
