package tests_test

import (
	"calculationServer/pkg/expressionparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		{"with brackets", "(2 + 2) + (2 + 2) + (2 + 2)",
			[]string{"2", "2", "+", "2", "2", "+", "2", "2", "+", "+", "+"}, false},
		{"float", "0.2 + 0.2", []string{"0.2", "0.2", "+"}, false},
	}

	ep := expressionparser.NewExpressionParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ep.ConvertInRPN(tt.in)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.out, actual)
			}
		})
	}
}

func TestReadRPN(t *testing.T) {
	type element struct {
		name      string
		in        []string
		out       []expressionparser.OperationOrNum
		wantError bool
	}

	tests := []element{
		{name: "simple", in: []string{"3", "4", "+"}, out: []expressionparser.OperationOrNum{
			{Data: 3},
			{Data: 4},
			{IsOperation: true, OperationID2: 1},
		}},
		{name: "complicated", in: []string{"3", "4", "2", "*", "1", "5", "-", "/", "+"},
			out: []expressionparser.OperationOrNum{
				{Data: 3},
				{Data: 4},
				{Data: 2},
				{IsOperation: true, OperationID1: 1, OperationID2: 2, Operator: 3},
				{Data: 1},
				{Data: 5},
				{IsOperation: true, OperationID1: 4, OperationID2: 5, Operator: 1},
				{IsOperation: true, OperationID1: 3, OperationID2: 6, Operator: 2},
				{IsOperation: true, OperationID2: 7},
			}},
		{name: "too many operators", in: []string{"1", "2", "+", "+"}, out: nil, wantError: true},
		{name: "too many numbers", in: []string{"1", "2", "3", "+"}, out: nil, wantError: true},
		{name: "unexpected symbol", in: []string{"1", "2", "+-"}, out: nil, wantError: true},
		{name: "unexpected symbol2", in: []string{"1", "2", "&"}, out: nil, wantError: true},
	}

	ep := expressionparser.NewExpressionParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ep.ReadRPN(tt.in)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.out, actual)
			}
		})
	}
}

func TestCalculation(t *testing.T) {
	type element struct {
		name      string
		in        []expressionparser.OperationOrNum
		out       float64
		wantError bool
	}
	tests := []element{
		{name: "simple", in: []expressionparser.OperationOrNum{
			{Data: 3},
			{Data: 4},
			{IsOperation: true, OperationID2: 1},
		}, out: 7},
	}

	ep := expressionparser.NewExpressionParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ep.CalculateRPNData(tt.in)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.InDelta(t, tt.out, actual, 0.001)
			}
		})
	}
}

func TestFullProcess(t *testing.T) {
	numberOfWorkers := 10
	timeCfg := expressionparser.ExecTimeConfig{
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

	ep := expressionparser.NewExpressionParser()
	err := ep.SetNumberOfWorkers(numberOfWorkers)
	require.NoError(t, err)
	err = ep.SetExecTimes(timeCfg)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			actual, _, err2 := ep.CalculateExpression(tt.in)
			if tt.wantError {
				require.Error(t, err2)
			} else {
				require.NoError(t, err2)
				assert.InDelta(t, tt.out, actual, 0.001)
			}
		})
	}
}
