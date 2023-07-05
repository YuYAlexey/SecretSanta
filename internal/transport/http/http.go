package http

import (
	"github.com/adYushinW/SecretSanta/internal/app"
	"github.com/gin-gonic/gin"
)

func Service(app *app.App) error {

	c := NewController(app)

	r := gin.Default()

	// Мне кажется нужно выработать какую-нибудь структуру endpoint'ов,
	// так как появилось много сущностей и образуется беспорядочность. Например
	// /auth/register
	// /auth/login
	// /auth/logout
	// /game/start
	// /game/stop
	// /gift/add
	// /gift/set
	// /gift/watch
	// ...etc
	r.GET("/ping", c.Ping)

	r.POST("/register", c.Register)
	r.POST("/login", c.Login)
	r.POST("/secretsanta", c.SecretSanta)
	authRoutes := r.Group("/logged").Use(c.AuthRoute)
	{
		authRoutes.GET("/check_cookie", c.CheckCookie)

		authRoutes.GET("/watch_gift", c.WatchGift)
		// Не совсем понятно название set рядом с add.
		// Можно попробовать переименовать в bindGift и оставить addGift, так понятнее
		authRoutes.POST("set_gift", c.SetGift)
		authRoutes.POST("/add_gift", c.AddGift)
		authRoutes.GET("/gift", c.GiftForWho)

		authRoutes.POST("/start", c.InGame)

		authRoutes.POST("/stop", c.OutGame)

		authRoutes.POST("/logout", c.Logout)
	}

	return r.Run(":8080")
}
