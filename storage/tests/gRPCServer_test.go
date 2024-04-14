package tests

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"io"
	"net"
	"storage/internal/api"
	"storage/internal/availableservers"
	"storage/internal/db"
	"storage/internal/expressionstorage"
	"storage/internal/gRPCServer"
	"sync"
	"testing"
	"time"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func setupgRPCServer(t *testing.T) (*grpc.Server, *expressionstorage.ExpressionStorage, *db.APIDb) {
	db, err := db.New()
	require.NoError(t, err)

	servers := &sync.Map{}

	expressions := expressionstorage.New(db, 1, servers)

	a := availableservers.New(expressions)

	timeConfig := &api.ExecTimeConfig{
		TimeAdd:      1,
		TimeSubtract: 1,
		TimeDivide:   1,
		TimeMultiply: 1,
	}

	lis = bufconn.Listen(bufSize)
	server := grpc.NewServer()

	gRPCServer.RegisterExpressionsServiceServer(server, gRPCServer.New(expressions, a, timeConfig, servers, db))

	go func() {
		err = server.Serve(lis)
		require.NoError(t, err)
	}()

	return server, expressions, db
}

func setupgRPCClient(t *testing.T) (gRPCServer.ExpressionsServiceClient, *grpc.ClientConn) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	require.NoError(t, err)

	return gRPCServer.NewExpressionsServiceClient(conn), conn
}

var RegisteredCountergRPC = 0

func createNewUser(t *testing.T, d *db.APIDb) int {
	newUser, err := d.AddUser(db.User{
		Login:    fmt.Sprintf("%v%vgrpc", time.Now().Unix(), RegisteredCountergRPC),
		Password: "test",
	})
	RegisteredCountergRPC++
	require.NoError(t, err)

	return newUser
}

func TestGetUpdates(t *testing.T) {
	server, expressions, d := setupgRPCServer(t)
	defer server.Stop()

	client, conn := setupgRPCClient(t)
	defer conn.Close()

	newUser := createNewUser(t, d)

	newExp, err := expressions.Add(db.Expression{
		Value:  "1+123",
		Answer: 2,
		User:   newUser,
	})
	require.NoError(t, err)

	res, err := client.GetUpdates(context.Background(), &gRPCServer.Empty{})
	require.NoError(t, err)

	var found bool
	for {
		exp, err := res.Recv()
		if err == io.EOF {
			break
		}
		if exp.Value == "1+123" {
			found = true
		}
		require.NoError(t, err)
	}

	assert.Equal(t, true, found)
	err = d.DeleteExpression(newExp)
	require.NoError(t, err)
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestConfirmStartCalculating(t *testing.T) {
	server, expressions, d := setupgRPCServer(t)
	defer server.Stop()

	client, conn := setupgRPCClient(t)
	defer conn.Close()

	newUser := createNewUser(t, d)

	newExp, err := expressions.Add(db.Expression{
		Value:  "1+123",
		Answer: 2,
		User:   newUser,
	})
	require.NoError(t, err)

	res, err := client.ConfirmStartCalculating(context.Background(), &gRPCServer.Expression{
		Id:     int64(newExp),
		UserId: int64(newUser),
	})
	require.NoError(t, err)

	assert.Equal(t, true, res.Confirm)

	err = d.DeleteExpression(newExp)
	require.NoError(t, err)
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestPostResult(t *testing.T) {
	server, expressions, d := setupgRPCServer(t)
	defer server.Stop()

	client, conn := setupgRPCClient(t)
	defer conn.Close()

	newUser := createNewUser(t, d)

	newExp, err := expressions.Add(db.Expression{
		Value: "1+123",
		User:  newUser,
	})
	require.NoError(t, err)

	res, err := client.PostResult(context.Background(), &gRPCServer.Expression{
		Id:     int64(newExp),
		UserId: int64(newUser),
		Answer: 2,
	})
	require.NoError(t, err)
	assert.NotEqual(t, "ok", res.Message)

	_, err = client.ConfirmStartCalculating(context.Background(), &gRPCServer.Expression{
		Id:     int64(newExp),
		UserId: int64(newUser),
	})
	require.NoError(t, err)

	res, err = client.PostResult(context.Background(), &gRPCServer.Expression{
		Id:     int64(newExp),
		UserId: int64(newUser),
		Answer: 2,
	})
	require.NoError(t, err)

	assert.Equal(t, "ok", res.Message)

	err = d.DeleteExpression(newExp)
	require.NoError(t, err)
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestKeepAlive(t *testing.T) {
	server, expressions, d := setupgRPCServer(t)
	defer server.Stop()

	client, conn := setupgRPCClient(t)
	defer conn.Close()

	newUser := createNewUser(t, d)

	newExp, err := expressions.Add(db.Expression{
		Value: "1+123",
		User:  newUser,
	})
	require.NoError(t, err)

	_, err = client.ConfirmStartCalculating(context.Background(), &gRPCServer.Expression{
		Id:     int64(newExp),
		UserId: int64(newUser),
	})
	require.NoError(t, err)

	_, err = client.KeepAlive(context.Background(), &gRPCServer.KeepAliveMsg{
		Expression: &gRPCServer.Expression{
			Id:     int64(newExp),
			UserId: int64(newUser),
		},
		StatusWorkers: "ok",
	})
	require.NoError(t, err)

	err = d.DeleteExpression(newExp)
	require.NoError(t, err)
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}

func TestGetOperationsAndTimesForClient(t *testing.T) {
	server, expressions, d := setupgRPCServer(t)
	defer server.Stop()

	client, conn := setupgRPCClient(t)
	defer conn.Close()

	newUser := createNewUser(t, d)

	operId, err := d.AddOperation(db.Operation{
		User:         newUser,
		TimeAdd:      1,
		TimeSubtract: 1,
		TimeDivide:   1,
		TimeMultiply: 1,
	})
	require.NoError(t, err)

	newExp, err := expressions.Add(db.Expression{
		Value: "1+123",
		User:  newUser,
	})
	require.NoError(t, err)

	_, err = client.ConfirmStartCalculating(context.Background(), &gRPCServer.Expression{
		Id:     int64(newExp),
		UserId: int64(newUser),
	})
	require.NoError(t, err)

	res, err := client.GetOperationsAndTimes(context.Background(), &gRPCServer.Expression{
		Id:     int64(newExp),
		UserId: int64(newUser),
	})
	require.NoError(t, err)

	assert.Equal(t, "ok", res.Message)

	err = d.DeleteOperation(operId)
	require.NoError(t, err)
	err = d.DeleteExpression(newExp)
	require.NoError(t, err)
	err = d.DeleteUser(newUser)
	require.NoError(t, err)
}
