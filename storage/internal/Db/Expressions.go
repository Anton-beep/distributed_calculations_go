package Db

import "log"

type Expression struct {
	Id     int    `db:"id" json:"id"`
	Value  string `db:"value" json:"value"`
	Answer int    `db:"answer" json:"answer"`
	Logs   string `db:"logs" json:"logs"`
	Ready  bool   `db:"ready" json:"ready"`
}

func (a *ApiDb) WriteExpression() {
	_, err := a.db.Exec("INSERT INTO expressions (value, answer, logs) VALUES ($1, $2, $3)", "test", 1, "test")
	if err != nil {
		log.Fatalln("error while writing expression: ", err)
	}
}
