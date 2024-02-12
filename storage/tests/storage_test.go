package tests_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"storage/internal/db"
	"storage/internal/expression_storage"
	"testing"
)

func TestAddGetStorage(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)
	e := expression_storage.New(d)

	newID, err := e.Add(db.Expression{
		ID:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionReady,
	})
	require.NoError(t, err)

	expression, err := e.GetByID(newID)
	require.NoError(t, err)

	assert.Equal(t, "2 + 2", expression.Value)
	assert.InDelta(t, float64(4), expression.Answer, 0.0001)
	assert.Equal(t, "ok", expression.Logs)
	assert.Equal(t, db.ExpressionReady, expression.Status)

	err = e.Delete(newID)
	require.NoError(t, err)
}

func TestGetAllStorage(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)
	e := expression_storage.New(d)

	newID1, err := e.Add(db.Expression{
		ID:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionReady,
	})
	require.NoError(t, err)

	newID2, err := e.Add(db.Expression{
		ID:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionReady,
	})
	require.NoError(t, err)

	expressions := e.GetAll()
	assert.NotEmpty(t, expressions)

	for _, expression := range expressions {
		if expression.ID == newID1 {
			assert.Equal(t, "2 + 2", expression.Value)
			assert.InDelta(t, float64(4), expression.Answer, 0.0001)
			assert.Equal(t, "ok", expression.Logs)
			assert.Equal(t, db.ExpressionReady, expression.Status)
		}
		if expression.ID == newID2 {
			assert.Equal(t, "2 + 2", expression.Value)
			assert.InDelta(t, float64(4), expression.Answer, 0.0001)
			assert.Equal(t, "ok", expression.Logs)
			assert.Equal(t, db.ExpressionReady, expression.Status)
		}
	}

	err = e.Delete(newID1)
	require.NoError(t, err)
	err = e.Delete(newID2)
	require.NoError(t, err)
}

func TestUpdatePending(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)
	e := expression_storage.New(d)

	newID, err := e.Add(db.Expression{
		ID:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionNotReady,
	})
	require.NoError(t, err)

	expression, err := e.GetByID(newID)
	require.NoError(t, err)

	expression.Status = db.ExpressionReady
	err = e.UpdatePendingExpression(expression)
	require.NoError(t, err)

	expression, err = e.GetByID(newID)
	require.NoError(t, err)
	assert.Equal(t, "2 + 2", expression.Value)
	assert.InDelta(t, float64(4), expression.Answer, 0.0001)
	assert.Equal(t, "ok", expression.Logs)
	assert.Equal(t, db.ExpressionReady, expression.Status)

	err = e.Delete(newID)
	require.NoError(t, err)
}

func TestPendingToReady(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)
	e := expression_storage.New(d)

	expression := db.Expression{
		ID:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionNotReady,
	}
	newID, err := e.Add(expression)
	require.NoError(t, err)

	expression.Status = db.ExpressionReady
	expression.ID = newID
	err = e.PendingToReady(expression)
	require.NoError(t, err)

	expression, err = e.GetByID(newID)
	require.NoError(t, err)
	assert.Equal(t, "2 + 2", expression.Value)
	assert.InDelta(t, float64(4), expression.Answer, 0.0001)
	assert.Equal(t, "ok", expression.Logs)
	assert.Equal(t, db.ExpressionReady, expression.Status)

	err = e.Delete(newID)
	require.NoError(t, err)
}
