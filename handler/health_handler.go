package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (*HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusCreated, "OK")
}
