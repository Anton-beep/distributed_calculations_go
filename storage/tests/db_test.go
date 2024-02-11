package tests_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"storage/internal/db"
	"testing"
)

func TestAddGetExpression(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	lastID := d.GetLastID()

	newID, err := d.AddExpression(db.Expression{
		ID:     lastID + 1,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionReady,
	})

	require.NoError(t, err)

	expression, err := d.GetExpressionByID(newID)
	require.NoError(t, err)
	assert.Equal(t, "2 + 2", expression.Value)
	assert.InDelta(t, float64(4), expression.Answer, 0.0001)
	assert.Equal(t, "ok", expression.Logs)
	assert.Equal(t, db.ExpressionReady, expression.Status)

	err = d.DeleteExpression(newID)
	require.NoError(t, err)
}

func TestUpdateExpression(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	lastID := d.GetLastID()

	newID, err := d.AddExpression(db.Expression{
		ID:     lastID + 1,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionReady,
	})
	require.NoError(t, err)

	err = d.UpdateExpression(db.Expression{
		ID:     newID,
		Value:  "2 + 2",
		Answer: 5,
		Logs:   "ok",
		Status: db.ExpressionReady,
	})
	require.NoError(t, err)

	expression, err := d.GetExpressionByID(newID)
	require.NoError(t, err)
	assert.Equal(t, "2 + 2", expression.Value)
	assert.InDelta(t, float64(5), expression.Answer, 0.0001)
	assert.Equal(t, "ok", expression.Logs)
	assert.Equal(t, db.ExpressionReady, expression.Status)

	err = d.DeleteExpression(newID)
	assert.NoError(t, err)
}

func TestGetAllExpressions(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	lastID := d.GetLastID()

	newID1, err := d.AddExpression(db.Expression{
		ID:     lastID + 1,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionReady,
	})
	require.NoError(t, err)

	lastID = d.GetLastID()

	newID2, err := d.AddExpression(db.Expression{
		ID:     lastID + 1,
		Value:  "2 + 3",
		Answer: 5,
		Logs:   "ok",
		Status: db.ExpressionReady,
	})
	require.NoError(t, err)

	expressions, err := d.GetAllExpressions()
	require.NoError(t, err)
	for _, expression := range expressions {
		if expression.ID == newID1 {
			assert.Equal(t, "2 + 2", expression.Value)
			assert.InDelta(t, float64(4), expression.Answer, 0.0001)
			assert.Equal(t, "ok", expression.Logs)
			assert.Equal(t, db.ExpressionReady, expression.Status)
		}
		if expression.ID == newID2 {
			assert.Equal(t, "2 + 3", expression.Value)
			assert.InDelta(t, float64(5), expression.Answer, 0.0001)
			assert.Equal(t, "ok", expression.Logs)
			assert.Equal(t, db.ExpressionReady, expression.Status)
		}
	}

	err = d.DeleteExpression(newID1)
	require.NoError(t, err)
	err = d.DeleteExpression(newID2)
	require.NoError(t, err)
}

func TestDeleteExpression(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	lastID := d.GetLastID()

	newID, err := d.AddExpression(db.Expression{
		ID:     lastID + 1,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionReady,
	})
	require.NoError(t, err)

	err = d.DeleteExpression(newID)
	require.NoError(t, err)

	_, err = d.GetExpressionByID(newID)
	require.Error(t, err)
}
