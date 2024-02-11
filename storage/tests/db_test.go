package tests

import (
	"github.com/stretchr/testify/assert"
	"storage/internal/Db"
	"testing"
)

func TestAddGetExpression(t *testing.T) {
	d, err := Db.New()
	assert.NoError(t, err)

	lastId, err := d.GetLastId()

	newId, err := d.AddExpression(Db.Expression{
		Id:     lastId + 1,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: Db.ExpressionReady,
	})

	assert.NoError(t, err)

	expression, err := d.GetExpressionById(newId)
	assert.NoError(t, err)
	assert.Equal(t, "2 + 2", expression.Value)
	assert.Equal(t, float64(4), expression.Answer)
	assert.Equal(t, "ok", expression.Logs)
	assert.Equal(t, Db.ExpressionReady, expression.Status)

	err = d.DeleteExpression(newId)
	assert.NoError(t, err)
}

func TestUpdateExpression(t *testing.T) {
	d, err := Db.New()
	assert.NoError(t, err)

	lastId, err := d.GetLastId()

	newId, err := d.AddExpression(Db.Expression{
		Id:     lastId + 1,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: Db.ExpressionReady,
	})
	assert.NoError(t, err)

	err = d.UpdateExpression(Db.Expression{
		Id:     newId,
		Value:  "2 + 2",
		Answer: 5,
		Logs:   "ok",
		Status: Db.ExpressionReady,
	})
	assert.NoError(t, err)

	expression, err := d.GetExpressionById(newId)
	assert.NoError(t, err)
	assert.Equal(t, "2 + 2", expression.Value)
	assert.Equal(t, float64(5), expression.Answer)
	assert.Equal(t, "ok", expression.Logs)
	assert.Equal(t, Db.ExpressionReady, expression.Status)

	err = d.DeleteExpression(newId)
	assert.NoError(t, err)
}

func TestGetAllExpressions(t *testing.T) {
	d, err := Db.New()
	assert.NoError(t, err)

	lastId, err := d.GetLastId()

	newId1, err := d.AddExpression(Db.Expression{
		Id:     lastId + 1,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: Db.ExpressionReady,
	})
	assert.NoError(t, err)

	lastId, err = d.GetLastId()
	assert.NoError(t, err)

	newId2, err := d.AddExpression(Db.Expression{
		Id:     lastId + 1,
		Value:  "2 + 3",
		Answer: 5,
		Logs:   "ok",
		Status: Db.ExpressionReady,
	})
	assert.NoError(t, err)

	expressions, err := d.GetAllExpressions()
	assert.NoError(t, err)
	for _, expression := range expressions {
		if expression.Id == newId1 {
			assert.Equal(t, "2 + 2", expression.Value)
			assert.Equal(t, float64(4), expression.Answer)
			assert.Equal(t, "ok", expression.Logs)
			assert.Equal(t, Db.ExpressionReady, expression.Status)
		}
		if expression.Id == newId2 {
			assert.Equal(t, "2 + 3", expression.Value)
			assert.Equal(t, float64(5), expression.Answer)
			assert.Equal(t, "ok", expression.Logs)
			assert.Equal(t, Db.ExpressionReady, expression.Status)
		}
	}

	err = d.DeleteExpression(newId1)
	assert.NoError(t, err)
	err = d.DeleteExpression(newId2)
	assert.NoError(t, err)
}

func TestDeleteExpression(t *testing.T) {
	d, err := Db.New()
	assert.NoError(t, err)

	lastId, err := d.GetLastId()

	newId, err := d.AddExpression(Db.Expression{
		Id:     lastId + 1,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: Db.ExpressionReady,
	})
	assert.NoError(t, err)

	err = d.DeleteExpression(newId)
	assert.NoError(t, err)

	_, err = d.GetExpressionById(newId)
	assert.Error(t, err)
}
