package agentcontroller

import (
	"ai-agent/model/aiagentmodel"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	ai aiagentmodel.AgentService
}

var (
	UserIDKey = "user_id"
)

func New(ai aiagentmodel.AgentService) *Controller {
	return &Controller{ai}
}

func (ctrl *Controller) SendMessage(c *gin.Context) {
	var rq aiagentmodel.RequestBody
	if err := c.BindJSON(&rq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "code": "1"})
		return
	}

	resp, err := ctrl.ai.SendMessageWithHistory(rq.Message)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
