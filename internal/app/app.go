package app

import (
	"github.com/adYushinW/SecretSanta/internal/db"
	"github.com/adYushinW/SecretSanta/internal/model"
)

type App struct {
	db db.Database
}

func New(db db.Database) *App {
	return &App{
		db: db,
	}
}

func (app *App) AddUser(login string, password string, firstName string, lastName string, sex string, age uint64) (bool, error) {
	return app.db.AddUser(login, password, firstName, lastName, sex, age)
}

func (app *App) Login(login string, password string) (bool, error) {
	return app.db.Login(login, password)
}

func (app *App) Participate(login string, isPlay bool) error {
	return app.db.Participate(login, isPlay)
}

func (app *App) WatchGift() ([]*model.Gift, error) {
	return app.db.WatchGift()
}

func (app *App) AddGift(name string, link string, description string) (bool, error) {
	return app.db.AddGift(name, link, description)
}

func (app *App) SetGift(login string, gift string) (bool, error) {
	return app.db.SetGift(login, gift)
}

func (app *App) GiftForWho(login string) (interface{}, error) {
	return app.db.GiftForWho(login)
}

func (app *App) SecretSanta() (peopleID []uint64, err error) {
	return app.db.SecretSanta()
}

func (app *App) SetGiverRecipient(giverRecipient map[uint64]*model.GiverRecipient) (bool, error) {
	return app.db.SetGiverRecipient(giverRecipient)
}
