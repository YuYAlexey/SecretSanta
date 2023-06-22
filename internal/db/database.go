package db

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/adYushinW/SecretSanta/internal/model"
	"github.com/jackc/pgx/v4"
)

var (
	ErrRowNotAdded  = errors.New("row not added")
	ErrConvertToSql = errors.New("failed to convert to SQL error")
	ErrBuildQuery   = errors.New("failed build query error")
)

type Database interface {
	AddUser(login string, password string, firstName string, lastName string, sex string, age uint64) (bool, error)
	Login(login string, password string) (bool, error)
	WatchGift() ([]*model.Gift, error)
	AddGift(name string, link string, description string) (bool, error)
}

type database struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) (Database, error) {
	return &database{
		conn: conn,
	}, nil
}

func (db *database) AddUser(login string, password string, firstName string, lastName string, sex string, age uint64) (bool, error) {
	qb := sq.Insert("users").
		Columns("login", "password", "first_name", "last_name", "sex", "age").
		Values(login, password, firstName, lastName, sex, age).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := qb.ToSql()
	if err != nil {
		err := fmt.Errorf("%s: %w", ErrConvertToSql, err)
		return false, err
	}

	row, err := db.conn.Exec(context.Background(), sql, args...)
	if err != nil {
		err := fmt.Errorf("%s: %w", ErrBuildQuery, err)
		return false, err
	}

	if row.RowsAffected() != 1 {
		return false, fmt.Errorf("%s", ErrRowNotAdded)
	}

	return true, nil

}

func (db *database) Login(login string, password string) (bool, error) {
	qb := sq.Select("id", "login", "password", "first_name, last_name, sex, age").
		From("users").
		Where("login = ?", login).
		Where("password = ?", password).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := qb.ToSql()
	if err != nil {
		err := fmt.Errorf("%s: %w", ErrConvertToSql, err)
		return false, err
	}

	row, err := db.conn.Query(context.Background(), sql, args...)
	if err != nil {
		err := fmt.Errorf("%s: %w", ErrBuildQuery, err)
		return false, err
	}

	user := new(model.Users)

	if err = row.Scan(&user.Login, &user.Password); err != pgx.ErrNoRows {
		err := fmt.Errorf("something went wrong while reading row error: %w", err)
		return false, err
	}

	return true, nil
}

func (db *database) WatchGift() ([]*model.Gift, error) {
	qb := sq.Select("id", "name", "link", "description", "is_selected").
		From("gift").
		Where("is_selected = false").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := qb.ToSql()
	if err != nil {
		err := fmt.Errorf("%s: %w", ErrConvertToSql, err)
		return nil, err
	}

	rows, err := db.conn.Query(context.Background(), sql, args...)
	if err != nil {
		err := fmt.Errorf("%s: %w", ErrBuildQuery, err)
		return nil, err
	}
	defer rows.Close()

	result := make([]*model.Gift, 0)

	for rows.Next() {
		gift := new(model.Gift)

		if err = rows.Scan(&gift.ID, &gift.Name, &gift.Link, &gift.Description, &gift.IsSelected); err != nil {
			continue
		}

		result = append(result, gift)
	}

	if err := rows.Err(); err != nil {
		err := fmt.Errorf("something went wrong while reading row error: %w", err)
		return result, err
	}

	return result, nil
}

func (db *database) AddGift(name string, link string, description string) (bool, error) {
	qb := sq.Insert("gift").
		Columns("name", "link", "description").
		Values(name, link, description).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := qb.ToSql()
	if err != nil {
		err := fmt.Errorf("%s: %w", ErrConvertToSql, err)
		return false, err
	}

	row, err := db.conn.Exec(context.Background(), sql, args...)
	if err != nil {
		err := fmt.Errorf("%s: %w", ErrBuildQuery, err)
		return false, err
	}

	if row.RowsAffected() != 1 {
		return false, fmt.Errorf("%s", ErrRowNotAdded)
	}

	return true, nil
}
