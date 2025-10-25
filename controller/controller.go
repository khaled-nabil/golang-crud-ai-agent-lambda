package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
}

func New() *Controller {
	return &Controller{}
}

func (*Controller) Health(c *gin.Context) {
	c.JSON(http.StatusCreated, "OK")
}
