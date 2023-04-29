package storage

import "database/sql"
import _ "github.com/jackc/pgx/v5/stdlib"

func Connect(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return conn, err
}
