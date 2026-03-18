package handler

import (
	"ai-agent/handler/handler_dto"
	_interface "ai-agent/interface"
	"log"
	"net/http"

	"ai-agent/entity/errormodels"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bu  _interface.BookUsecase
	au _interface.AgentAdapter
	du _interface.DomainUsecase
	errorMgr errormodels.Errors
}

const (
	BookDomain = "book"
)

func NewBookHandler(bu _interface.BookUsecase, au _interface.AgentAdapter, du _interface.DomainUsecase, e errormodels.Errors) *BookHandler {
	return &BookHandler{
		bu:       bu,
		au:       au,
		du:       du,
		errorMgr: e,
	}
}

func (b *BookHandler) CreateBook(c *gin.Context) {
	var rq handler_dto.BookRequest
	if err := c.BindJSON(&rq); err != nil {
		b.handleError(c, errormodels.ErrBadRequest, err.Error())
		return
	}

	err := b.bu.Insert(rq.ToBookEntity())
	if err != nil {
		b.handleError(c, errormodels.ErrGeneric, err.Error())
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

func (b *BookHandler) GetBookRecommendations(c *gin.Context) {
	prompt := c.Query("prompt")
	if prompt == "" {
		b.handleError(c, errormodels.ErrBadRequest, "prompt is required", "Get Parameter \"prompt\" is required")
		return
	}

	books, err := b.bu.GetBookRecommendations(prompt)
	if err != nil {
		b.handleError(c, errormodels.ErrGeneric, err.Error())
		return
	}

	sysIns, err := b.du.GetInstructions(BookDomain)
	if err != nil {
		b.handleError(c, errormodels.ErrGeneric, err.Error())
		return
	}

	recommendations, err := b.au.RecommendBookFromList(prompt, sysIns, books)
	if err != nil {
		b.handleError(c, errormodels.ErrGeneric, err.Error())
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

func (ctrl *BookHandler) handleError(c *gin.Context, code errormodels.ErrorCodes, debugMsg string, userMessage ...string) {
	log.Printf("Error [%s]: %s", code, debugMsg)

	formatted := ctrl.errorMgr.GetFormattedError(code, userMessage...)

	c.AbortWithStatusJSON(formatted.Code, gin.H{
		"message": formatted.Message,
	})
}
