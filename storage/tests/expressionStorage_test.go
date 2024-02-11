package tests

import (
	"github.com/stretchr/testify/assert"
	"storage/internal/Db"
	"storage/internal/ExpressionStorage"
	"testing"
)

func TestAddGetStorage(t *testing.T) {
	d, err := Db.New()
	assert.NoError(t, err)
	e := ExpressionStorage.New(d)

	newId, err := e.Add(Db.Expression{
		Id:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: Db.ExpressionReady,
	})
	assert.NoError(t, err)

	expression, err := e.GetById(newId)
	assert.NoError(t, err)

	assert.Equal(t, "2 + 2", expression.Value)
	assert.Equal(t, float64(4), expression.Answer)
	assert.Equal(t, "ok", expression.Logs)
	assert.Equal(t, Db.ExpressionReady, expression.Status)

	err = e.Delete(newId)
	assert.NoError(t, err)
}

func TestGetAllStorage(t *testing.T) {
	d, err := Db.New()
	assert.NoError(t, err)
	e := ExpressionStorage.New(d)

	newId1, err := e.Add(Db.Expression{
		Id:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: Db.ExpressionReady,
	})
	assert.NoError(t, err)

	newId2, err := e.Add(Db.Expression{
		Id:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: Db.ExpressionReady,
	})
	assert.NoError(t, err)

	expressions := e.GetAll()
	assert.NotEmpty(t, expressions)

	for _, expression := range expressions {
		if expression.Id == newId1 {
			assert.Equal(t, "2 + 2", expression.Value)
			assert.Equal(t, float64(4), expression.Answer)
			assert.Equal(t, "ok", expression.Logs)
			assert.Equal(t, Db.ExpressionReady, expression.Status)
		}
		if expression.Id == newId2 {
			assert.Equal(t, "2 + 2", expression.Value)
			assert.Equal(t, float64(4), expression.Answer)
			assert.Equal(t, "ok", expression.Logs)
			assert.Equal(t, Db.ExpressionReady, expression.Status)
		}
	}

	err = e.Delete(newId1)
	assert.NoError(t, err)
	err = e.Delete(newId2)
	assert.NoError(t, err)
}

func TestUpdatePending(t *testing.T) {
	d, err := Db.New()
	assert.NoError(t, err)
	e := ExpressionStorage.New(d)

	newId, err := e.Add(Db.Expression{
		Id:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: Db.ExpressionNotReady,
	})
	assert.NoError(t, err)

	expression, err := e.GetById(newId)

	expression.Status = Db.ExpressionReady
	err = e.UpdatePendingExpression(expression)
	assert.NoError(t, err)

	expression, err = e.GetById(newId)
	assert.NoError(t, err)
	assert.Equal(t, "2 + 2", expression.Value)
	assert.Equal(t, float64(4), expression.Answer)
	assert.Equal(t, "ok", expression.Logs)
	assert.Equal(t, Db.ExpressionReady, expression.Status)

	err = e.Delete(newId)
	assert.NoError(t, err)
}

func TestPendingToReady(t *testing.T) {
	d, err := Db.New()
	assert.NoError(t, err)
	e := ExpressionStorage.New(d)

	expression := Db.Expression{
		Id:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: Db.ExpressionNotReady,
	}
	newId, err := e.Add(expression)
	assert.NoError(t, err)

	expression.Status = Db.ExpressionReady
	expression.Id = newId
	err = e.PendingToReady(expression)
	assert.NoError(t, err)

	expression, err = e.GetById(newId)
	assert.NoError(t, err)
	assert.Equal(t, "2 + 2", expression.Value)
	assert.Equal(t, float64(4), expression.Answer)
	assert.Equal(t, "ok", expression.Logs)
	assert.Equal(t, Db.ExpressionReady, expression.Status)

	err = e.Delete(newId)
	assert.NoError(t, err)
}
