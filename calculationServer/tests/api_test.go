package tests

import (
	"calculationServer/internal/Api"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	a := Api.NewApi(STORAGE_URL, SECRET_KEY)
	router := a.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", BASE_PATH+"/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	pong := &Api.Pong{Message: "pong"}
	b, err := json.Marshal(pong)

	assert.NoError(t, err)
	assert.Equal(t, string(b), w.Body.String())
}
