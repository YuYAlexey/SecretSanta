package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/adYushinW/SecretSanta/internal/app"
	"github.com/adYushinW/SecretSanta/internal/model"
	"github.com/adYushinW/SecretSanta/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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
	var password string

	user := model.Users{}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	switch {
	case strings.TrimSpace(user.Login) == "":
		{
			ctx.JSON(http.StatusBadRequest, "Login cannot be empty")
			return
		}
	case strings.TrimSpace(user.Password) == "":
		{
			ctx.JSON(http.StatusBadRequest, "Password cannot be empty")
			return
		}
	case strings.TrimSpace(user.FirstName) == "":
		{
			ctx.JSON(http.StatusBadRequest, "FirstName cannot be empty")
			return
		}
	case strings.TrimSpace(user.LastName) == "":
		{
			ctx.JSON(http.StatusBadRequest, "LastName cannot be empty")
			return
		}
	}

	if user.Password != "" {
		password, err = utils.GenerateHashPassword(user.Password)
		if err != nil {
			ctx.JSON(http.StatusConflict, "Password cannot be hashed")
			return
		}
	}

	if user.Age != "" {
		age, err = strconv.ParseUint(user.Age, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "Age must be a number")
			return
		}
	}

	newuser, err := c.app.AddUser(user.Login, password, user.FirstName, user.LastName, user.Sex, age)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	ctx.JSON(http.StatusOK, newuser)
}

func (c *Controller) Login(ctx *gin.Context) {

	user := model.Users{}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	login, err := c.app.Login(user.Login, user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	if login {
		ctx.SetCookie(user.Login, "yes", 60, "/", "", false, true)
	}

	ctx.JSON(http.StatusOK, fmt.Sprint("Login success!", user.Login))
}

func (c *Controller) CheckCookie(ctx *gin.Context) {

	user := model.Users{}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	_, err := ctx.Cookie(user.Login)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprint("No cookie with current login: ", user.Login))
		return
	}
	ctx.JSON(http.StatusOK, fmt.Sprint("Cookie Get Success for: ", user.Login))
}

func (c *Controller) Logout(ctx *gin.Context) {

	user := model.Users{}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	_, err := ctx.Cookie("user")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprint("User is already logout: ", user.Login))
		return
	}

	ctx.SetCookie(user.Login, "", -1, "/", "", false, true)
	ctx.JSON(http.StatusOK, fmt.Sprint("You are logout!", user.Login))
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

func (c *Controller) Cookie(ctx *gin.Context) {

	session := sessions.Default(ctx)

	if session.Get("hello") != "world" {
		session.Set("hello", "world")
		session.Save()
	}

	ctx.JSON(http.StatusOK, session.Get("hello"))
}
