package http

import (
	"net/http"

	"github.com/adYushinW/SecretSanta/internal/db"
	"github.com/adYushinW/SecretSanta/internal/model"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	db *db.Database
}

func NewController(db *db.Database) *Controller {
	// REVIEW: Почему Contraller по значет от базе дыннх!? Где слой бизнес логики?
	// Переделать!!!!! 
	return &Controller{db: db}
}

func (c *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "PONG")
}

func (c *Controller) Register(ctx *gin.Context) {

	user := model.Users{}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, "register")
}

func (c *Controller) Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "login")
}
