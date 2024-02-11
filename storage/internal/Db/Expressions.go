package Db

import (
	"database/sql"
	"go.uber.org/zap"
)

const (
	ExpressionNotReady = 0
	ExpressionWorking  = 1
	ExpressionReady    = 2
)

type Expression struct {
	Id     int     `db:"id" json:"id"`
	Value  string  `db:"value" json:"value"`
	Answer float64 `db:"answer" json:"answer"`
	Logs   string  `db:"logs" json:"logs"`
	Status int     `db:"ready" json:"ready"` // 0 - not ready, 1 - working, 2 - ready
}

func (a *ApiDb) GetAllExpressions() ([]Expression, error) {
	expressions := make([]Expression, 0)
	rows, err := a.db.Query("SELECT * FROM expressions")
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			zap.S().Error(err)
		}
	}(rows)

	for rows.Next() {
		expression := Expression{}
		err := rows.Scan(&expression.Id, &expression.Value, &expression.Answer, &expression.Logs, &expression.Status)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expression)
	}
	return expressions, nil
}

func (a *ApiDb) GetExpressionById(id int) (Expression, error) {
	expression := Expression{}
	err := a.db.QueryRow("SELECT * FROM expressions WHERE id=$1", id).
		Scan(&expression.Id, &expression.Value, &expression.Answer, &expression.Logs, &expression.Status)
	if err != nil {
		return expression, err
	}
	return expression, nil
}

func (a *ApiDb) AddExpression(expression Expression) (int, error) {
	var id int
	err := a.db.QueryRow("INSERT INTO expressions(value, answer, logs, ready) VALUES($1, $2, $3, $4) RETURNING id",
		expression.Value, expression.Answer, expression.Logs, expression.Status).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *ApiDb) UpdateExpression(expression Expression) error {
	_, err := a.db.Exec("UPDATE expressions SET value=$1, answer=$2, logs=$3, ready=$4 WHERE id=$5",
		expression.Value, expression.Answer, expression.Logs, expression.Status, expression.Id)
	return err
}

func (a *ApiDb) GetLastId() (int, error) {
	var id int
	err := a.db.QueryRow("SELECT MAX(id) FROM expressions").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *ApiDb) DeleteExpression(id int) error {
	_, err := a.db.Exec("DELETE FROM expressions WHERE id = $1", id)
	return err
}
