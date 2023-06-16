package db

import "github.com/jackc/pgx/v4"

type Database interface {
}

type database struct {
	conn *pgx.Conn
}

func New() (Database, error) {
	conn, err := newConnect()
	if err != nil {
		return nil, err
	}

	return &database{
		conn: conn,
	}, nil
}
