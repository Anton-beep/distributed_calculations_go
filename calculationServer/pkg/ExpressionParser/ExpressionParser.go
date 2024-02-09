package ExpressionParser

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type OperationOrNum struct {
	IsOperation  bool
	OperationId1 int
	OperationId2 int
	Operator     int
	Data         int
}

type ExecTimeConfig struct {
	TimeAdd      time.Duration
	TimeSubtract time.Duration
	TimeDivide   time.Duration
	TimeMultiply time.Duration
}

type ExpressionParser struct {
	data            map[int]OperationOrNum // here will be stored simple functions, which will return a simple number and operations calculate, both must return channel to read int from
	amountOfWorkers int
	execTimeConfig  ExecTimeConfig
	lock            sync.RWMutex
}

func NewExpressionParser() *ExpressionParser {
	return &ExpressionParser{
		amountOfWorkers: 1,
	}
}

func IsExecTimeConfigCorrect(execTimeConfig ExecTimeConfig) (bool, error) {
	if execTimeConfig.TimeAdd < 0 || execTimeConfig.TimeSubtract < 0 || execTimeConfig.TimeDivide < 0 ||
		execTimeConfig.TimeMultiply < 0 {
		return false, errors.New("execution time cannot be smaller than 0")
	}
	return true, nil
}

func (e *ExpressionParser) SetExecTimes(execTimeConfig ExecTimeConfig) error {
	if _, err := IsExecTimeConfigCorrect(execTimeConfig); err != nil {
		return err
	}
	e.execTimeConfig = execTimeConfig
	return nil
}

func isByteNumber(b byte) bool {
	if b > []byte("0")[0] && b < []byte("9")[0] {
		return true
	}
	return false
}

func convertByteToOperator(b byte) (int, error) {
	switch b {
	case []byte("+")[0]:
		return ADD, nil
	case []byte("-")[0]:
		return SUBTRACT, nil
	case []byte("*")[0]:
		return MULTIPLY, nil
	case []byte("/")[0]:
		return DIVIDE, nil
	}
	return 0, errors.New("not an Operator")
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

// ConvertExpressionInRPN see an animation https://somethingorotherwhatever.com/shunting-yard-animation/
func (e *ExpressionParser) ConvertExpressionInRPN(expression string) ([]string, error) {
	stack := make([]string, 0)
	out := make([]string, 0)

	lastNum := false
	lastOper := true

	switchSign := false

	for i := 0; i < len(expression); i++ {
		if isByteNumber(expression[i]) {
			lastOper = false

			if lastNum {
				return nil, fmt.Errorf("unexpected number %v, pos: %v", string(expression[i]), i)
			}
			bufNum := make([]byte, 0)
			if switchSign {
				bufNum = append(bufNum, '-')
				switchSign = false
			}
			for i < len(expression) && isByteNumber(expression[i]) {
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
				// While there is an Operator o₂ at the top of the stack with greater precedence, or with equal precedence and o₁ is left associative, push o₂ from the stack to the output
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

	return out, nil
}

func (e *ExpressionParser) ReadRPN(expressionRPN []string) ([]OperationOrNum, error) {
	stack := make([]int, 0)
	data := make([]OperationOrNum, len(expressionRPN))
	id := 0

	for ind, el := range expressionRPN {
		if val, err := strconv.Atoi(el); err == nil {
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
				return nil, fmt.Errorf("not enought arguments for Operator, pos: %v", ind)
			}
			data[id] = OperationOrNum{
				IsOperation:  true,
				OperationId1: stack[len(stack)-2],
				OperationId2: stack[len(stack)-1],
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

func (e *ExpressionParser) CalculateRPNData(data []OperationOrNum) (int, error) {
	return e.calculatorOrganizer(data)
	/*
		select {
		case calcRes := <-ansChan:
			return calcRes, nil
		case err := <-errChan:
			// control errors
			return 0, err
		}
	*/
}

func (e *ExpressionParser) isOperator(operator int) bool {
	return operator == ADD || operator == SUBTRACT || operator == MULTIPLY || operator == DIVIDE
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

func (e *ExpressionParser) calculateOperation(num1, num2, operator int) (int, error) {
	if !e.isOperator(operator) {
		return 0, fmt.Errorf("%v is not an operator", operator)
	}

	e.lock.RLock()
	duration, _ := e.getTimeForOperator(operator)
	e.lock.RUnlock()

	switch operator {
	case ADD:
		time.Sleep(duration)

		return num1 + num2, nil
	case SUBTRACT:
		time.Sleep(duration)

		return num1 - num2, nil
	case DIVIDE:
		time.Sleep(duration)

		if num2 == 0 {
			return 0, errors.New("division by zero")
		}
		return num1 / num2, nil
	case MULTIPLY:
		time.Sleep(duration)

		return num1 * num2, nil
	}

	return 0, errors.New("unknown operator")
}

// aka workerPool
func (e *ExpressionParser) calculatorOrganizer(data []OperationOrNum) (int, error) {
	// pool will control number of workers at the same time
	if e.amountOfWorkers < 1 {
		return 0, errors.New("number of workers must be bigger than 0")
	}

	errChan := make(chan error)

	mu := sync.Mutex{}
	running := 0
	readyChan := make(chan bool, e.amountOfWorkers*2)
	// to understand, when this goroutine will go through all data elements

	for ind, el := range data {
		el := el
		ind := ind
		if !el.IsOperation {
			continue
		} else {

			for {
				mu.Lock()
				fmt.Println(data, el, running, e.amountOfWorkers)
				if !data[el.OperationId1].IsOperation && !data[el.OperationId2].IsOperation && running < e.amountOfWorkers {
					// we can start new worker
					mu.Unlock()
					break
				}
				mu.Unlock()
				select {
				case <-readyChan:
					// wait one worker to be done
					continue
				case err := <-errChan:
					return 0, err
				}
			}

			running++
			go func() {
				fmt.Println("start operation", ind, data[el.OperationId1].IsOperation, data[el.OperationId2].IsOperation)
				// write something in pool, so max number of working goroutines will be equal to e.amountOfWorkers
				fmt.Printf("operation %v start \n", ind)

				mu.Lock()
				// take numbers from data
				num1, num2 := data[el.OperationId1].Data, data[el.OperationId2].Data
				mu.Unlock()
				// calculate with delays
				outOper, err := e.calculateOperation(num1, num2, el.Operator)
				fmt.Println("operation", num1, num2, el.Operator, outOper)
				// if error, write in a channel
				if err != nil {
					errChan <- err
					return
				}
				mu.Lock()
				// write result of an operation
				data[ind] = OperationOrNum{Data: outOper}
				mu.Unlock()

				fmt.Printf("operation %v end \n", ind)

				mu.Lock()
				running--
				mu.Unlock()
				// read one element from pool, so new goroutine can turn on
				readyChan <- true
			}()
		}
	}

	for running > 0 {
		select {
		case <-readyChan:
			continue
		case err := <-errChan:
			return 0, err
		}
	}

	fmt.Println("ans:", data[len(data)-1].Data, len(data)-1)
	return data[len(data)-1].Data, nil
}

func (e *ExpressionParser) CalculateExpression(in string) (int, error) {
	// convert to RPN
	rpn, err := e.ConvertExpressionInRPN(in)
	if err != nil {
		return 0, err
	}
	// read rpn, setup for calculator
	data, err := e.ReadRPN(rpn)
	if err != nil {
		return 0, err
	}
	// calculate
	res, err := e.CalculateRPNData(data)
	if err != nil {
		return 0, err
	}
	return res, nil
}
