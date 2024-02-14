package tests_test

import (
	"calculationServer/internal/storageclient"
	"calculationServer/pkg/expressionparser"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func clientSetup(t *testing.T, servURL string) *storageclient.Client {
	t.Setenv("STORAGE_URL", servURL)
	t.Setenv("NUMBER_OF_CALCULATORS", "1")
	t.Setenv("SEND_ALIVE_DURATION", "1")
	client, err := storageclient.New()
	require.NoError(t, err)
	return client
}

func TestGetUpdates(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprint(w, `{"tasks":[{"id":1,"value":"2+2","answer":0,"logs":"","ready":0}], "message":"ok"}`)
			if err != nil {
				require.NoError(t, err)
			}
		}))
	defer server.Close()

	client := clientSetup(t, server.URL)

	expressions, err := client.GetUpdates()
	require.NoError(t, err)

	assert.Equal(t, []storageclient.Expression{{ID: 1, Value: "2+2", Answer: 0, Logs: "", Status: 0}}, expressions)
}

func TestConfirmTrue(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var ans storageclient.SendConfirmStartOfCalculating
			err = json.Unmarshal(body, &ans)
			require.NoError(t, err)
			assert.Equal(t, storageclient.SendConfirmStartOfCalculating{Expression: storageclient.Expression{
				ID:     1,
				Value:  "2+2",
				Answer: 0,
				Logs:   "",
				Status: 0,
			}}, ans)

			_, err = fmt.Fprint(w, `{"confirm": true, "message": "ok"}`)
			if err != nil {
				require.NoError(t, err)
			}
			w.WriteHeader(http.StatusOK)
		}))
	defer server.Close()

	client := clientSetup(t, server.URL)

	ans := storageclient.Expression{ID: 1, Value: "2+2", Answer: 0, Logs: "", Status: 0}
	confirm, err := client.TryToConfirm(ans)
	require.NoError(t, err)
	assert.True(t, confirm)
}

func TestConfirmFalse(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := fmt.Fprint(w, `{"confirm": false, "message": "expression is not in pending"}`)
			if err != nil {
				require.NoError(t, err)
			}
			w.WriteHeader(http.StatusOK)
		}))
	defer server.Close()

	client := clientSetup(t, server.URL)

	ans := storageclient.Expression{ID: 1, Value: "2+2", Answer: 0, Logs: "", Status: 0}
	confirm, err := client.TryToConfirm(ans)
	require.NoError(t, err)
	assert.False(t, confirm)
}

func TestGetOperationsAndTimes(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := fmt.Fprint(w, `{"data": {"+":100, "-":200, "*":300, "/":400}, "message": "ok"}`)
			if err != nil {
				require.NoError(t, err)
			}
			w.WriteHeader(http.StatusOK)
		}))
	defer server.Close()

	client := clientSetup(t, server.URL)

	operations, err := client.GetOperationsAndTimes()
	require.NoError(t, err)
	assert.Equal(t, expressionparser.ExecTimeConfig{
		TimeAdd:      100 * time.Millisecond,
		TimeSubtract: 200 * time.Millisecond,
		TimeMultiply: 300 * time.Millisecond,
		TimeDivide:   400 * time.Millisecond,
	}, operations)
}
