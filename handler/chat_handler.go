package handler

import (
	"ai-agent/handler/handler_dto"
	"ai-agent/interface"
	"log"
	"net/http"

	"ai-agent/entity/errormodels"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	ai       _interface.AgentUsecase
	errorMgr errormodels.Errors
}

func NewChatHandler(ai _interface.AgentUsecase, e errormodels.Errors) *ChatHandler {
	return &ChatHandler{
		ai:       ai,
		errorMgr: e,
	}
}

func (ctrl *ChatHandler) ChatWithHistory(c *gin.Context) {
	var rq handler_dto.ChatRequest
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

	c.JSON(http.StatusCreated, handler_dto.ChatResponse{
		Message: resp,
	})
}

func (ctrl *ChatHandler) handleError(c *gin.Context, code errormodels.ErrorCodes, debugMsg string) {
	log.Printf("Error [%s]: %s", code, debugMsg)

	formatted := ctrl.errorMgr.GetFormattedError(code)

	c.AbortWithStatusJSON(formatted.Code, gin.H{
		"message": formatted.Message,
	})
}
