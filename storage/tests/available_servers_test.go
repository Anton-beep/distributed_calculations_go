package tests

import (
	"github.com/stretchr/testify/require"
	"storage/internal/availableservers"
	"storage/internal/db"
	"storage/internal/expressionstorage"
	"sync"
	"testing"
)

func TestGetExpressions(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)
	e := expressionstorage.New(d, 1, &sync.Map{})
	a := availableservers.New(e)

	newUser := CreateTestUser(t, d)

	newID, err := e.Add(db.Expression{
		ID:         0,
		Value:      "2 + 2",
		Answer:     4,
		Logs:       "ok",
		Status:     db.ExpressionReady,
		User:       newUser,
		Servername: "server1",
	})

	require.NoError(t, err)

	a.Add("server1")

	expressions := a.GetExpressions(newUser, "server1")
	require.Len(t, expressions, 1)

	servers := a.GetAll()
	require.Len(t, servers, 1)

	a.Remove("server1")

	err = e.Delete(newID)
	require.NoError(t, err)
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}
