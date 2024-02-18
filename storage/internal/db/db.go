package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	// postgresql driver.
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"os"
)

type APIDb struct {
	db *sql.DB
}

func New() (*APIDb, error) {
	// this is for test purposes
	err := godotenv.Load("../.env")
	if err != nil {
		zap.S().Warn("this warning is normal if you are running the server in production")
		zap.S().Warn(err)
	}

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("POSTGRESQL_USER"),
		os.Getenv("POSTGRESQL_PASSWORD"),
		os.Getenv("POSTGRESQL_HOST"),
		os.Getenv("POSTGRESQL_PORT"),
		os.Getenv("POSTGRESQL_NAME"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// check connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	a := &APIDb{db}
	// database will be wiped if in .env variable RESET_POSTGRESQL = TRUE
	if os.Getenv("RESET_POSTGRESQL") != "TRUE" {
		if res, _ := a.IsDBCorrect(); res {
			zap.S().Info("database is correct")
		} else {
			zap.S().Warn("database is not correct, resetting")
			a.ResetDatabase()
		}
	} else {
		zap.S().Warn("resetting database because of RESET_POSTGRESQL = TRUE")
		a.ResetDatabase()
	}
	return a, nil
}

func (a *APIDb) ResetDatabase() {
	_, err := a.db.Exec("DROP TABLE IF EXISTS expressions")
	if err != nil {
		zap.S().Fatal(err)
	}
	_, err = a.db.Exec("CREATE TABLE expressions (id SERIAL PRIMARY KEY, value TEXT, answer FLOAT, logs TEXT," +
		" ready INT, alive_expires_at BIGINT, creation_time TEXT, end_calculation_time TEXT, server_name TEXT)")
	if err != nil {
		zap.S().Fatal(err)
	}
}

func (a *APIDb) IsDBCorrect() (bool, error) {
	var exists bool
	err := a.db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE tablename = 'expressions')").Scan(&exists)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, errors.New("table expressions does not exist")
	}

	// Check if table has correct fields
	rows, err := a.db.Query("SELECT column_name FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = 'expressions'")
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// Define the correct fields
	correctFields := map[string]bool{
		"id": true, "value": true, "answer": true, "logs": true,
		"ready": true, "alive_expires_at": true, "creation_time": true,
		"end_calculation_time": true, "server_name": true,
	}

	for rows.Next() {
		var columnName string
		if err = rows.Scan(&columnName); err != nil {
			return false, err
		}
		if _, ok := correctFields[columnName]; !ok {
			return false, fmt.Errorf("unexpected column %s in table expressions", columnName)
		}
		delete(correctFields, columnName)
	}

	if err = rows.Err(); err != nil {
		return false, err
	}

	if len(correctFields) > 0 {
		return false, fmt.Errorf("missing columns in table expressions: %v", correctFields)
	}

	return true, nil
}
