package Db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"os"
)

type ApiDb struct {
	db *sql.DB
}

func New() (*ApiDb, error) {
	// this is for test purposes
	err := godotenv.Load("../.env")
	if err != nil {
		zap.S().Warn(err)
	}

	//zap.S().Info(os.Getenv("POSTGRESQL_USER"), os.Getenv("POSTGRESQL_PASSWORD"), os.Getenv("POSTGRESQL_HOST"), os.Getenv("POSTGRESQL_PORT"), os.Getenv("POSTGRESQL_NAME"))
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

	a := &ApiDb{db}
	// database will be wiped if in .env variable RESET_POSTGRESQL = TRUE
	a.ResetDatabase()
	return a, nil
}

func (a *ApiDb) ResetDatabase() {
	if os.Getenv("RESET_POSTGRESQL") != "TRUE" {
		return
	}

	zap.S().Warn("resetting database")

	_, err := a.db.Exec("DROP TABLE IF EXISTS expressions")
	if err != nil {
		zap.S().Fatal(err)
	}
	_, err = a.db.Exec("CREATE TABLE expressions (id SERIAL PRIMARY KEY, value TEXT, answer FLOAT, logs TEXT, ready INT)")
	if err != nil {
		zap.S().Fatal(err)
	}
}
