package agentcontroller

import (
	"ai-agent/model/dtomodels"
	"ai-agent/model/errormodels"
	"ai-agent/model/servicemodels"
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	ai       servicemodels.AgentService
	errorMgr errormodels.Errors
}

func New(ai servicemodels.AgentService, e errormodels.Errors) *Controller {
	return &Controller{
		ai:       ai,
		errorMgr: e,
	}
}

func (ctrl *Controller) SendMessage(c *gin.Context) {
	var rq servicemodels.RequestBody
	if err := c.BindJSON(&rq); err != nil {
		ctrl.handleError(c, errormodels.ErrBadRequest, err.Error())
		return
	}

	if rq.UserID == "" {
		ctrl.handleError(c, errormodels.ErrUnauthorized, "user not authenticated")
		return
	}

	resp, err := ctrl.ai.SendMessageWithHistory(rq.UserID, rq.Message)
	if err != nil {
		ctrl.handleError(c, errormodels.ErrGeneric, err.Error())
		return
	}

	c.JSON(http.StatusCreated, dtomodels.GenAIResponse{
		Message: resp,
	})
}

func (ctrl *Controller) handleError(c *gin.Context, code errormodels.ErrorCodes, debugMsg string) {
	log.Printf("Error [%s]: %s", code, debugMsg)

	formatted := ctrl.errorMgr.GetFormattedError(code)

	c.AbortWithStatusJSON(formatted.Code, gin.H{
		"message": formatted.Message,
	})
}
