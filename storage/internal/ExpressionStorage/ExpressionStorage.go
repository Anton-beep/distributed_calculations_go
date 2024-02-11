package ExpressionStorage

import (
	"errors"
	"go.uber.org/zap"
	"storage/internal/Db"
	"sync"
)

type ExpressionStorage struct {
	readyExpressions   sync.Map
	pendingExpressions sync.Map
	db                 *Db.ApiDb
}

func New(db *Db.ApiDb) *ExpressionStorage {
	e := &ExpressionStorage{
		db: db,
	}

	// check saved data in database
	expressions, err := db.GetAllExpressions()
	if err != nil {
		zap.S().Error(err)
	}
	for _, expression := range expressions {
		if expression.Status == Db.ExpressionReady {
			e.readyExpressions.Store(expression.Id, expression)
		} else {
			e.pendingExpressions.Store(expression.Id, expression)
		}
	}

	return e
}

func (e *ExpressionStorage) Add(expression Db.Expression) (int, error) {
	// sync with database
	newId, err := e.db.AddExpression(expression)
	if err != nil {
		return 0, err
	}
	expression.Id = newId

	e.pendingExpressions.Store(newId, expression)
	return newId, nil
}

func (e *ExpressionStorage) GetAll() []Db.Expression {
	expressions := make([]Db.Expression, 0)
	e.readyExpressions.Range(func(key, value interface{}) bool {
		expressions = append(expressions, value.(Db.Expression))
		return true
	})
	e.pendingExpressions.Range(func(key, value interface{}) bool {
		expressions = append(expressions, value.(Db.Expression))
		return true
	})
	return expressions
}

func (e *ExpressionStorage) GetById(id int) (Db.Expression, error) {
	if expression, ok := e.readyExpressions.Load(id); ok {
		return expression.(Db.Expression), nil
	}
	if expression, ok := e.pendingExpressions.Load(id); ok {
		return expression.(Db.Expression), nil
	}
	return Db.Expression{}, errors.New("expression is not found")
}

func (e *ExpressionStorage) GetNotWorkingExpressions() []Db.Expression {
	expressions := make([]Db.Expression, 0)
	e.pendingExpressions.Range(func(key, value interface{}) bool {
		if value.(Db.Expression).Status == Db.ExpressionNotReady {
			expressions = append(expressions, value.(Db.Expression))
		}
		return true
	})
	return expressions
}

func (e *ExpressionStorage) UpdatePendingExpression(expression Db.Expression) error {
	if _, ok := e.pendingExpressions.Load(expression.Id); !ok {
		return errors.New("expression is not found")
	}
	e.pendingExpressions.Store(expression.Id, expression)
	// sync with database
	if err := e.db.UpdateExpression(expression); err != nil {
		return err
	}
	return nil
}

func (e *ExpressionStorage) PendingToReady(expression Db.Expression) error {
	if expression.Status != Db.ExpressionReady {
		return errors.New("expression is not ready")
	}

	if _, ok := e.pendingExpressions.Load(expression.Id); !ok {
		return errors.New("expression is not found")
	}
	e.pendingExpressions.Delete(expression.Id)
	e.readyExpressions.Store(expression.Id, expression)
	// sync with database
	if err := e.db.UpdateExpression(expression); err != nil {
		return err
	}
	return nil
}

func (e *ExpressionStorage) IsExpressionWorking(id int) (bool, error) {
	if expression, ok := e.pendingExpressions.Load(id); ok {
		return expression.(Db.Expression).Status == Db.ExpressionWorking, nil
	}
	return false, errors.New("expression is not found")
}

func (e *ExpressionStorage) IsExpressionNotReady(id int) (bool, error) {
	if expression, ok := e.pendingExpressions.Load(id); ok {
		return expression.(Db.Expression).Status == Db.ExpressionNotReady, nil
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
