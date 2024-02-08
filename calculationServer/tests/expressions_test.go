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
				assert.NoError(t, err)
				assert.Equal(t, tt.out, actual)
			}
		})
	}
}

func TestReadRPN(t *testing.T) {
	type element struct {
		name      string
		in        []string
		out       map[int]ExpressionParser.OperationOrNum
		wantError bool
	}

	tests := []element{
		{name: "simple", in: []string{"3", "4", "+"}, out: map[int]ExpressionParser.OperationOrNum{
			0: {false, 0, 0, 0, 3},
			1: {false, 0, 0, 0, 4},
			2: {true, 0, 1, 0, 0},
		}},
		{name: "complicated", in: []string{"3", "4", "2", "*", "1", "5", "-", "/", "+"}, out: map[int]ExpressionParser.OperationOrNum{
			0: {false, 0, 0, 0, 3},
			1: {false, 0, 0, 0, 4},
			2: {false, 0, 0, 0, 2},
			3: {true, 1, 2, 3, 0},
			4: {false, 0, 0, 0, 1},
			5: {false, 0, 0, 0, 5},
			6: {true, 4, 5, 1, 0},
			7: {true, 3, 6, 2, 0},
			8: {true, 0, 7, 0, 0},
		}},
		{name: "too many operators", in: []string{"1", "2", "+", "+"}, out: nil, wantError: true},
		{name: "too many numbers", in: []string{"1", "2", "3", "+"}, out: nil, wantError: true},
		{name: "unexpected symbol", in: []string{"1", "2", "+-"}, out: nil, wantError: true},
		{name: "unexpected symbol2", in: []string{"1", "2", "&"}, out: nil, wantError: true},
	}

	ep := ExpressionParser.NewExpressionParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ep.ReadRPN(tt.in)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.out, actual)
			}
		})
	}
}
