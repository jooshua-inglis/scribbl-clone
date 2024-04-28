package data

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB = nil

func GetDb() *sqlx.DB {
	if db == nil {
		slog.Debug("connecting to database...")
		connStr := "user=postgres dbname=postgres sslmode=disable password=password"
		_db, err := sqlx.Open("postgres", connStr)
		if err != nil {
			slog.Error(err.Error())
		}
		slog.Debug("successfully connected to database")
		db = _db
	}
	return db
}
