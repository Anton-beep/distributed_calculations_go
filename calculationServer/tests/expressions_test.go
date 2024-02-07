package tests

import (
	"calculationServer/pkg/ExpressionParser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertToRPN(t *testing.T) {
	type element struct {
		name      string
		in        string
		out       []string
		wantError bool
	}
	tests := []element{
		{"simple", "3+4", []string{"3", "4", "+"}, false},
		{"complicated", "3 + 4 * 2 / (1 - 5)", []string{"3", "4", "2", "*", "1", "5", "-", "/", "+"}, false},
		{"double brackets", "2 * ((1 + 1) + 1)", []string{"2", "1", "1", "+", "1", "+", "*"}, false},
		{"wrong (unexpected symbol)", "3 + 4 * 2 / (1s - 5)", nil, true},
		{"wrong (two operators in a row)", "3 + 4 * 2 / (1 - - 5)", nil, true},
		{"wrong (two numbers in a row)", "3 + 4 * 2 2 / (1 - 5)", nil, true},
		{"extra brackets", "((1 + 2))", []string{"1", "2", "+"}, false},
		{"operator from stack (o1 is +, o2 is *)", "4 * 3 + 2", []string{"4", "3", "*", "2", "+"}, false},
		{"operator from stack (o1 is -, o2 is -)", "4 - 3 - 2", []string{"4", "3", "-", "2", "-"}, false},
	}

	ep := ExpressionParser.NewExpressionParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ep.ConvertExpressionInRPN(tt.in)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.out, actual)
			}
		})
	}
}
