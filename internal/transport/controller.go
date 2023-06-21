package http

import (
	"net/http"
	"strconv"

	"github.com/adYushinW/SecretSanta/internal/app"
	"github.com/adYushinW/SecretSanta/internal/model"
	"github.com/gin-gonic/gin"
)

// REVIEW: почему файл не в пакете controller!? перенести

type Controller struct {
	app *app.App
}

func NewController(app *app.App) *Controller {
	return &Controller{app: app}
}

func (c *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "PONG")
}

func (c *Controller) Register(ctx *gin.Context) {
	var age uint64
	var err error

	user := model.Users{}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if user.Age != "" {
		age, err = strconv.ParseUint(user.Age, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "Age must be a number")
			return
		}
	}

	newuser, err := c.app.AddUser(user.Login, user.Password, user.FirstName, user.LastName, user.Sex, age)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	ctx.JSON(http.StatusOK, newuser)
}

func (c *Controller) Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "login")
}

func (c *Controller) WatchGift(ctx *gin.Context) {
	gift, err := c.app.WatchGift()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Bad Request")
		return
	}
	ctx.JSON(http.StatusOK, gift)
}

func (c *Controller) AddGift(ctx *gin.Context) {

	name := ctx.Query("name")
	link := ctx.Query("link")
	description := ctx.Query("description")

	gift, err := c.app.AddGift(name, link, description)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Bad Request")
		return
	}
	ctx.JSON(http.StatusOK, gift)
}
