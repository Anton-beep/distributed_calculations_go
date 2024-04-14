package tests

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"storage/internal/api"
	"storage/internal/availableservers"
	"storage/internal/db"
	"storage/internal/expressionstorage"
	"strings"
	"sync"
	"testing"
	"time"
)

func CreateApi(t *testing.T) (*db.APIDb, *api.API) {
	d, err := db.New()
	require.NoError(t, err)

	statusWorkers := &sync.Map{}

	e := expressionstorage.New(d, 1, statusWorkers)

	servers := availableservers.New(e)

	timeConfig := &api.ExecTimeConfig{
		TimeAdd:      1,
		TimeSubtract: 1,
		TimeDivide:   1,
		TimeMultiply: 1,
	}

	return d, api.New(d, e, statusWorkers, servers, timeConfig)
}

var RegisteredCounter = 0

func CreateRegisteredUser(t *testing.T, r *gin.Engine) string {
	userData := api.InRegister{
		Login:    fmt.Sprintf("%v%v", time.Now().Unix(), RegisteredCounter),
		Password: "test",
	}
	RegisteredCounter++

	body, _ := json.Marshal(userData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/register", strings.NewReader(string(body)))
	r.ServeHTTP(w, req)

	var out api.OutRegister
	err := json.Unmarshal(w.Body.Bytes(), &out)
	require.NoError(t, err)

	assert.Equal(t, "ok", out.Message)
	return out.Access
}

func TestPing(t *testing.T) {
	_, a := CreateApi(t)
	router := a.Start()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}

func TestRegisterLoginUpdate(t *testing.T) {
	d, a := CreateApi(t)
	router := a.Start()

	userData := api.InRegister{
		Login:    fmt.Sprintf("%vnosuchuser", time.Now().Unix()),
		Password: "test",
	}

	body, _ := json.Marshal(userData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/register", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	var out api.OutRegister
	err := json.Unmarshal(w.Body.Bytes(), &out)
	require.NoError(t, err)

	assert.Equal(t, "ok", out.Message)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/register", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	var out2 api.OutRegister
	err = json.Unmarshal(w.Body.Bytes(), &out2)
	require.NoError(t, err)

	assert.NotEqual(t, "ok", out2.Message)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	var out3 api.OutRegister
	err = json.Unmarshal(w.Body.Bytes(), &out3)
	require.NoError(t, err)

	assert.Equal(t, "ok", out3.Message)

	token := out3.Access

	// wrong login
	userData = api.InRegister{
		Login:    fmt.Sprintf("%vnosuchuser", time.Now().Unix()),
		Password: "testWrong",
	}

	body, _ = json.Marshal(userData)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	var out4 api.OutRegister
	err = json.Unmarshal(w.Body.Bytes(), &out4)
	require.NoError(t, err)

	assert.NotEqual(t, "ok", out4.Message)

	// update
	userUpdate := api.InUpdateUser{
		Login:       fmt.Sprintf("%vnewlogin", time.Now().Unix()),
		OldPassword: "test",
		NewPassword: "testnew",
	}

	body, _ = json.Marshal(userUpdate)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/updateUser", strings.NewReader(string(body)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, req)

	var out5 api.OutUpdateUser
	err = json.Unmarshal(w.Body.Bytes(), &out5)
	require.NoError(t, err)

	assert.Equal(t, "ok", out5.Message)

	token = out5.Access

	// get user
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/getUser", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, req)

	var out6 api.OutGetUser
	err = json.Unmarshal(w.Body.Bytes(), &out6)
	require.NoError(t, err)

	assert.Equal(t, userUpdate.Login, out6.Login)

	// delete user
	user, err := d.GetUserByUsername(userUpdate.Login)
	require.NoError(t, err)

	err = d.DeleteByUserId(user.ID)
	err = d.DeleteUser(user.ID)
	require.NoError(t, err)
}

func TestPostGetExpression(t *testing.T) {
	_, a := CreateApi(t)
	router := a.Start()

	token := CreateRegisteredUser(t, router)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(api.InPostExpression{
		Expression: "2+2",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/expression", strings.NewReader(string(body)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out1 api.OutPostExpression
	err := json.Unmarshal(w.Body.Bytes(), &out1)
	require.NoError(t, err)

	assert.Equal(t, "ok", out1.Message)
	id := out1.ID

	var in api.InGetExpressionByID
	in.ID = id
	body, _ = json.Marshal(in)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/expressionById", strings.NewReader(string(body)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out2 api.OutGetExpressionByID
	err = json.Unmarshal(w.Body.Bytes(), &out2)
	require.NoError(t, err)

	assert.Equal(t, "ok", out2.Message)

	assert.Equal(t, "2+2", out2.Expression.Value)
}

func TestGetOperationsAndTimes(t *testing.T) {
	_, a := CreateApi(t)
	router := a.Start()

	token := CreateRegisteredUser(t, router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/getOperationsAndTimes", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out api.OutGetOperationsAndTimes
	err := json.Unmarshal(w.Body.Bytes(), &out)
	require.NoError(t, err)

	assert.Equal(t, "ok", out.Message)
}

func TestPostOperationsAndTimes(t *testing.T) {
	_, a := CreateApi(t)
	router := a.Start()

	token := CreateRegisteredUser(t, router)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(map[string]int{
		"+": 1,
		"-": 1,
		"/": 1,
		"*": 1,
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/postOperationsAndTimes", strings.NewReader(string(body)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out api.OutPostOperationsAndTimes
	err := json.Unmarshal(w.Body.Bytes(), &out)
	require.NoError(t, err)

	assert.NotEqual(t, "ok", out.Message)
}

func TestGetExpressionsByServer(t *testing.T) {
	_, a := CreateApi(t)
	router := a.Start()

	token := CreateRegisteredUser(t, router)

	var in api.InGetExpressionByServer
	in.ServerName = "server"

	body, _ := json.Marshal(in)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/getExpressionsByServer", strings.NewReader(string(body)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out api.OutGetExpressionByServer
	err := json.Unmarshal(w.Body.Bytes(), &out)
	require.NoError(t, err)

	assert.Equal(t, "ok", out.Message)
}

func TestGetComputingPowers(t *testing.T) {
	_, a := CreateApi(t)
	router := a.Start()

	token := CreateRegisteredUser(t, router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/getComputingPowers", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out api.OutGetComputingPowers
	err := json.Unmarshal(w.Body.Bytes(), &out)
	require.NoError(t, err)

	assert.Equal(t, "ok", out.Message)
}
