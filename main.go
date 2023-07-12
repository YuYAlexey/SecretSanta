package main

import (
	"github.com/YuYAlexey/SecretSanta/internal/app"
	"github.com/YuYAlexey/SecretSanta/internal/db"
	"github.com/YuYAlexey/SecretSanta/internal/transport/http"
)

func main() {

	conn, err := db.NewConnect()
	if err != nil {
		panic(err)
	}

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
