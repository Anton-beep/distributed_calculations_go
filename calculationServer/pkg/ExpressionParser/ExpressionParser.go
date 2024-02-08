package ExpressionParser

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type OperationOrNum struct {
	IsOperation  bool
	OperationId1 int
	OperationId2 int
	Operator     int
	Data         int
}

type ExpressionParser struct {
	data             map[int]OperationOrNum // here will be stored simple functions, which will return a simple number and operations calculate, both must return channel to read int from
	amountOfWorkers  int
	execTimeAdd      time.Duration
	execTimeSubtract time.Duration
	execTimeDivide   time.Duration
	execTimeMultiply time.Duration
}

func NewExpressionParser() *ExpressionParser {
	return &ExpressionParser{}
}

func (e *ExpressionParser) SetExecTimes(execTimeAdd, execTimeSubtract, execTimeDivide, execTimeMultiply time.Duration) error {
	if execTimeAdd < 0 {
		return errors.New("execution time cannot be smaller than 0")
	}
	if execTimeSubtract < 0 {
		return errors.New("execution time cannot be smaller than 0")
	}
	if execTimeDivide < 0 {
		return errors.New("execution time cannot be smaller than 0")
	}
	if execTimeMultiply < 0 {
		return errors.New("execution time cannot be smaller than 0")
	}
	e.execTimeAdd = execTimeAdd
	e.execTimeSubtract = execTimeSubtract
	e.execTimeDivide = execTimeDivide
	e.execTimeMultiply = execTimeMultiply
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
	lastOper := false

	for i := 0; i < len(expression); i++ {
		if isByteNumber(expression[i]) {
			lastOper = false

			if lastNum {
				return nil, fmt.Errorf("unexpected number %v, pos: %v", string(expression[i]), i)
			}
			bufNum := make([]byte, 0)
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
				return nil, fmt.Errorf("unexpected Operator %v, pos: %v", string(expression[i]), i)
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

func (e *ExpressionParser) ReadRPN(expressionRPN []string) (map[int]OperationOrNum, error) {
	stack := make([]int, 0)
	data := make(map[int]OperationOrNum)
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

func (e *ExpressionParser) Calculator(data map[int]func()) int {
	return 0
}
