package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Service() error {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "PONG")
	})

	return r.Run(":8080")
}
