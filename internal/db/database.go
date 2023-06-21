package db

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/adYushinW/SecretSanta/internal/model"
	"github.com/jackc/pgx/v4"
)

// REVIEW: пример сообщения об ошибке
// var (
// 	ErrRowNotAdded = errors.New("row not added")
// )

type Database interface {
	AddUser(login string, password string, firstName string, lastName string, sex string, age uint64) (bool, error)
	WatchGift() ([]*model.Gift, error)
	AddGift(name string, link string, description string) (bool, error)
}

type database struct {
	conn *pgx.Conn
}

func New() (Database, error) {
	// REVIEW: не стоит использовать инициализировать подключение на слое с базой
	// Лучше в функцию передать соединение, а соединение создать в main файле
	// Пример: New(conn *pgx.Conn)
	conn, err := newConnect()
	if err != nil {
		return nil, err
	}

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
		// REVIEW: оберни ошибку и напиши пояснение к ней
		return false, err
	}

	row, err := db.conn.Exec(context.Background(), sql, args...)
	if err != nil {
		// REVIEW: оберни ошибку и напиши пояснение к ней
		return false, err
	}

	if row.RowsAffected() != 1 {
		// REVIEW: в GO соощкния об ошибке начинабтся с буквы нижнего регистра
		// fmt.Errorf("row not added")

		// REVIEW: Стоит избегать дублирование сообщение об ошибке,
		// у тебя такое сообщение есть в методе AddGift.
		// Я бы вынес такую ошибку в глобальну переменную внутри пакет (пример в начале файла)
		return false, fmt.Errorf("Row not added")
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
		// REVIEW: оберни ошибку и напиши пояснение к ней
		return nil, err
	}

	rows, err := db.conn.Query(context.Background(), sql, args...)
	if err != nil {
		// REVIEW: оберни ошибку и напиши пояснение к ней
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
		// REVIEW: оберни ошибку и напиши пояснение к ней
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
		// REVIEW: оберни ошибку и напиши пояснение к ней
		return false, err
	}

	row, err := db.conn.Exec(context.Background(), sql, args...)
	if err != nil {
		// REVIEW: оберни ошибку и напиши пояснение к ней
		return false, err
	}

	if row.RowsAffected() != 1 {
		// REVIEW: Стоит избегать дублирование сообщение об ошибке,
		// у тебя такое сообщение есть в методе AddUser.
		// Я бы вынес такую ошибку в глобальну переменную внутри пакет (пример в начале файла)
		return false, fmt.Errorf("Row not added")
	}

	return true, nil
}
