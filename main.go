package main

import (
	"github.com/adYushinW/SecretSanta/internal/db"
	http "github.com/adYushinW/SecretSanta/internal/transport"
)

func main() {

	conn, err := db.New()
	if err != nil {
		panic(err)
	}

	if err := http.Service(&conn); err != nil {
		panic(err)
	}
}
