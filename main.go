package main

import (
	"github.com/adYushinW/SecretSanta/internal/app"
	"github.com/adYushinW/SecretSanta/internal/db"
	http "github.com/adYushinW/SecretSanta/internal/transport"
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
