package http

import (
	"github.com/adYushinW/SecretSanta/internal/db"
	"github.com/gin-gonic/gin"
)

func Service(db *db.Database) error {

	c := NewController(db)

	r := gin.Default()

	r.GET("/ping", c.Ping)

	return r.Run(":8080")
}
