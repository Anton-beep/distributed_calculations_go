package expressionparser

import (
	"calculationServer/internal/expressionlogger"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ADD = iota
	SUBTRACT
	DIVIDE
	MULTIPLY
)

type OperationOrNum struct {
	IsOperation  bool
	OperationID1 int
	OperationID2 int
	Operator     int
	Data         float64
}

type ExecTimeConfig struct {
	TimeAdd      time.Duration
	TimeSubtract time.Duration
	TimeDivide   time.Duration
	TimeMultiply time.Duration
}

type ExpressionParser struct {
	numberOfWorkers int
	execTimeConfig  ExecTimeConfig
	mu              sync.Mutex
	logs            *expressionlogger.ExpLogger
	running         int
}

func isByteNumberOrPoint(b byte) bool {
	if (b >= '0' && b <= '9') || b == '.' {
		return true
	}
	return false
}

func isOperatorGreater(b1 string, b2 string) bool {
	// if an Operator b2 has greater precedence than b1, or they have equal precedence and b1 is left associative
	if b2 == "(" {
		return false
	}

	if b1 == "+" || b1 == "-" {
		if b2 == "*" || b2 == "/" {
			return true
		}

		if b1 == "-" {
			return true
		}
	}
	if b1 == "/" && (b2 == "*" || b2 == "/") {
		return true
	}
	return false
}

func convertByteToOperator(b byte) (int, error) {
	switch b {
	case '+':
		return ADD, nil
	case '-':
		return SUBTRACT, nil
	case '*':
		return MULTIPLY, nil
	case '/':
		return DIVIDE, nil
	}
	return 0, errors.New("not an Operator")
}

func convertOperatorToString(oper int) (string, error) {
	switch oper {
	case ADD:
		return "+", nil
	case SUBTRACT:
		return "-", nil
	case DIVIDE:
		return "/", nil
	case MULTIPLY:
		return "*", nil
	default:
		return "", errors.New("no such operator")
	}
}

func IsExecTimeConfigCorrect(execTimeConfig ExecTimeConfig) (bool, error) {
	if execTimeConfig.TimeAdd < 0 || execTimeConfig.TimeSubtract < 0 || execTimeConfig.TimeDivide < 0 ||
		execTimeConfig.TimeMultiply < 0 {
		return false, errors.New("execution time cannot be smaller than 0")
	}
	return true, nil
}

func New() *ExpressionParser {
	return &ExpressionParser{
		numberOfWorkers: 1,
		logs:            expressionlogger.New(),
	}
}

func (e *ExpressionParser) getTimeForOperator(operator int) (time.Duration, error) {
	switch operator {
	case ADD:
		return e.execTimeConfig.TimeAdd, nil
	case SUBTRACT:
		return e.execTimeConfig.TimeSubtract, nil
	case MULTIPLY:
		return e.execTimeConfig.TimeMultiply, nil
	case DIVIDE:
		return e.execTimeConfig.TimeDivide, nil
	default:
		return 0, errors.New("not an operator")
	}
}

func (e *ExpressionParser) SetExecTimes(execTimeConfig ExecTimeConfig) error {
	if _, err := IsExecTimeConfigCorrect(execTimeConfig); err != nil {
		return err
	}
	e.execTimeConfig = execTimeConfig
	return nil
}

func (e *ExpressionParser) SetNumberOfWorkers(in int) error {
	if in < 1 {
		return errors.New("number of workers must be bigger than 0")
	}
	e.numberOfWorkers = in
	return nil
}

