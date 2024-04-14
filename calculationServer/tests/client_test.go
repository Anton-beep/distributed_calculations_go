package tests

import (
	"calculationServer/internal/storageclient"
	"calculationServer/pkg/expressionparser"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetUpdates(t *testing.T) {
	client, s, conn := ClientAndServerSetup(t)
	defer s.Stop()
	defer conn.Close()

	GetUpdatesValues = []*storageclient.Expression{
		{
			Id:     0,
			Value:  "1+1",
			UserId: 0,
		},
		{
			Id:     1,
			Value:  "2+2",
			UserId: 1,
		},
	}

	resp, err := client.GetUpdates()
	assert.NoError(t, err)

	for _, exp := range resp {
		if exp.Id == 0 {
			assert.Equal(t, "1+1", exp.Value)
		} else {
			assert.Equal(t, "2+2", exp.Value)
		}
	}
}

func TestConfirm(t *testing.T) {
	client, s, conn := ClientAndServerSetup(t)
	defer s.Stop()
	defer conn.Close()

	ConfirmValue = &storageclient.Confirm{Confirm: true}

	resp, err := client.TryToConfirm(&storageclient.Expression{})
	if err != nil {
		assert.NoError(t, err)
	}

	assert.True(t, resp)

	ConfirmValue = &storageclient.Confirm{Confirm: false}

	resp, err = client.TryToConfirm(&storageclient.Expression{})
	assert.NoError(t, err)

	assert.False(t, resp)
}

func TestGetOperationsAndTimes(t *testing.T) {
	client, s, conn := ClientAndServerSetup(t)
	defer s.Stop()
	defer conn.Close()

	OperationsAndTimesValue = &storageclient.OperationsAndTimes{
		TimeAdd:      1000,
		TimeSubtract: 1000,
		TimeDivide:   1000,
		TimeMultiply: 1000,
	}

	resp, err := client.GetOperationsAndTimes(&storageclient.Expression{})
	assert.NoError(t, err)

	assert.Equal(t, expressionparser.ExecTimeConfig{
		TimeAdd:      time.Duration(1000) * time.Millisecond,
		TimeSubtract: time.Duration(1000) * time.Millisecond,
		TimeDivide:   time.Duration(1000) * time.Millisecond,
		TimeMultiply: time.Duration(1000) * time.Millisecond,
	}, resp)
}

func TestFullRun(t *testing.T) {
	client, s, conn := ClientAndServerSetup(t)
	defer s.Stop()
	defer conn.Close()

	GetUpdatesValues = []*storageclient.Expression{
		{
			Id:     0,
			Value:  "1+1",
			UserId: 0,
		},
	}

	ConfirmValue = &storageclient.Confirm{Confirm: true}

	OperationsAndTimesValue = &storageclient.OperationsAndTimes{
		TimeAdd:      1000,
		TimeSubtract: 1000,
		TimeDivide:   1000,
		TimeMultiply: 1000,
	}

	PostResultValue = &storageclient.Message{Message: "ok"}
	PostResultChannel = make(chan *storageclient.Expression, 1)

	go client.Run()
	resExp := <-PostResultChannel

	assert.Equal(t, "1+1", resExp.Value)
	assert.Equal(t, float64(2), resExp.Answer)
	assert.Equal(t, int64(0), resExp.UserId)
	assert.Equal(t, int64(0), resExp.Id)
}
