package agentcontroller

import (
	"ai-agent/model/agentmodel"
	"ai-agent/model/geminimodel"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	ai geminimodel.Gemini
}

var (
	UserIDKey = "user_id"
)

func New(ai geminimodel.Gemini) *Controller {
	return &Controller{ai}
}

func (ctrl *Controller) SendMessage(c *gin.Context) {
	var rq agentmodel.RequestBody
	if err := c.BindJSON(&rq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "code": "1"})
		return
	}

	resp, _, err := ctrl.ai.Chat(rq.Message, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
