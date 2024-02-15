package tests_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"storage/internal/api"
	"storage/internal/db"
	"strings"
	"testing"
)

func TestPing(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	a := api.New(d)
	router := a.Start()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}

func TestPostGetExpression(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	a := api.New(d)
	router := a.Start()

	w := httptest.NewRecorder()
	body, _ := json.Marshal(api.InPostExpression{Expression: "2+2"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/expression", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out1 api.OutPostExpression
	err = json.Unmarshal(w.Body.Bytes(), &out1)
	require.NoError(t, err)

	assert.Equal(t, "ok", out1.Message)
	id := out1.ID

	var in api.InGetExpressionByID
	in.ID = id
	body, _ = json.Marshal(in)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/expressionById", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out2 api.OutGetExpressionByID
	err = json.Unmarshal(w.Body.Bytes(), &out2)
	require.NoError(t, err)

	assert.Equal(t, "ok", out2.Message)

	assert.Equal(t, "2+2", out2.Expression.Value)
}

func TestGetOperationsAndTimes(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	a := api.New(d)
	router := a.Start()

	w := httptest.NewRecorder()
	in := map[string]int{
		"+": 100,
		"-": 200,
		"*": 300,
		"/": 400,
	}
	body, _ := json.Marshal(in)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/postOperationsAndTimes", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	out := api.OutGetOperationsAndTimes{}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/getOperationsAndTimes", nil)
	router.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &out)
	require.NoError(t, err)

	assert.Equal(t, "ok", out.Message)
	assert.Equal(t, 100, out.Data["+"])
	assert.Equal(t, 200, out.Data["-"])
	assert.Equal(t, 300, out.Data["*"])
	assert.Equal(t, 400, out.Data["/"])

	assert.Equal(t, 200, w.Code)
}

func TestGetUpdates(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	a := api.New(d)
	router := a.Start()

	w := httptest.NewRecorder()
	body, _ := json.Marshal(api.InPostExpression{Expression: "2+2"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/expression", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/getUpdates", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out api.OutGetUpdates
	err = json.Unmarshal(w.Body.Bytes(), &out)
	require.NoError(t, err)
	assert.Equal(t, "ok", out.Message)
	assert.Equal(t, "2+2", out.Expressions[0].Value)
}

func TestConfirmStartCalculating(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	a := api.New(d)
	router := a.Start()

	w := httptest.NewRecorder()
	body, _ := json.Marshal(api.InPostExpression{Expression: "2+2"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/expression", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out1 api.OutPostExpression
	err = json.Unmarshal(w.Body.Bytes(), &out1)
	require.NoError(t, err)

	assert.Equal(t, "ok", out1.Message)
	id := out1.ID

	var in api.InConfirmStartOfCalculating
	in.Expression.ID = id
	body, _ = json.Marshal(in)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/confirmStartCalculating", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out2 api.OutConfirmStartOfCalculating
	err = json.Unmarshal(w.Body.Bytes(), &out2)
	require.NoError(t, err)

	assert.Equal(t, "ok", out2.Message)
	assert.True(t, out2.Confirm)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/confirmStartCalculating", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var out3 api.OutConfirmStartOfCalculating
	err = json.Unmarshal(w.Body.Bytes(), &out3)
	require.NoError(t, err)

	assert.False(t, out3.Confirm)
}

func TestPostResult(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	a := api.New(d)
	router := a.Start()

	w := httptest.NewRecorder()
	body, _ := json.Marshal(api.InPostExpression{Expression: "2+2"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/expression", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out api.OutPostExpression
	err = json.Unmarshal(w.Body.Bytes(), &out)
	require.NoError(t, err)

	var in api.InConfirmStartOfCalculating
	in.Expression.ID = out.ID
	body, _ = json.Marshal(in)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/confirmStartCalculating", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out2 api.OutConfirmStartOfCalculating
	err = json.Unmarshal(w.Body.Bytes(), &out2)
	require.NoError(t, err)

	assert.Equal(t, "ok", out2.Message)

	var in2 api.InPostResult
	in2.Expression.ID = out.ID
	in2.Expression.Status = db.ExpressionReady
	in2.Expression.Answer = 4
	in2.Expression.Logs = "ok"
	body, _ = json.Marshal(in2)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/postResult", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out3 api.OutPostResult
	err = json.Unmarshal(w.Body.Bytes(), &out3)
	require.NoError(t, err)

	assert.Equal(t, "ok", out3.Message)
}

func TestAlive(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	a := api.New(d)
	router := a.Start()

	w := httptest.NewRecorder()
	var in api.InPostExpression
	in.Expression = "2+2"
	body, _ := json.Marshal(in)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/expression", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	var in2 api.InConfirmStartOfCalculating
	in2.Expression.ID = 1
	body, _ = json.Marshal(in2)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/confirmStartCalculating", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	var in3 api.InKeepAlive
	in3.Expression.ID = 1
	body, _ = json.Marshal(in3)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/keepAlive", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	var in4 api.InGetExpressionByID
	in4.ID = 1
	body, _ = json.Marshal(in4)
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/expressionById", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out api.OutGetExpressionByID
	err = json.Unmarshal(w.Body.Bytes(), &out)
	require.NoError(t, err)

	assert.Equal(t, "ok", out.Message)
	assert.Equal(t, 1, out.Expression.ID)
	assert.Greater(t, out.Expression.AliveExpiresAt, 0)
}

func TestGetServers(t *testing.T) {
	d, err := db.New()
	require.NoError(t, err)

	a := api.New(d)
	router := a.Start()

	w := httptest.NewRecorder()
	var in api.InPostExpression
	in.Expression = "2+2"
	body, _ := json.Marshal(in)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/expression", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	var in2 api.InConfirmStartOfCalculating
	in2.Expression.ID = 1
	in2.Expression.Servername = "server1"
	body, _ = json.Marshal(in2)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/confirmStartCalculating", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var in3 api.InGetExpressionByServer
	in3.ServerName = "server1"
	body, _ = json.Marshal(in3)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/getExpressionsByServer", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var out2 api.OutGetExpressionByServer
	err = json.Unmarshal(w.Body.Bytes(), &out2)
	require.NoError(t, err)

	assert.Equal(t, "ok", out2.Message)
	assert.Equal(t, "server1", out2.Expressions[0].Servername)
}
