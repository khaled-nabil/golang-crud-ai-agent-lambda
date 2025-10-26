package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
}

func New() *Controller {
	return &Controller{}
}

func (*Controller) Health(c *gin.Context) {
	log.Println("Health check OK")

	c.JSON(http.StatusCreated, "OK")
}
