package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Session struct {
	UserID    string
	ExpiredAt time.Time
}

func AuthRequeired(ctx *gin.Context) {
	session := sessions.Default(ctx)
	user := session.Get("user")
	if user == nil {
		log.Println("User not logged in")
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		ctx.Abort()
		return
	}
	ctx.Next()
}
