package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"storage/internal/db"
	"testing"
)

func CreateTestUser(t *testing.T, d *db.APIDb) int {
	lastID := d.GetLastID()

	newID, err := d.AddUser(db.User{
		ID:    lastID + 1,
		Login: "test",
	})
	require.NoError(t, err)
	return newID
}

func TestAddGetExpression(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	newUser := CreateTestUser(t, d)

	lastID := d.GetLastID()

	newID, err := d.AddExpression(db.Expression{
		ID:     lastID + 1,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionReady,
		User:   newUser,
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
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestUpdateExpression(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	newUser := CreateTestUser(t, d)

	lastID := d.GetLastID()

	newID, err := d.AddExpression(db.Expression{
		ID:     lastID + 1,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionReady,
		User:   newUser,
	})
	require.NoError(t, err)

	err = d.UpdateExpression(db.Expression{
		ID:     newID,
		Value:  "2 + 2",
		Answer: 5,
		Logs:   "ok",
		Status: db.ExpressionReady,
		User:   newUser,
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
	err = d.DeleteUser(newUser)
	assert.NoError(t, err)
}

func TestGetAllExpressions(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	newUser := CreateTestUser(t, d)

	lastID := d.GetLastID()

	newID1, err := d.AddExpression(db.Expression{
		ID:     lastID + 1,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionReady,
		User:   newUser,
	})
	require.NoError(t, err)

	lastID = d.GetLastID()

	newID2, err := d.AddExpression(db.Expression{
		ID:     lastID + 1,
		Value:  "2 + 3",
		Answer: 5,
		Logs:   "ok",
		Status: db.ExpressionReady,
		User:   newUser,
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
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestDeleteExpression(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	newUser := CreateTestUser(t, d)

	lastID := d.GetLastID()

	newID, err := d.AddExpression(db.Expression{
		ID:     lastID + 1,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionReady,
		User:   newUser,
	})
	require.NoError(t, err)

	err = d.DeleteExpression(newID)
	require.NoError(t, err)

	_, err = d.GetExpressionByID(newID)
	require.Error(t, err)

	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestOperations(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	newUser := CreateTestUser(t, d)

	lastID := d.GetLastID()

	newID, err := d.AddOperation(db.Operation{
		ID:           lastID + 1,
		TimeAdd:      1,
		TimeSubtract: 2,
		TimeDivide:   3,
		TimeMultiply: 4,
		User:         newUser,
	})
	require.NoError(t, err)

	operation, err := d.GetUserOperations(newUser)
	require.NoError(t, err)
	assert.Equal(t, 1, operation.TimeAdd)
	assert.Equal(t, 2, operation.TimeSubtract)
	assert.Equal(t, 3, operation.TimeDivide)
	assert.Equal(t, 4, operation.TimeMultiply)

	err = d.UpdateOperation(db.Operation{
		ID:           newID,
		TimeAdd:      2,
		TimeSubtract: 3,
		TimeDivide:   4,
		TimeMultiply: 5,
	})
	require.NoError(t, err)
	operation, err = d.GetUserOperations(newUser)
	require.NoError(t, err)
	assert.Equal(t, 2, operation.TimeAdd)
	assert.Equal(t, 3, operation.TimeSubtract)
	assert.Equal(t, 4, operation.TimeDivide)
	assert.Equal(t, 5, operation.TimeMultiply)

	err = d.DeleteOperation(newID)
	require.NoError(t, err)
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestUsers(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	newUser := CreateTestUser(t, d)

	err = d.UpdateUser(db.User{
		ID:    newUser,
		Login: "test2",
	})
	require.NoError(t, err)

	user, err := d.GetUserByUsername("test2")
	require.Equal(t, newUser, user.ID)
	require.NoError(t, err)

	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}
