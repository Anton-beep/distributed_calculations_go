package expressionstorage

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"storage/internal/db"
	"sync"
	"time"
)

type ExpressionStorage struct {
	expressions  sync.Map
	db           *db.APIDb
	checkAlive   time.Duration
	serverStatus *sync.Map
}

func New(indb *db.APIDb, checkAlive time.Duration, serverStatus *sync.Map) *ExpressionStorage {
	e := &ExpressionStorage{
		db: indb,
	}

	// check saved data in database and uploads it to memory
	expressions, err := indb.GetAllExpressions()
	if err != nil {
		zap.S().Error(err)
	}
	for _, expression := range expressions {
		e.expressions.Store(expression.ID, expression)
	}

	e.checkAlive = checkAlive
	e.serverStatus = serverStatus
	go e.keepAliveExpressions()
	return e
}

func (e *ExpressionStorage) Add(expression db.Expression) (int, error) {
	// sync with database
	newID, err := e.db.AddExpression(expression)
	if err != nil {
		return 0, err
	}
	expression.ID = newID

	e.expressions.Store(newID, expression)
	return newID, nil
}

func (e *ExpressionStorage) GetAll(userID int) []db.Expression {
	expressions := make([]db.Expression, 0)
	e.expressions.Range(func(_, value interface{}) bool {
		if value.(db.Expression).User == userID {
			expressions = append(expressions, value.(db.Expression))
		}
		return true
	})
	return expressions
}

func (e *ExpressionStorage) GetByID(id int) (db.Expression, error) {
	if expression, ok := e.expressions.Load(id); ok {
		return expression.(db.Expression), nil
	}
	return db.Expression{}, errors.New("expression is not found")
}

func (e *ExpressionStorage) GetByUserAndID(userID int, id int) (db.Expression, error) {
	expression, err := e.GetByID(id)
	if err != nil {
		return db.Expression{}, err
	}
	if expression.User != userID {
		return db.Expression{}, errors.New("expression is not found")
	}
	return expression, nil
}

// GetNotWorkingExpressions returns all expressions that have Status == ExpressionNotReady.
func (e *ExpressionStorage) GetNotWorkingExpressions() []db.Expression {
	expressions := make([]db.Expression, 0)
	e.expressions.Range(func(_, value interface{}) bool {
		if value.(db.Expression).Status == db.ExpressionNotReady {
			expressions = append(expressions, value.(db.Expression))
		}
		return true
	})
	return expressions
}

// UpdateExpression updates expression in pendingExpressions and sync with database.
func (e *ExpressionStorage) UpdateExpression(expression db.Expression) error {
	if _, ok := e.expressions.Load(expression.ID); !ok {
		return errors.New("expression is not found")
	}
	e.expressions.Store(expression.ID, expression)
	// sync with database
	if err := e.db.UpdateExpression(expression); err != nil {
		return err
	}
	return nil
}

// IsExpressionWorking returns true if expression is in pendingExpressions and has Status == ExpressionWorking.
func (e *ExpressionStorage) IsExpressionWorking(id int) (bool, error) {
	if expression, ok := e.expressions.Load(id); ok {
		return expression.(db.Expression).Status == db.ExpressionWorking, nil
	}
	return false, errors.New("expression is not found")
}

// IsExpressionNotReady returns true if expression is in pendingExpressions and has Status == ExpressionNotReady.
func (e *ExpressionStorage) IsExpressionNotReady(id int) (bool, error) {
	if expression, ok := e.expressions.Load(id); ok {
		return expression.(db.Expression).Status == db.ExpressionNotReady, nil
	}
	return false, errors.New("expression is not found")
}

func (e *ExpressionStorage) Delete(id int) error {
	if _, ok := e.expressions.Load(id); ok {
		e.expressions.Delete(id)
	}
	// sync with database
	if err := e.db.DeleteExpression(id); err != nil {
		return err
	}
	return nil
}

// keepAliveExpressions checks all expressions and if aliveExpiresAt is less than now, then change to not ready,
// so it will be calculated again via getUpdates.
func (e *ExpressionStorage) keepAliveExpressions() {
	// check all expressions and if aliveExpiresAt is less than now, then change to not ready
	for {
		time.Sleep(e.checkAlive)
		e.expressions.Range(func(key, value interface{}) bool {
			expression, ok := value.(db.Expression)
			if !ok {
				zap.S().Error("expression is not found")
			}

			if expression.Status == db.ExpressionWorking && expression.AliveExpiresAt < int(time.Now().Unix()) {
				// change to not ready, so it will be calculated again
				zap.S().Info(fmt.Sprintf("expression ID %v is not alive, change to not ready."+
					" Dead server: %v", expression.ID, expression.Servername))

				e.serverStatus.Range(func(key, _ interface{}) bool {
					zap.S().Info(fmt.Sprintf("%v", key))
					if expression.Servername == key.(string) {
						e.serverStatus.Store(key, fmt.Sprintf("%v -> server %v is not alive",
							time.Now().Format("01-02-2006 15:04:05"), expression.Servername))
					}
					return true
				})
				expression.Status = db.ExpressionNotReady
				e.expressions.Store(key, expression)
				// sync with database
				if err := e.db.UpdateExpression(expression); err != nil {
					zap.S().Error(err)
				}
			}
			return true
		})
	}
}

// GetByServer returns all expressions that have ServerName == server.
func (e *ExpressionStorage) GetByServer(userID int, server string) []db.Expression {
	expressions := make([]db.Expression, 0)
	e.expressions.Range(func(_, value interface{}) bool {
		if value.(db.Expression).Servername == server && value.(db.Expression).User == userID {
			expressions = append(expressions, value.(db.Expression))
		}
		return true
	})
	return expressions
}
