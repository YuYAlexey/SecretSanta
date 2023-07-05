package main

import (
	"github.com/adYushinW/SecretSanta/internal/app"
	"github.com/adYushinW/SecretSanta/internal/db"
	"github.com/adYushinW/SecretSanta/internal/transport/http"
)

func main() {

	conn, err := db.NewConnect()
	if err != nil {
		panic(err)
	}
	// В переменной db теперь хранится реализация структуры Database
	// Если захочешь где-то ниже вызвать функции из `package db`, файла `connect.go`, то будут проблемы
	// Лучше не переопределять сущности, так как в дальнейшем возможно придётся рефакторить и
	// переименовывать переменную либо до её определения, либо после
	db, err := db.New(conn)
	if err != nil {
		panic(err)
	}

	app := app.New(db)
	if err != nil {
		panic(err)
	}

	if err := http.Service(app); err != nil {
		panic(err)
	}
}
