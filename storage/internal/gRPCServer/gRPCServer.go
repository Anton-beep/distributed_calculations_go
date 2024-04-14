package gRPCServer

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"storage/internal/api"
	"storage/internal/availableservers"
	"storage/internal/db"
	"storage/internal/expressionstorage"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	expressions    *expressionstorage.ExpressionStorage
	servers        *availableservers.AvailableServers
	checkAlive     int
	execTimeConfig *api.ExecTimeConfig
	ExpressionsServiceServer
	statusWorkers *sync.Map
	db            *db.APIDb
}

func New(expressions *expressionstorage.ExpressionStorage, servers *availableservers.AvailableServers, execTimeConfig *api.ExecTimeConfig, statusWorkers *sync.Map, db *db.APIDb) *Server {
	num, err := strconv.Atoi(os.Getenv("CHECK_SERVER_DURATION"))
	if err != nil {
		zap.S().Fatal(err)
	}
	return &Server{
		expressions:    expressions,
		servers:        servers,
		checkAlive:     num,
		execTimeConfig: execTimeConfig,
		statusWorkers:  statusWorkers,
		db:             db,
	}
}

func dbExpressionTogRPCExpression(expression db.Expression) *Expression {
	return &Expression{
		Id:                 int64(expression.ID),
		Value:              expression.Value,
		Answer:             expression.Answer,
		Logs:               expression.Logs,
		Status:             int32(expression.Status),
		AliveExpiresAt:     int64(expression.AliveExpiresAt),
		CreationTime:       expression.CreationTime,
		EndCalculationTime: expression.EndCalculationTime,
		ServerName:         expression.Servername,
		UserId:             int64(expression.User),
	}
}

func gRPCExpressionTodbExpression(expression *Expression) db.Expression {
	return db.Expression{
		ID:                 int(expression.Id),
		Value:              expression.Value,
		Answer:             expression.Answer,
		Logs:               expression.Logs,
		Status:             int(expression.Status),
		AliveExpiresAt:     int(expression.AliveExpiresAt),
		CreationTime:       expression.CreationTime,
		EndCalculationTime: expression.EndCalculationTime,
		Servername:         expression.ServerName,
		User:               int(expression.UserId),
	}
}

func (s *Server) GetUpdates(_ *Empty, stream ExpressionsService_GetUpdatesServer) error {
	expressions := s.expressions.GetNotWorkingExpressions()
	for _, expression := range expressions {
		err := stream.Send(dbExpressionTogRPCExpression(expression))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) ConfirmStartCalculating(_ context.Context, e *Expression) (*Confirm, error) {
	expression := gRPCExpressionTodbExpression(e)
	ok, err := s.expressions.IsExpressionNotReady(expression.ID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("expression is not in pending")
	}

	// change to working
	expression.Status = db.ExpressionWorking
	expression.AliveExpiresAt = int(time.Now().Add(time.Duration(s.checkAlive) * time.Second).Unix())
	if err = s.expressions.UpdateExpression(expression); err != nil {
		return &Confirm{
			Confirm: false,
		}, err
	}

	// add server
	s.servers.Add(expression.Servername)
	return &Confirm{
		Confirm: true,
	}, nil
}

func (s *Server) PostResult(_ context.Context, e *Expression) (*Message, error) {
	// check if expression is in working
	expression := gRPCExpressionTodbExpression(e)
	ok, err := s.expressions.IsExpressionWorking(expression.ID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &Message{
			Message: "expression is not in working",
		}, err
	}

	expression.EndCalculationTime = time.Now().Format("2006-01-02 15:04:05")
	if err = s.expressions.UpdateExpression(expression); err != nil {
		return nil, err
	}

	// add server
	s.servers.Add(expression.Servername)
	s.statusWorkers.Store(expression.Servername, fmt.Sprintf("%v -> server %v finished calculating %v",
		time.Now().Format("01-02-2006 15:04:05"), expression.Servername, expression.Value))
	return &Message{
		Message: "ok",
	}, nil
}

func (s *Server) KeepAlive(_ context.Context, msg *KeepAliveMsg) (*Empty, error) {
	expression := gRPCExpressionTodbExpression(msg.Expression)
	expression, err := s.expressions.GetByID(expression.ID)
	if err != nil {
		return nil, err
	}

	expression.AliveExpiresAt = int(time.Now().Add(time.Duration(s.checkAlive) * time.Second).Unix())
	err = s.expressions.UpdateExpression(expression)
	if err != nil {
		return nil, err
	}

	s.statusWorkers.Store(expression.Servername, msg.StatusWorkers)
	return nil, nil
}

func (s *Server) GetOperationsAndTimes(_ context.Context, e *Expression) (*OperationsAndTimes, error) {
	operations, err := s.db.GetUserOperations(int(e.UserId))
	if err != nil {
		return nil, err
	}
	return &OperationsAndTimes{
		TimeAdd:      int64(operations.TimeAdd),
		TimeSubtract: int64(operations.TimeSubtract),
		TimeMultiply: int64(operations.TimeMultiply),
		TimeDivide:   int64(operations.TimeDivide),
		Message:      "ok",
	}, nil
}
