package db

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/adYushinW/SecretSanta/internal/model"
	"github.com/jackc/pgx/v4"
)

type Database interface {
	// REVIEW: не правильное навзание перпменных с точки зерения Go. Перименовать: first_name -> firstName last_name -> lastName
	AddUser(login string, password string, first_name string, last_name string, sex string, age uint8) ([]*model.Users, error)
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

// REVIEW: не правильное навзание перпменных с точки зерения Go. Перименовать: first_name -> firstName last_name -> lastName
func (db *database) AddUser(login string, password string, first_name string, last_name string, sex string, age uint8) ([]*model.Users, error) {
	// REVIEW: не вижу необходимости возвращаться slice пользователей. Хватить и одного!

	// REVIEW: ты в RETURNIGN возвращаешь всю информвцию? можно же сделать запрос через Exec и поянть добавились данные или нет и не нужен RETURNIGN
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

	// REVIEW: defer с cloese не стоит отделять строкой от блока где он создан
	defer row.Close()

	result := make([]*model.Users, 0)

	user := new(model.Users)

	// REVIEW: Я бы перписал сканирование так
	//	if err := row.Scan(&user.Id); errors.Is(err, pgx.ErrNoRows) {
	//		return nil, err
	//	}

	// REVIEW: А какой  ID Ты пытаешься разпарсить если ты запросил login, password, first_name, last_name, sex, age - там нет поле id
	err = row.Scan(&user.Id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	result = append(result, user)

	return result, nil
}
