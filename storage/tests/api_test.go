package tests

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"storage/internal/Api"
	"storage/internal/Db"
	"strings"
	"testing"
)

func TestPing(t *testing.T) {
	d, err := Db.New()
	assert.NoError(t, err)

	a := Api.New(d)
	router := a.Start()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}

func TestPostGetExpression(t *testing.T) {
	d, err := Db.New()
	assert.NoError(t, err)

	a := Api.New(d)
	router := a.Start()

	w := httptest.NewRecorder()
	body, _ := json.Marshal(Api.InPostExpression{Expression: "2+2"})
	req, _ := http.NewRequest("POST", "/api/v1/expression", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out1 Api.OutPostExpression
	err = json.Unmarshal(w.Body.Bytes(), &out1)
	assert.NoError(t, err)

	assert.Equal(t, "ok", out1.Message)
	id := out1.Id

	var in Api.InGetExpressionById
	in.Id = id
	body, _ = json.Marshal(in)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/expressionById", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out2 Api.OutGetExpressionById
	err = json.Unmarshal(w.Body.Bytes(), &out2)
	assert.NoError(t, err)

	assert.Equal(t, "ok", out2.Message)

	assert.Equal(t, "2+2", out2.Expression.Value)
}

func TestGetOperationsAndTimes(t *testing.T) {
	d, err := Db.New()
	assert.NoError(t, err)

	a := Api.New(d)
	router := a.Start()

	w := httptest.NewRecorder()
	in := map[string]int{
		"+": 100,
		"-": 200,
		"*": 300,
		"/": 400,
	}
	body, _ := json.Marshal(in)

	req, _ := http.NewRequest("POST", "/api/v1/postOperationsAndTimes", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	out := Api.OutGetOperationsAndTimes{}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/getOperationsAndTimes", nil)
	router.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &out)
	assert.NoError(t, err)

	assert.Equal(t, "ok", out.Message)
	assert.Equal(t, 100, out.Data["+"])
	assert.Equal(t, 200, out.Data["-"])
	assert.Equal(t, 300, out.Data["*"])
	assert.Equal(t, 400, out.Data["/"])

	assert.Equal(t, 200, w.Code)
}

func TestGetUpdates(t *testing.T) {
	d, err := Db.New()
	assert.NoError(t, err)

	a := Api.New(d)
	router := a.Start()

	w := httptest.NewRecorder()
	body, _ := json.Marshal(Api.InPostExpression{Expression: "2+2"})
	req, _ := http.NewRequest("POST", "/api/v1/expression", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/getUpdates", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out Api.OutGetUpdates
	err = json.Unmarshal(w.Body.Bytes(), &out)
	assert.NoError(t, err)
	assert.Equal(t, "ok", out.Message)
	assert.Equal(t, "2+2", out.Expressions[0].Value)
}

func TestConfirmStartCalculating(t *testing.T) {
	d, err := Db.New()
	assert.NoError(t, err)

	a := Api.New(d)
	router := a.Start()

	w := httptest.NewRecorder()
	body, _ := json.Marshal(Api.InPostExpression{Expression: "2+2"})
	req, _ := http.NewRequest("POST", "/api/v1/expression", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out1 Api.OutPostExpression
	err = json.Unmarshal(w.Body.Bytes(), &out1)
	assert.NoError(t, err)

	assert.Equal(t, "ok", out1.Message)
	id := out1.Id

	var in Api.InConfirmStartOfCalculating
	in.Expression.Id = id
	body, _ = json.Marshal(in)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/confirmStartCalculating", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out2 Api.OutConfirmStartOfCalculating
	err = json.Unmarshal(w.Body.Bytes(), &out2)
	assert.NoError(t, err)

	assert.Equal(t, "ok", out2.Message)
	assert.Equal(t, true, out2.Confirm)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/confirmStartCalculating", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var out3 Api.OutConfirmStartOfCalculating
	err = json.Unmarshal(w.Body.Bytes(), &out3)
	assert.NoError(t, err)

	assert.Equal(t, false, out3.Confirm)
}

func TestPostResult(t *testing.T) {
	d, err := Db.New()
	assert.NoError(t, err)

	a := Api.New(d)
	router := a.Start()

	w := httptest.NewRecorder()
	body, _ := json.Marshal(Api.InPostExpression{Expression: "2+2"})
	req, _ := http.NewRequest("POST", "/api/v1/expression", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out Api.OutPostExpression
	err = json.Unmarshal(w.Body.Bytes(), &out)
	assert.NoError(t, err)

	var in Api.InConfirmStartOfCalculating
	in.Expression.Id = out.Id
	body, _ = json.Marshal(in)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/confirmStartCalculating", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out2 Api.OutConfirmStartOfCalculating
	err = json.Unmarshal(w.Body.Bytes(), &out2)
	assert.NoError(t, err)

	assert.Equal(t, "ok", out2.Message)

	var in2 Api.InPostResult
	in2.Expression.Id = out.Id
	in2.Expression.Status = Db.ExpressionReady
	in2.Expression.Answer = 4
	in2.Expression.Logs = "ok"
	body, _ = json.Marshal(in2)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/postResult", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out3 Api.OutPostResult
	err = json.Unmarshal(w.Body.Bytes(), &out3)
	assert.NoError(t, err)

	assert.Equal(t, "ok", out3.Message)
}
