package agentcontroller

import (
	"ai-agent/model/servicemodels"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	ai servicemodels.AgentService
}

func New(ai servicemodels.AgentService) *Controller {
	return &Controller{ai}
}

func (ctrl *Controller) SendMessage(c *gin.Context) {
	var rq servicemodels.RequestBody
	if err := c.BindJSON(&rq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "code": "1"})
		return
	}

	if rq.UserID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	resp, err := ctrl.ai.SendMessageWithHistory(rq.UserID, rq.Message)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
