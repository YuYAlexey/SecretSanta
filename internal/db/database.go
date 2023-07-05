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
	SetGift(login string, gift string) (bool, error)
	GiftForWho(login string) (*recipientGift, error)
	SecretSanta() (peopleID []uint64, err error)
	Participate(login string, isPlay bool) error
	SetGiverRecipient(giverRecipient map[uint64]*model.GiverRecipient) (bool, error)
}

type recipientGift struct {
	FirstName string
	LastName  string
	GiftName  string
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
	qb := sq.Select("login", "password").
		From("users").
		Where("login = ? AND password = ?", login, password).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := qb.ToSql()
	if err != nil {
		err := fmt.Errorf("%s: %w", ErrConvertToSql, err)
		return false, err
	}

	row := db.conn.QueryRow(context.Background(), sql, args...)

	user := new(model.Users)
	if err := row.Scan(&user.Login, &user.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
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

func (db *database) GiftForWho(login string) (*recipientGift, error) {

	sqb := sq.Select("gr.recipient").
		From("users u").
		LeftJoin("giverrecipient gr ON u.id = gr.giver").
		Where("login = ?", login).PlaceholderFormat(sq.Dollar)

	qb := sq.Select("u.first_name", "u.last_name", "g.name").
		From("users u").
		LeftJoin("gift g ON g.id = u.gift").
		Where("u.id = ?", sqb).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := qb.ToSql()
	if err != nil {
		err := fmt.Errorf("%s: %w", ErrConvertToSql, err)
		return nil, err
	}

	row := db.conn.QueryRow(context.Background(), sql, args...)

	user := new(recipientGift)
	if err = row.Scan(&user.FirstName, &user.LastName, &user.GiftName); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		err := fmt.Errorf("something went wrong while reading row error: %w", err)
		return nil, err
	}

	return user, nil
}

func (db *database) SetGift(login string, gift string) (bool, error) {

	sqb := sq.Select("id").
		From("gift").
		Where("lower(name) = ?", gift)
	//Не уменьшает текст для поя Name, хотя в SQL работает

	qb := sq.Update("users").
		Set("gift", sqb).
		Where("login = ?", login).
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

func (db *database) SecretSanta() (peopleID []uint64, err error) {
	var players int

	qb := sq.Select("id").
		From("users").
		Where("is_player = TRUE AND gift NOTNULL").
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

	result := make([]uint64, 0, players)

	for rows.Next() {
		var playerId uint64
		if err := rows.Scan(&playerId); err != nil {
			continue
		}
		result = append(result, playerId)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (db *database) SetGiverRecipient(giverRecipient map[uint64]*model.GiverRecipient) (bool, error) {

	qb := sq.Insert("giverrecipient").
		Columns("giver", "recipient")

	for i := range giverRecipient {
		qb = qb.Values(giverRecipient[i].Giver, giverRecipient[i].Recipient)
	}

	qb = qb.PlaceholderFormat(sq.Dollar)

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

	if row.RowsAffected() == 0 {
		return false, fmt.Errorf("%s", ErrRowNotAdded)
	}

	return true, nil
}

func (db *database) Participate(login string, isPlay bool) error {
	qb := sq.Update("users").
		Set("is_player", isPlay).
		Where("login = ?", login).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := qb.ToSql()
	if err != nil {
		err := fmt.Errorf("%s: %w", ErrConvertToSql, err)
		return err
	}

	row, err := db.conn.Exec(context.Background(), sql, args...)
	if err != nil {
		err := fmt.Errorf("%s: %w", ErrBuildQuery, err)
		return err
	}

	if row.RowsAffected() != 1 {
		return fmt.Errorf("%s", ErrRowNotAdded)
	}

	return nil
}
