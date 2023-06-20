package db

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/adYushinW/SecretSanta/internal/model"
	"github.com/jackc/pgx/v4"
)

type Database interface {
	AddUser(login string, password string, first_name string, last_name string, sex string, age uint8) ([]*model.Users, error)
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

func (db *database) AddUser(login string, password string, first_name string, last_name string, sex string, age uint8) ([]*model.Users, error) {
	qb := sq.Insert("users").
		Columns("login", "password", "first_name", "last_name", "sex", "age").
		Values(login, password, first_name, last_name, sex, age).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNIGN login, password, first_name, last_name, sex, age")

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	row, err := db.conn.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	result := make([]*model.Users, 0)

	user := new(model.Users)

	err = row.Scan(&user.Id)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	result = append(result, user)

	return result, nil

}
