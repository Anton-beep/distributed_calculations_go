package tests

import (
	"calculationServer/pkg/ExpressionParser"
	"fmt"
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
		{"strange notation (two operators - in a row)", "1 - -1", []string{"1", "-1", "-"}, false},
		{"wrong (two numbers in a row)", "3 + 4 * 2 2 / (1 - 5)", nil, true},
		{"extra brackets", "((1 + 2))", []string{"1", "2", "+"}, false},
		{"operator from stack (o1 is +, o2 is *)", "4 * 3 + 2", []string{"4", "3", "*", "2", "+"}, false},
		{"operator from stack (o1 is -, o2 is -)", "4 - 3 - 2", []string{"4", "3", "-", "2", "-"}, false},
		{"without brackets", "2 + 2 + 2 + 2 + 2 + 2", []string{"2", "2", "2", "2", "2", "2", "+", "+", "+", "+", "+"}, false},
		{"with brackets", "(2 + 2) + (2 + 2) + (2 + 2)", []string{"2", "2", "+", "2", "2", "+", "2", "2", "+", "+", "+"}, false},
		{"float", "0.2 + 0.2", []string{"0.2", "0.2", "+"}, false},
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
		out       []ExpressionParser.OperationOrNum
		wantError bool
	}

	tests := []element{
		{name: "simple", in: []string{"3", "4", "+"}, out: []ExpressionParser.OperationOrNum{
			{false, 0, 0, 0, 3},
			{false, 0, 0, 0, 4},
			{true, 0, 1, 0, 0},
		}},
		{name: "complicated", in: []string{"3", "4", "2", "*", "1", "5", "-", "/", "+"}, out: []ExpressionParser.OperationOrNum{
			{false, 0, 0, 0, 3},
			{false, 0, 0, 0, 4},
			{false, 0, 0, 0, 2},
			{true, 1, 2, 3, 0},
			{false, 0, 0, 0, 1},
			{false, 0, 0, 0, 5},
			{true, 4, 5, 1, 0},
			{true, 3, 6, 2, 0},
			{true, 0, 7, 0, 0},
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

func TestCalculation(t *testing.T) {
	type element struct {
		name      string
		in        []ExpressionParser.OperationOrNum
		out       float64
		wantError bool
	}
	tests := []element{
		{"simple", []ExpressionParser.OperationOrNum{
			{false, 0, 0, 0, 3},
			{false, 0, 0, 0, 4},
			{true, 0, 1, 0, 0},
		}, 7, false},
	}

	ep := ExpressionParser.NewExpressionParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ep.CalculateRPNData(tt.in)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.out, actual)
			}
		})
	}
}

func TestFullProcess(t *testing.T) {
	numberOfWorkers := 10
	timeCfg := ExpressionParser.ExecTimeConfig{
		TimeAdd:      50,
		TimeSubtract: 50,
		TimeDivide:   50,
		TimeMultiply: 50,
	}

	type element struct {
		in        string
		out       float64
		wantError bool
	}
	tests := []element{
		{"1 + 1", 2, false},
		{"2 + 2", 4, false},
		{"2 + 2 * 2", 6, false},
		{"2 * 2 * 2", 8, false},
		{"(2 + 2) * 2", 8, false},
		{"3 + 4 * 2 / (1 - 5)", 1, false},
		{"1", 1, false},
		{"-1", -1, false},
		{"2 * (-1)", -2, false},
		{"2 2 - 2", 0, true},
		{"1 - + 1", 0, false},
		{"1 + ()", 0, true},
		{"1 + (1)", 2, false},
		{"(2 + 2) + (2 + 2) + (2 + 2)", 12, false},
		{"2 + 2 + 2 + 2 + 2 + 2 + 2 + 2", 16, false},
		{"0.1 + 0.9", 1, false},
	}

	ep := ExpressionParser.NewExpressionParser()
	err := ep.SetNumberOfWorkers(numberOfWorkers)
	assert.NoError(t, err)
	err = ep.SetExecTimes(timeCfg)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			actual, logs, err := ep.CalculateExpression(tt.in)
			fmt.Println(logs)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.out, actual)
			}
		})
	}
}
