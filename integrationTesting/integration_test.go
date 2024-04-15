package integrationTesting

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func startDocker(t *testing.T) {
	err := os.Chdir("../")
	require.NoError(t, err)
	cmd := exec.Command("docker-compose", "up", "-d")
	err = cmd.Run()
	require.NoError(t, err)

	//wait for the server to start
	for {
		resp, err := http.Get("http://localhost:8080/api/v1/ping")
		if err == nil {
			resp.Body.Close()
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func clearDocker(t *testing.T) {
	err := os.Chdir("integrationTesting")
	require.NoError(t, err)
	cmd := exec.Command("docker-compose", "down")
	err = cmd.Run()
	require.NoError(t, err)
}

type Expression struct {
	ID                 int     `json:"id"`
	Value              string  `json:"value"`
	Answer             float64 `json:"answer"`
	Logs               string  `json:"logs"`
	Status             int     `json:"ready"`
	AliveExpiresAt     int     `json:"alive_expires_at"`
	CreationTime       string  `json:"creation_time"`
	EndCalculationTime string  `json:"end_calculation_time"`
	Servername         string  `json:"server_name"`
	User               int     `json:"user_id"`
}

func TestSimpleIntegration(t *testing.T) {
	startDocker(t)

	login := "test"
	password := "test"

	// Register user
	body := fmt.Sprintf(`{"login": "%s", "password": "%s"}`, login, password)

	resp, err := http.Post("http://localhost:8080/api/v1/register", "application/json", strings.NewReader(body))
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	out, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	registerResult := struct {
		Access  string `json:"access"`
		Message string `json:"message"`
	}{}

	err = json.Unmarshal(out, &registerResult)
	require.NoError(t, err)

	// Try to login

	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/login", strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+registerResult.Access)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)

	out, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	loginResult := struct {
		Access  string `json:"access"`
		Message string `json:"message"`
	}{}

	err = json.Unmarshal(out, &loginResult)
	require.NoError(t, err)
	// Create expression

	req, err = http.NewRequest("POST", "http://localhost:8080/api/v1/expression", strings.NewReader(`{"expression": "2+2"}`))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+loginResult.Access)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)

	out, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	postExpressionResult := struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	}{}

	err = json.Unmarshal(out, &postExpressionResult)
	require.NoError(t, err)
	expressionID := postExpressionResult.ID

	// Wait for an answer

	for {
		body = fmt.Sprintf(`{"id": %d}`, expressionID)
		req, err = http.NewRequest("GET", "http://localhost:8080/api/v1/expressionById", strings.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+loginResult.Access)

		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)

		if resp.StatusCode == 200 {
			break
		}

		out, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		resp.Body.Close()

		getExpression := struct {
			Expression `json:"expression"`
			Message    string `json:"message"`
		}{}

		err = json.Unmarshal(out, &getExpression)
		require.NoError(t, err)

		if getExpression.Status == 2 {
			assert.InDelta(t, 4.0, getExpression.Answer, 0.0001)
			break
		}
	}

	// Edit time configuration

	req, err = http.NewRequest("POST", "http://localhost:8080/api/v1/postOperationsAndTimes", strings.NewReader(`{"+": 2000, "-": 2000, "*": 2000, "/": 2000}`))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+loginResult.Access)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)

	out, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	postOperationsAndTimesResult := struct {
		Message string `json:"message"`
	}{}

	err = json.Unmarshal(out, &postOperationsAndTimesResult)
	require.NoError(t, err)

	assert.Contains(t, postOperationsAndTimesResult.Message, "changed for")

	clearDocker(t)
}
