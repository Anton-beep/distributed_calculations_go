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

func GetSQLFromFile(name string) (string, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (a *APIDb) ResetDatabase() {
	command, err := GetSQLFromFile("sqlScripts/resetDB.sql")
	if err != nil {
		zap.S().Fatal(err)
	}
	_, err = a.db.Exec(command)
	if err != nil {
		zap.S().Fatal(err)
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

	err = a.CheckFields("expressions", correctFieldsExpressions)
	if err != nil {
		return false, err
	}
	err = a.CheckFields("users", correctFieldsExpressionsUsers)
	if err != nil {
		return false, err
	}
	return true, nil
}
