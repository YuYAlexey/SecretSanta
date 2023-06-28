package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/adYushinW/SecretSanta/internal/app"
	"github.com/adYushinW/SecretSanta/internal/model"
	"github.com/adYushinW/SecretSanta/internal/utils"
	"github.com/gin-gonic/gin"
)

var key *string
var session = NewSession[CookieLogin](secretKey)
var sLogin = "user"

const (
	//all time in seconds
	rememberMeExpTime     = 60 * 60 * 24 * 365
	standartCookieExpTime = 60 * 10
	cookieName            = "gin_cookie_auth_"
	secretKey             = "secret_key"
)

type CookieLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Remember string `json:"remember"`
}

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
	var remember bool
	var err error
	user := new(CookieLogin)

	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, err = ctx.Cookie(cookieName + user.Login)
	if err != http.ErrNoCookie {
		ctx.JSON(http.StatusBadRequest, "User is already logged")
		return
	}

	if user.Remember != "" {
		remember, err = strconv.ParseBool(user.Remember)
		if err != nil {
			return
		}
	}

	login, err := c.app.Login(user.Login, user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	if login && remember {
		session_key := session.Set(*user, rememberMeExpTime)
		key = &session_key
		ctx.SetCookie(cookieName+user.Login, "yes", rememberMeExpTime, "/", "", false, true)
	} else if login && !remember {
		session_key := session.Set(*user, standartCookieExpTime)
		key = &session_key
		ctx.SetCookie(cookieName+user.Login, "yes", standartCookieExpTime, "/", "", false, true)
	}

	ctx.JSON(http.StatusOK, "Login success!")
}

func (c *Controller) CheckCookie(ctx *gin.Context) {

	user := &CookieLogin{}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	_, err := ctx.Cookie(user.Login)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "No cookie for current login")
		return
	}
	ctx.JSON(http.StatusOK, fmt.Sprint("Cookie Get Succeed for: ", user.Login))
}

func (c *Controller) Logout(ctx *gin.Context) {

	user := CookieLogin{}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	_, err := ctx.Cookie(cookieName + user.Login)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "User is already logout")
		return
	}

	ctx.SetCookie(cookieName+user.Login, "", -1, "/", "", false, true)
	ctx.JSON(http.StatusOK, "You are logout!")
}

func (c *Controller) WatchGift(ctx *gin.Context) {

	gift, err := c.app.WatchGift()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Bad Request")
		return
	}
	ctx.JSON(http.StatusOK, gift)
}

func (c *Controller) InGame(ctx *gin.Context) {

	user := CookieLogin{}

	_, err := c.app.StartParticipate(user.Login, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Bad Request")
		return
	}
	ctx.JSON(http.StatusOK, "You are in game")
}

func (c *Controller) OutGame(ctx *gin.Context) {

	user := &CookieLogin{}

	_, err := c.app.StopParticipate(user.Login, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Bad Request")
		return
	}
	ctx.JSON(http.StatusOK, "You are not play anymore")
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

func (c *Controller) AuthRoute(ctx *gin.Context) {

	user := sLogin

	_, err := ctx.Cookie(cookieName + user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "user unauthorized")
		ctx.Abort()
		return
	}

	username, ok := session.Get(*key)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, "user unauthorized")
		ctx.Abort()
		return
	}
	ctx.Set("username", username)
	ctx.Next()
}
