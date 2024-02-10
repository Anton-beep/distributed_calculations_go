package ExpressionStorage

import (
	"storage/internal/Db"
	"sync"
)

type ExpressionStorage struct {
	expressions map[int]Db.Expression
	mu          sync.Mutex
}

func (e *ExpressionStorage) Add(expression Db.Expression) {

}

func (e *ExpressionStorage) GetAll() []Db.Expression {
	return nil
}

func (e *ExpressionStorage) GetById(id int) Db.Expression {
	return Db.Expression{}
}
