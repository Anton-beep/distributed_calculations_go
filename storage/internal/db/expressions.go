package db

const (
	ExpressionNotReady = 0
	ExpressionWorking  = 1
	ExpressionReady    = 2
	ExpressionError    = 3
)

type Expression struct {
	ID                 int     `db:"id" json:"id"`
	Value              string  `db:"value" json:"value"`
	Answer             float64 `db:"answer" json:"answer"`
	Logs               string  `db:"logs" json:"logs"`
	Status             int     `db:"ready" json:"ready"` // 0 - not ready, 1 - working, 2 - ready, 3 - error
	AliveExpiresAt     int     `db:"alive_expires_at" json:"alive_expires_at"`
	CreationTime       string  `db:"creation_time" json:"creation_time"`
	EndCalculationTime string  `db:"end_calculation_time" json:"end_calculation_time"`
	Servername         string  `db:"server_name" json:"server_name"`
	User               int     `db:"user_id" json:"user_id"`
}

func (a *APIDb) GetAllExpressions() ([]Expression, error) {
	expressions := make([]Expression, 0)
	rows, err := a.db.Query("SELECT * FROM expressions")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		expression := Expression{}
		err = rows.Scan(&expression.ID, &expression.Value, &expression.Answer, &expression.Logs, &expression.Status,
			&expression.AliveExpiresAt, &expression.CreationTime, &expression.EndCalculationTime, &expression.Servername,
			&expression.User)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expression)
	}

	return expressions, nil
}

func (a *APIDb) GetExpressionByID(id int) (Expression, error) {
	expression := Expression{}
	err := a.db.QueryRow("SELECT * FROM expressions WHERE id=$1", id).
		Scan(&expression.ID, &expression.Value, &expression.Answer, &expression.Logs, &expression.Status,
			&expression.AliveExpiresAt, &expression.CreationTime, &expression.EndCalculationTime, &expression.Servername,
			&expression.User)
	if err != nil {
		return expression, err
	}
	return expression, nil
}

func (a *APIDb) AddExpression(expression Expression) (int, error) {
	var id int
	err := a.db.QueryRow("INSERT INTO expressions(value, answer, logs, ready, alive_expires_at, creation_time,"+
		" end_calculation_time, server_name, user_id) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
		expression.Value, expression.Answer, expression.Logs, expression.Status, expression.AliveExpiresAt,
		expression.CreationTime, expression.EndCalculationTime, expression.Servername, expression.User).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *APIDb) UpdateExpression(expression Expression) error {
	_, err := a.db.Exec("UPDATE expressions SET value=$1, answer=$2, logs=$3, ready=$4, alive_expires_at=$5,"+
		" creation_time=$6, end_calculation_time=$7, server_name=$8, user_id=$9 WHERE id=$10",
		expression.Value, expression.Answer, expression.Logs, expression.Status, expression.AliveExpiresAt,
		expression.CreationTime, expression.EndCalculationTime, expression.Servername, expression.User, expression.ID)
	return err
}

func (a *APIDb) GetLastID() int {
	var id int
	err := a.db.QueryRow("SELECT MAX(id) FROM expressions").Scan(&id)
	if err != nil {
		return 0
	}
	return id
}

func (a *APIDb) DeleteExpression(id int) error {
	_, err := a.db.Exec("DELETE FROM expressions WHERE id = $1", id)
	return err
}
