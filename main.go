package main

import (
	"github.com/adYushinW/SecretSanta/internal/transport/http"
)

func main() {

	if err := http.Service(); err != nil {
		panic(err)
	}
}