// ConvertInRPN see an animation https://somethingorotherwhatever.com/shunting-yard-animation/
func (e *ExpressionParser) ConvertInRPN(expression string) ([]string, error) {
	e.logs.Add("Start conversion to reversed polish notation")

	stack := make([]string, 0)
	out := make([]string, 0)

	lastNum := false
	lastOper := true

	switchSign := false

	for i := 0; i < len(expression); i++ {
		if isByteNumberOrPoint(expression[i]) {
			lastOper = false

			if lastNum {
				return nil, fmt.Errorf("unexpected number %v, pos: %v", string(expression[i]), i)
			}
			bufNum := make([]byte, 0)
			if switchSign {
				bufNum = append(bufNum, '-')
				switchSign = false
			}
			for i < len(expression) && isByteNumberOrPoint(expression[i]) {
				// while a number, push it to the output
				bufNum = append(bufNum, expression[i])
				i++
			}

			out = append(out, string(bufNum))
			lastNum = true
			if i >= len(expression) {
				break
			}
		}

		if _, err := convertByteToOperator(expression[i]); err == nil {
			// an Operator is found
			if lastOper {
				if expression[i] == '+' {
					// unnecessary plus before the number
					continue
				}
				if expression[i] == '-' {
					// minus before the number
					switchSign = !switchSign
					continue
				}
				return nil, fmt.Errorf("unexpected operator %v, pos: %v", string(expression[i]), i)
			}
			lastNum = false

			if len(stack) > 0 {
				// While there is an Operator o₂ at the top of the stack with greater precedence,
				// or with equal precedence and o₁ is left associative, push o₂ from the stack to the output.
				for len(stack) > 0 && isOperatorGreater(string(expression[i]), stack[len(stack)-1]) {
					out = append(out, stack[len(stack)-1])
					stack = stack[:len(stack)-1]
				}
			}
			stack = append(stack, string(expression[i]))

			lastOper = true
			continue
		}

		if string(expression[i]) == "(" {
			stack = append(stack, string(expression[i]))
			continue
		}

		if string(expression[i]) == ")" {
			for stack[len(stack)-1] != "(" {
				out = append(out, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1]
			continue
		}

		if string(expression[i]) == " " {
			continue
		}

		return nil, fmt.Errorf("unexpected symbol %v, pos: %v", string(expression[i]), i)
	}

	for len(stack) > 0 {
		out = append(out, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	e.logs.Add("Result: " + strings.Join(out, " "))
	return out, nil
}

// ReadRPN read reversed polish notation and convert to slice, ao it can be calculated later.
func (e *ExpressionParser) ReadRPN(expressionRPN []string) ([]OperationOrNum, error) {
	stack := make([]int, 0)
	data := make([]OperationOrNum, len(expressionRPN))
	id := 0

	for ind, el := range expressionRPN {
		if val, err := strconv.ParseFloat(el, 64); err == nil {
			data[id] = OperationOrNum{Data: val}
			stack = append(stack, id)
			id++
			continue
		}
		if len(el) > 1 {
			return nil, fmt.Errorf("unexpected symbol %v, pos: %v", el, ind)
		}
		if oper, err := convertByteToOperator(el[0]); err == nil {
			if len(stack) < 2 {
				return nil, fmt.Errorf("not enought arguments for Operator, pos: %v (need 2 number for an operator)", ind)
			}
			data[id] = OperationOrNum{
				IsOperation:  true,
				OperationID1: stack[len(stack)-2],
				OperationID2: stack[len(stack)-1],
				Operator:     oper,
			}
			stack = stack[:len(stack)-2]
			stack = append(stack, id)
			id++
			continue
		}
		return nil, fmt.Errorf("unexpected symbol %v, pos: %v", el, ind)
	}
	if len(stack) > 1 {
		return nil, errors.New("unexpected numbers")
	}
	return data, nil
}

func (e *ExpressionParser) CalculateOperation(num1, num2 float64, operator int) (float64, error) {
	e.mu.Lock()
	duration, err := e.getTimeForOperator(operator)
	e.mu.Unlock()
	if err != nil {
		return 0, fmt.Errorf("%v is not an operator", operator)
	}

	res := make(chan float64)
	switch operator {
	case ADD:
		go func() {
			res <- num1 + num2
		}()
		time.Sleep(duration)

		return <-res, nil
	case SUBTRACT:
		go func() {
			res <- num1 - num2
		}()
		time.Sleep(duration)

		return <-res, nil
	case DIVIDE:
		if num2 == 0 {
			return 0, errors.New("division by zero")
		}

		go func() {
			res <- num1 / num2
		}()
		time.Sleep(duration)

		return <-res, nil
	case MULTIPLY:
		go func() {
			res <- num1 * num2
		}()
		time.Sleep(duration)

		return <-res, nil
	}

	return 0, errors.New("unknown operator")
}

// CalculateRPNData aka workerPool.
func (e *ExpressionParser) CalculateRPNData(data []OperationOrNum) (float64, error) {
	// pool will control number of workers at the same time
	e.logs.Add("Start of calculations")
	if e.numberOfWorkers < 1 {
		return 0, errors.New("number of workers must be bigger than 0")
	}

	errChan := make(chan error)
	e.running = 0
	readyChan := make(chan bool, e.numberOfWorkers+1)
	// to understand, when this goroutine will go through all data elements

	for ind, el := range data {
		el := el
		ind := ind
		if !el.IsOperation {
			continue
		}
		for {
			e.mu.Lock()
			if !data[el.OperationID1].IsOperation && !data[el.OperationID2].IsOperation && e.running < e.numberOfWorkers {
				// we can start new worker
				e.mu.Unlock()
				break
			}
			e.mu.Unlock()
			select {
			case <-readyChan:
				// wait one worker to be done
				continue
			case err := <-errChan:
				return 0, err
			}
		}

		e.running++
		go func() {
			e.mu.Lock()
			// take numbers from data
			num1, num2 := data[el.OperationID1].Data, data[el.OperationID2].Data
			e.mu.Unlock()

			strOper, err := convertOperatorToString(el.Operator)
			// if error, write in a channel
			if err != nil {
				errChan <- err
				return
			}
			e.logs.Add(fmt.Sprintf("Start worker with id %v; work: %v %v %v", ind, num1, strOper, num2))

			// calculate with delays
			outOper, err := e.CalculateOperation(num1, num2, el.Operator)

			e.logs.Add(fmt.Sprintf("End of worker with id %v; work was %v %v %v; result is %v",
				ind, num1, strOper, num2, outOper))
			// if error, write in a channel
			if err != nil {
				errChan <- err
				return
			}
			e.mu.Lock()
			// write result of an operation
			data[ind] = OperationOrNum{Data: outOper}
			e.mu.Unlock()

			e.mu.Lock()
			e.running--
			e.mu.Unlock()
			// read one element from pool, so new goroutine can turn on
			readyChan <- true
		}()
	}

	for e.running > 0 {
		select {
		case <-readyChan:
			continue
		case err := <-errChan:
			return 0, err
		}
	}

	e.logs.Add(fmt.Sprintf("All workers are stopped; the final result is %v", data[len(data)-1].Data))

	return data[len(data)-1].Data, nil
}

func (e *ExpressionParser) CalculateExpression(in string) (float64, string, error) {
	e.logs.Reset()
	// convert to RPN
	rpn, err := e.ConvertInRPN(in)
	if err != nil {
		return 0, "", err
	}
	// read rpn, setup for calculator
	data, err := e.ReadRPN(rpn)
	if err != nil {
		return 0, "", err
	}
	// calculate
	res, err := e.CalculateRPNData(data)
	if err != nil {
		return 0, "", err
	}
	return res, e.logs.Get(), nil
}

func (e *ExpressionParser) GetWorkingWorkers() int {
	return e.running
}

func (e *ExpressionParser) GetTotalNumberOfWorkers() int {
	return e.numberOfWorkers
}
