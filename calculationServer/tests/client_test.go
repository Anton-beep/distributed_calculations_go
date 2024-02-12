package tests

import (
	"calculationServer/internal/storage_client"
	"calculationServer/pkg/expression_parser"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func clientSetup(t *testing.T, servUrl string) *storage_client.Client {
	err := os.Setenv("STORAGE_URL", servUrl)
	require.NoError(t, err)
	err = os.Setenv("NUMBER_OF_CALCULATORS", "1")
	require.NoError(t, err)
	client, err := storage_client.New()
	require.NoError(t, err)
	return client
}

func TestGetUpdates(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprint(w, `{"tasks":[{"id":1,"value":"2+2","answer":0,"logs":"","status":0}], "message":"ok"`)
			if err != nil {
				require.NoError(t, err)
			}
		}))
	defer server.Close()

	client := clientSetup(t, server.URL)

	expressions, err := client.GetUpdates()
	require.NoError(t, err)

	assert.Equal(t, []storage_client.Expression{{ID: 1, Value: "2+2", Answer: 0, Logs: "", Status: 0}}, expressions)
}

func TestConfirmTrue(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var ans storage_client.Expression
			err = json.Unmarshal(body, &ans)
			require.NoError(t, err)
			assert.Equal(t, storage_client.Expression{ID: 1, Value: "2+2", Answer: 0, Logs: "", Status: 0}, ans)

			_, err = fmt.Fprint(w, `{"confirm": true, "message": "ok"}`)
			if err != nil {
				require.NoError(t, err)
			}
			w.WriteHeader(http.StatusOK)
		}))
	defer server.Close()

	client := clientSetup(t, server.URL)

	ans := storage_client.Expression{ID: 1, Value: "2+2", Answer: 0, Logs: "", Status: 0}
	confirm, err := client.TryToConfirm(ans)
	require.NoError(t, err)
	assert.Equal(t, true, confirm)
}

func TestConfirmFalse(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := fmt.Fprint(w, `{"confirm": false, "message": "expression is not in pending"}`)
			if err != nil {
				require.NoError(t, err)
			}
			w.WriteHeader(http.StatusOK)
		}))
	defer server.Close()

	client := clientSetup(t, server.URL)

	ans := storage_client.Expression{ID: 1, Value: "2+2", Answer: 0, Logs: "", Status: 0}
	confirm, err := client.TryToConfirm(ans)
	require.NoError(t, err)
	assert.Equal(t, false, confirm)
}

func TestGetOperationsAndTimes(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := fmt.Fprint(w, `{"+":100, "-":200, "*":300, "/":400}`)
			if err != nil {
				require.NoError(t, err)
			}
			w.WriteHeader(http.StatusOK)
		}))
	defer server.Close()

	client := clientSetup(t, server.URL)

	operations, err := client.GetOperationsAndTimes()
	require.NoError(t, err)
	assert.Equal(t, expression_parser.ExecTimeConfig{
		TimeAdd:      100 * time.Millisecond,
		TimeSubtract: 200 * time.Millisecond,
		TimeMultiply: 300 * time.Millisecond,
		TimeDivide:   400 * time.Millisecond,
	}, operations)
}
