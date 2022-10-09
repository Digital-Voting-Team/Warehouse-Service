package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/sijms/go-ora/v2"
	"log"
)

func main() {
	db, err := connectToOracle("oracle://SYS:admin@localhost:1521/XEPDB1")

	if err != nil {
		log.Panicln(err)
	}

	query, err := db.Queryx("SELECT table_name  from all_tables where owner = 'WAREHOUSE'")

	if err != nil {
		log.Panicln(err)
	}

	var s string

	for query.Next() {
		err = query.Scan(&s)

		if err != nil {
			log.Panicln(err)
		}

		log.Println(s)
	}
}

func connectToOracle(connectionString string) (*sqlx.DB, error) {
	db, err := sqlx.Open("oracle", connectionString)

	if err != nil {
		return nil, err
	}

	return db, err
}
