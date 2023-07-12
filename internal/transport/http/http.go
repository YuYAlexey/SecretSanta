package http

import (
	"github.com/YuYAlexey/SecretSanta/internal/app"
	"github.com/gin-gonic/gin"
)

func Service(app *app.App) error {

	c := NewController(app)

	r := gin.Default()

	r.GET("/ping", c.Ping)

	r.POST("/register", c.Register)
	r.POST("/login", c.Login)
	r.POST("/secretsanta", c.SecretSanta)
	authRoutes := r.Group("/logged").Use(c.AuthRoute)
	{
		authRoutes.GET("/check_cookie", c.CheckCookie)

		authRoutes.GET("/watch_gift", c.WatchGift)
		authRoutes.POST("set_gift", c.SetGift)
		authRoutes.POST("/add_gift", c.AddGift)
		authRoutes.GET("/gift", c.GiftForWho)

		authRoutes.POST("/start", c.InGame)

		authRoutes.POST("/stop", c.OutGame)

		authRoutes.POST("/logout", c.Logout)
	}

	return r.Run(":8080")
}
