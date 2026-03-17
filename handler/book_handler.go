package handler

import (
	"ai-agent/handler/handler_dto"
	"ai-agent/interface"
	"log"
	"net/http"

	"ai-agent/entity/errormodels"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	usecase  _interface.BookUsecase
	errorMgr errormodels.Errors
}

func NewBookHandler(uc _interface.BookUsecase, e errormodels.Errors) *BookHandler {
	return &BookHandler{
		usecase:  uc,
		errorMgr: e,
	}
}

func (ctrl *BookHandler) CreateBook(c *gin.Context) {
	var rq handler_dto.BookRequest
	if err := c.BindJSON(&rq); err != nil {
		ctrl.handleError(c, errormodels.ErrBadRequest, err.Error())
		return
	}

	err := ctrl.usecase.Insert(rq.ToBookEntity())
	if err != nil {
		ctrl.handleError(c, errormodels.ErrGeneric, err.Error())
		return
	}

	c.JSON(http.StatusCreated, handler_dto.BookRequest{
		Title:         rq.Title,
		Subtitle:      rq.Subtitle,
		Authors:       rq.Authors,
		Categories:    rq.Categories,
		Description:   rq.Description,
		Year:          rq.Year,
		AverageRating: rq.AverageRating,
		RatingsCount:  rq.RatingsCount,
		PagesCount:    rq.PagesCount,
	})
}

func (ctrl *BookHandler) handleError(c *gin.Context, code errormodels.ErrorCodes, debugMsg string) {
	log.Printf("Error [%s]: %s", code, debugMsg)

	formatted := ctrl.errorMgr.GetFormattedError(code)

	c.AbortWithStatusJSON(formatted.Code, gin.H{
		"message": formatted.Message,
	})
}
