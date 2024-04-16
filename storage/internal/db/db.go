package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"time"

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

	var db *sql.DB

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("POSTGRESQL_USER"),
		os.Getenv("POSTGRESQL_PASSWORD"),
		os.Getenv("POSTGRESQL_HOST"),
		os.Getenv("POSTGRESQL_PORT"),
		os.Getenv("POSTGRESQL_NAME"))

	for i := 0; i < 3; i++ {
		zap.S().Warn(fmt.Sprintf("Attempt %d: Connecting to database: %v", i+1, connStr))
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			zap.S().Warn(fmt.Sprintf("Failed to connect to database: %v", err))
			if i < 2 { // Don't sleep after the last attempt
				time.Sleep(5 * time.Second)
			}
		} else {
			break
		}
	}

	a := &APIDb{db}
	// database will be wiped if in .env variable RESET_POSTGRESQL = TRUE
	if os.Getenv("RESET_POSTGRESQL") != "TRUE" {
		if res, err := a.IsDBCorrect(); res {
			zap.S().Info("database is correct")
		} else {
			zap.S().Warn(fmt.Sprintf("database is not correct: %v, resetting", err))
			a.ResetDatabase()
		}
	} else {
		zap.S().Warn("resetting database because of RESET_POSTGRESQL = TRUE")
		a.ResetDatabase()
	}
	return a, nil
}

func (a *APIDb) ResetDatabase() {
	for i := 0; i < 5; i++ {
		zap.S().Warn(fmt.Sprintf("Attempt %d: Resetting database", i+1))
		command := "DROP TABLE IF EXISTS expressions;\nDROP TABLE IF EXISTS users;\n\nCREATE TABLE users\n(\n    id       SERIAL PRIMARY KEY,\n    login    TEXT,\n    password TEXT\n);\n\nCREATE TABLE expressions\n(\n    id                   SERIAL PRIMARY KEY,\n    value                TEXT,\n    answer               FLOAT,\n    logs                 TEXT,\n    ready                INT,\n    alive_expires_at     BIGINT,\n    creation_time        TEXT,\n    end_calculation_time TEXT,\n    server_name          TEXT,\n    user_id              INT,\n    CONSTRAINT fk_user\n        FOREIGN KEY (user_id)\n            REFERENCES users (id)\n);\n\nCREATE TABLE operations\n(\n    id            SERIAL PRIMARY KEY,\n    time_add      INT,\n    time_subtract INT,\n    time_divide   INT,\n    time_multiply INT,\n    user_id       INT,\n    CONSTRAINT fk_user\n        FOREIGN KEY (user_id)\n            REFERENCES users (id)\n);"
		_, err := a.db.Exec(command)
		if err != nil {
			zap.S().Warn(fmt.Sprintf("Failed to reset database: %v", err))
			if i < 2 { // Don't sleep after the last attempt
				time.Sleep(2 * time.Second)
			}
		} else {
			break
		}
	}
}

func (a *APIDb) CheckFields(tableName string, correctFields []string) error {
	rows, err := a.db.Query("SELECT column_name FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = $1", tableName)
	if err != nil {
		return err
	}
	defer rows.Close()

	fieldsMap := make(map[string]bool)
	for _, el := range correctFields {
		fieldsMap[el] = true
	}

	for rows.Next() {
		var columnName string
		if err = rows.Scan(&columnName); err != nil {
			return err
		}
		if _, ok := fieldsMap[columnName]; !ok {
			return fmt.Errorf("unexpected column %s in table expressions", columnName)
		}
		delete(fieldsMap, columnName)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	if len(fieldsMap) > 0 {
		return fmt.Errorf("missing columns in table expressions: %v", correctFields)
	}

	return nil
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

	correctFieldsExpressions := []string{
		"id", "value", "answer", "logs", "ready", "alive_expires_at", "creation_time", "end_calculation_time", "server_name", "user_id",
	}
	correctFieldsExpressionsUsers := []string{
		"id", "login", "password",
	}
	correctFieldsOperarions := []string{
		"id", "time_add", "time_subtract", "time_divide", "time_multiply", "user_id",
	}

	err = a.CheckFields("expressions", correctFieldsExpressions)
	if err != nil {
		return false, err
	}
	err = a.CheckFields("users", correctFieldsExpressionsUsers)
	if err != nil {
		return false, err
	}
	err = a.CheckFields("operations", correctFieldsOperarions)
	if err != nil {
		return false, err
	}
	return true, nil
}
