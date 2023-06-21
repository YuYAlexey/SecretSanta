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

func (app *App) WatchGift() ([]*model.Gift, error) {
	return app.db.WatchGift()
}

func (app *App) AddGift(name string, link string, description string) (bool, error) {
	return app.db.AddGift(name, link, description)
}