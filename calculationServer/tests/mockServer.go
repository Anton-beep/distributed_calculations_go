package tests

import (
	"calculationServer/internal/storageclient"
	"context"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

type mockServer struct {
	storageclient.ExpressionsServiceServer
}

func initServer() *grpc.Server {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	storageclient.RegisterExpressionsServiceServer(s, &mockServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	return s
}

func ClientAndServerSetup(t *testing.T) (*storageclient.Client, *grpc.Server, *grpc.ClientConn) {
	t.Setenv("NUMBER_OF_CALCULATORS", "1")
	t.Setenv("SEND_ALIVE_DURATION", "1")
	client, err := storageclient.New()
	require.NoError(t, err)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	client.SetConnection(conn)

	server := initServer()
	return client, server, conn
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

var GetUpdatesValues []*storageclient.Expression

func (m *mockServer) GetUpdates(_ *storageclient.Empty, stream storageclient.ExpressionsService_GetUpdatesServer) error {
	for _, value := range GetUpdatesValues {
		err := stream.Send(value)
		if err != nil {
			return err
		}
	}
	return nil
}

var ConfirmValue *storageclient.Confirm

func (m *mockServer) ConfirmStartCalculating(_ context.Context, _ *storageclient.Expression) (*storageclient.Confirm, error) {
	return ConfirmValue, nil
}

var OperationsAndTimesValue *storageclient.OperationsAndTimes

func (m *mockServer) GetOperationsAndTimes(_ context.Context, _ *storageclient.Expression) (*storageclient.OperationsAndTimes, error) {
	return OperationsAndTimesValue, nil
}

func (m *mockServer) KeepAlive(_ context.Context, _ *storageclient.KeepAliveMsg) (*storageclient.Empty, error) {
	return &storageclient.Empty{}, nil
}

var PostResultChannel chan *storageclient.Expression
var PostResultValue *storageclient.Message

func (m *mockServer) PostResult(_ context.Context, exp *storageclient.Expression) (*storageclient.Message, error) {
	PostResultChannel <- exp
	return PostResultValue, nil
}
