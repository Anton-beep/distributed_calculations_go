package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"storage/internal/db"
	"storage/internal/expressionstorage"
	"sync"
	"testing"
	"time"
)

func TestAddGetStorage(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)
	e := expressionstorage.New(d, 1, &sync.Map{})

	newUser := CreateTestUser(t, d)

	newID, err := e.Add(db.Expression{
		ID:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionReady,
		User:   newUser,
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
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestGetAllStorage(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)
	e := expressionstorage.New(d, 1, &sync.Map{})

	newUser := CreateTestUser(t, d)

	newID1, err := e.Add(db.Expression{
		ID:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionReady,
		User:   newUser,
	})
	require.NoError(t, err)

	newID2, err := e.Add(db.Expression{
		ID:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionReady,
		User:   newUser,
	})
	require.NoError(t, err)

	expressions := e.GetAll(newUser)
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
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestUpdatePending(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)
	e := expressionstorage.New(d, 1, &sync.Map{})

	newUser := CreateTestUser(t, d)

	newID, err := e.Add(db.Expression{
		ID:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionNotReady,
		User:   newUser,
	})
	require.NoError(t, err)

	expression, err := e.GetByID(newID)
	require.NoError(t, err)

	expression.Status = db.ExpressionReady
	err = e.UpdateExpression(expression)
	require.NoError(t, err)

	expression, err = e.GetByID(newID)
	require.NoError(t, err)
	assert.Equal(t, "2 + 2", expression.Value)
	assert.InDelta(t, float64(4), expression.Answer, 0.0001)
	assert.Equal(t, "ok", expression.Logs)
	assert.Equal(t, db.ExpressionReady, expression.Status)

	err = e.Delete(newID)
	require.NoError(t, err)
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestPendingToReady(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)
	e := expressionstorage.New(d, 1, &sync.Map{})

	newUser := CreateTestUser(t, d)

	expression := db.Expression{
		ID:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionNotReady,
		User:   newUser,
	}
	newID, err := e.Add(expression)
	require.NoError(t, err)

	expression.Status = db.ExpressionReady
	expression.ID = newID
	err = e.UpdateExpression(expression)
	require.NoError(t, err)

	expression, err = e.GetByID(newID)
	require.NoError(t, err)
	assert.Equal(t, "2 + 2", expression.Value)
	assert.InDelta(t, float64(4), expression.Answer, 0.0001)
	assert.Equal(t, "ok", expression.Logs)
	assert.Equal(t, db.ExpressionReady, expression.Status)

	err = e.Delete(newID)
	require.NoError(t, err)
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestIsExpressionWorking(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)
	e := expressionstorage.New(d, 1, &sync.Map{})

	newUser := CreateTestUser(t, d)

	expression := db.Expression{
		ID:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionWorking,
		User:   newUser,
	}
	newID, err := e.Add(expression)
	require.NoError(t, err)

	ok, err := e.IsExpressionWorking(newID)
	require.NoError(t, err)
	assert.True(t, ok)

	err = e.Delete(newID)
	require.NoError(t, err)
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestIsExpressionNotReady(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)
	e := expressionstorage.New(d, 1, &sync.Map{})

	newUser := CreateTestUser(t, d)

	expression := db.Expression{
		ID:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionNotReady,
		User:   newUser,
	}
	newID, err := e.Add(expression)
	require.NoError(t, err)

	ok, err := e.IsExpressionNotReady(newID)
	require.NoError(t, err)
	assert.True(t, ok)

	err = e.Delete(newID)
	require.NoError(t, err)
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestNotWorkingExpressions(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)
	e := expressionstorage.New(d, 1, &sync.Map{})

	newUser := CreateTestUser(t, d)

	expression := db.Expression{
		ID:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionNotReady,
		User:   newUser,
	}
	newID, err := e.Add(expression)
	require.NoError(t, err)

	expressions := e.GetNotWorkingExpressions()
	assert.NotEmpty(t, expressions)

	err = e.Delete(newID)
	require.NoError(t, err)
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestGetByUserAndId(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)
	e := expressionstorage.New(d, 1, &sync.Map{})

	newUser := CreateTestUser(t, d)

	expression := db.Expression{
		ID:     0,
		Value:  "2 + 2",
		Answer: 4,
		Logs:   "ok",
		Status: db.ExpressionNotReady,
		User:   newUser,
	}
	newID, err := e.Add(expression)
	require.NoError(t, err)

	expression, err = e.GetByUserAndID(newUser, newID)
	require.NoError(t, err)
	assert.Equal(t, "2 + 2", expression.Value)
	assert.InDelta(t, float64(4), expression.Answer, 0.0001)
	assert.Equal(t, "ok", expression.Logs)
	assert.Equal(t, db.ExpressionNotReady, expression.Status)

	err = e.Delete(newID)
	require.NoError(t, err)
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestGetByServer(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)
	e := expressionstorage.New(d, 1, &sync.Map{})

	newUser := CreateTestUser(t, d)

	expression := db.Expression{
		ID:         0,
		Value:      "2 + 2",
		Answer:     4,
		Logs:       "ok",
		Status:     db.ExpressionNotReady,
		User:       newUser,
		Servername: "server",
	}
	newID, err := e.Add(expression)
	require.NoError(t, err)

	expressions := e.GetByServer(newUser, "server")
	assert.NotEmpty(t, expressions)

	err = e.Delete(newID)
	require.NoError(t, err)
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestKeepAliveExpressions(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	servers := &sync.Map{}
	servers.Store("server", "")

	e := expressionstorage.New(d, 1, servers)

	newUser := CreateTestUser(t, d)

	expression := db.Expression{
		ID:             0,
		Value:          "2 + 2",
		Answer:         4,
		Logs:           "ok",
		Status:         db.ExpressionWorking,
		User:           newUser,
		Servername:     "server",
		AliveExpiresAt: int(time.Now().Unix()),
	}
	newID, err := e.Add(expression)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	expression, err = e.GetByID(newID)
	require.NoError(t, err)
	assert.Equal(t, db.ExpressionNotReady, expression.Status)

	val, ok := servers.Load("server")
	require.True(t, ok)
	assert.NotNil(t, val)

	err = e.Delete(newID)
	require.NoError(t, err)
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}
