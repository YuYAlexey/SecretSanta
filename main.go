package main

import (
	"github.com/adYushinW/SecretSanta/internal/app"
	"github.com/adYushinW/SecretSanta/internal/db"
	http "github.com/adYushinW/SecretSanta/internal/transport"
)

func main() {

	db, err := db.New()
	if err != nil {
		panic(err)
	}

	app := app.New(db)

	if err := http.Service(app); err != nil {
		panic(err)
	}
}
