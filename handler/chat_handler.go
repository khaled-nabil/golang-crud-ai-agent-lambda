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
	du       _interface.DomainUsecase
	errorMgr errormodels.Errors
}

const (
	chatDomain = "chat"
)

func NewChatHandler(ai _interface.AgentUsecase, du _interface.DomainUsecase, e errormodels.Errors) *ChatHandler {
	return &ChatHandler{
		ai:       ai,
		du:       du,
		errorMgr: e,
	}
}

func (c *ChatHandler) ChatWithHistory(ctx *gin.Context) {
	var rq handler_dto.ChatRequest
	if err := ctx.BindJSON(&rq); err != nil {
		c.handleError(ctx, errormodels.ErrBadRequest, err.Error())
		return
	}

	if rq.UserID == "" {
		c.handleError(ctx, errormodels.ErrUnauthorized, "user not authenticated")
		return
	}

	si, err := c.du.GetInstructions(chatDomain)
	if err != nil {
		c.handleError(ctx, errormodels.ErrGeneric, err.Error())
		return
	}

	resp, err := c.ai.SendMessageWithHistory(rq.UserID, rq.Message, si)
	if err != nil {
		c.handleError(ctx, errormodels.ErrGeneric, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, handler_dto.ChatResponse{
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
