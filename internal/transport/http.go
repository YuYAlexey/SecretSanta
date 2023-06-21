package http

import (
	"github.com/adYushinW/SecretSanta/internal/app"
	"github.com/gin-gonic/gin"
)

func Service(app *app.App) error {

	c := NewController(app)

	r := gin.Default()

	r.GET("/ping", c.Ping)

	r.POST("/register", c.Register)
	r.POST("/login", c.Login)

	r.GET("/watch_gift", c.WatchGift)
	r.POST("/add_gift", c.AddGift)

	return r.Run(":8080")
}
