package http

import (
	"net/http"

	"github.com/adYushinW/SecretSanta/internal/db"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	db *db.Database
}

func NewController(db *db.Database) *Controller {
	return &Controller{db: db}
}

func (c *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "PONG")
}
