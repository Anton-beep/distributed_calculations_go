package ExpressionParser

import (
	"errors"
	"fmt"
)

type ExpressionParser struct {
	data            map[int]func() // here will be stored simple functions, which will return a simple number and operations calculate, both must return channel to read int from
	amountOfWorkers int
}

func NewExpressionParser() *ExpressionParser {
	return &ExpressionParser{}
}

func isByteAllowed(b byte) bool {
	if isByteNumber(b) {
		return true
	}
	_, err := convertByteToOperator(b)
	if err == nil {
		return true
	}
	if b == []byte("(")[0] || b == []byte(")")[0] {
		return true
	}
	return false
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
	return 0, errors.New("not an operator")
}

func isOperatorGreater(b1 string, b2 string) bool {
	// if an operator b2 has greater precedence than b1, or they have equal precedence and b1 is left associative
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
			// an operator is found
			if lastOper {
				return nil, fmt.Errorf("unexpected operator %v, pos: %v", string(expression[i]), i)
			}
			lastNum = false

			if len(stack) > 0 {
				// While there is an operator o₂ at the top of the stack with greater precedence, or with equal precedence and o₁ is left associative, push o₂ from the stack to the output
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

func (e *ExpressionParser) ReadRPN(expressionRPN []string) (map[int]func(), error) {
	return nil, nil
}

func (e *ExpressionParser) Calculator(data map[int]func()) int {
	return 0
}
