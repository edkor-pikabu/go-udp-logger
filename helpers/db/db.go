package db

import (
	"database/sql"
	"server/helpers/config"
	"time"
)

func New(config *config.Config) *sql.DB {
	db, err := sql.Open("mysql", config.MysqlDsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}
