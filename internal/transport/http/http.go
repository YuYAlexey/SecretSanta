package http

import (
	"github.com/adYushinW/SecretSanta/internal/app"
	"github.com/gin-gonic/gin"
)

var Secret = []byte("secret")

func Service(app *app.App) error {

	c := NewController(app)

	r := gin.Default()

	r.GET("/ping", c.Ping)

	r.POST("/register", c.Register)
	r.POST("/login", c.Login)

	//	authRoutes := r.Group("/logged").Use(middlware.AuthRequeired)

	r.GET("/CheckCookie", c.CheckCookie)
	r.GET("/watch_gift", c.WatchGift)
	r.POST("/add_gift", c.AddGift)
	r.POST("/logout", c.Logout)

	return r.Run(":8080")
}
