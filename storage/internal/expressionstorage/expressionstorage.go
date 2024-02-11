package expressionstorage

import (
	"errors"
	"go.uber.org/zap"
	"storage/internal/db"
	"sync"
)

type ExpressionStorage struct {
	readyExpressions   sync.Map // calculated and ready to return expressions
	pendingExpressions sync.Map // expressions that are not ready yet and need to be calculated
	db                 *db.APIDb
}

func New(indb *db.APIDb) *ExpressionStorage {
	e := &ExpressionStorage{
		db: indb,
	}

	// check saved data in database and uploads it to memory
	expressions, err := indb.GetAllExpressions()
	if err != nil {
		zap.S().Error(err)
	}
	for _, expression := range expressions {
		if expression.Status == db.ExpressionReady {
			e.readyExpressions.Store(expression.ID, expression)
		} else {
			e.pendingExpressions.Store(expression.ID, expression)
		}
	}

	return e
}

func (e *ExpressionStorage) Add(expression db.Expression) (int, error) {
	// sync with database
	newID, err := e.db.AddExpression(expression)
	if err != nil {
		return 0, err
	}
	expression.ID = newID

	e.pendingExpressions.Store(newID, expression)
	return newID, nil
}

func (e *ExpressionStorage) GetAll() []db.Expression {
	expressions := make([]db.Expression, 0)
	e.readyExpressions.Range(func(_, value interface{}) bool {
		expressions = append(expressions, value.(db.Expression))
		return true
	})
	e.pendingExpressions.Range(func(_, value interface{}) bool {
		expressions = append(expressions, value.(db.Expression))
		return true
	})
	return expressions
}

func (e *ExpressionStorage) GetByID(id int) (db.Expression, error) {
	if expression, ok := e.readyExpressions.Load(id); ok {
		return expression.(db.Expression), nil
	}
	if expression, ok := e.pendingExpressions.Load(id); ok {
		return expression.(db.Expression), nil
	}
	return db.Expression{}, errors.New("expression is not found")
}

// GetNotWorkingExpressions returns all expressions that have Status == ExpressionNotReady.
func (e *ExpressionStorage) GetNotWorkingExpressions() []db.Expression {
	expressions := make([]db.Expression, 0)
	e.pendingExpressions.Range(func(_, value interface{}) bool {
		if value.(db.Expression).Status == db.ExpressionNotReady {
			expressions = append(expressions, value.(db.Expression))
		}
		return true
	})
	return expressions
}

// UpdatePendingExpression updates expression in pendingExpressions and sync with database.
func (e *ExpressionStorage) UpdatePendingExpression(expression db.Expression) error {
	if _, ok := e.pendingExpressions.Load(expression.ID); !ok {
		return errors.New("expression is not found")
	}
	e.pendingExpressions.Store(expression.ID, expression)
	// sync with database
	if err := e.db.UpdateExpression(expression); err != nil {
		return err
	}
	return nil
}

// PendingToReady changes expression from pendingExpressions to readyExpressions and sync with database.
func (e *ExpressionStorage) PendingToReady(expression db.Expression) error {
	if expression.Status != db.ExpressionReady {
		return errors.New("expression is not ready")
	}

	if _, ok := e.pendingExpressions.Load(expression.ID); !ok {
		return errors.New("expression is not found")
	}
	e.pendingExpressions.Delete(expression.ID)
	e.readyExpressions.Store(expression.ID, expression)
	// sync with database
	if err := e.db.UpdateExpression(expression); err != nil {
		return err
	}
	return nil
}

// IsExpressionWorking returns true if expression is in pendingExpressions and has Status == ExpressionWorking.
func (e *ExpressionStorage) IsExpressionWorking(id int) (bool, error) {
	if expression, ok := e.pendingExpressions.Load(id); ok {
		return expression.(db.Expression).Status == db.ExpressionWorking, nil
	}
	return false, errors.New("expression is not found")
}

// IsExpressionNotReady returns true if expression is in pendingExpressions and has Status == ExpressionNotReady.
func (e *ExpressionStorage) IsExpressionNotReady(id int) (bool, error) {
	if expression, ok := e.pendingExpressions.Load(id); ok {
		return expression.(db.Expression).Status == db.ExpressionNotReady, nil
	}
	return false, errors.New("expression is not found")
}

func (e *ExpressionStorage) Delete(id int) error {
	if _, ok := e.readyExpressions.Load(id); ok {
		e.readyExpressions.Delete(id)
	}
	if _, ok := e.pendingExpressions.Load(id); ok {
		e.pendingExpressions.Delete(id)
	}
	// sync with database
	if err := e.db.DeleteExpression(id); err != nil {
		return err
	}
	return nil
}
