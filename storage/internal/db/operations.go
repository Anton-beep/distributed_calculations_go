package db

type Operation struct {
	ID           int `db:"id" json:"id"`
	TimeAdd      int `db:"time_add" json:"time_add"`
	TimeSubtract int `db:"time_subtract" json:"time_subtract"`
	TimeDivide   int `db:"time_divide" json:"time_divide"`
	TimeMultiply int `db:"time_multiply" json:"time_mutiply"`
	User         int `db:"user_id" json:"user_id"`
}

func (a *APIDb) GetUserOperations(userID int) (Operation, error) {
	operation := Operation{}
	err := a.db.QueryRow("SELECT * FROM operations WHERE user_id=$1", userID).
		Scan(&operation.ID, &operation.TimeAdd, &operation.TimeSubtract, &operation.TimeDivide, &operation.TimeMultiply, &operation.User)
	if err != nil {
		return operation, err
	}
	return operation, nil
}

func (a *APIDb) AddOperation(operation Operation) (int, error) {
	var id int
	err := a.db.QueryRow("INSERT INTO operations(time_add, time_subtract, time_divide, time_multiply, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id", operation.TimeAdd, operation.TimeSubtract, operation.TimeDivide, operation.TimeMultiply, operation.User).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *APIDb) UpdateOperation(operation Operation) error {
	_, err := a.db.Exec("UPDATE operations SET time_add=$1, time_subtract=$2, time_divide=$3, time_multiply=$4 WHERE id=$5", operation.TimeAdd, operation.TimeSubtract, operation.TimeDivide, operation.TimeMultiply, operation.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *APIDb) DeleteOperation(id int) error {
	_, err := a.db.Exec("DELETE FROM operations WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func (a *APIDb) DeleteByUserId(userId int) error {
	_, err := a.db.Exec("DELETE FROM operations WHERE user_id=$1", userId)
	if err != nil {
		return err
	}
	return nil
}
