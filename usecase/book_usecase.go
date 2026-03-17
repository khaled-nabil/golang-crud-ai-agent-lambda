package usecase

import (
	"ai-agent/entity"
	_interface "ai-agent/interface"
	"fmt"
)

type (
	BookUsecase struct {
		em _interface.EmbeddingModel
		cr _interface.BookRepo
	}
)

func NewBookUsecase(em _interface.EmbeddingModel, cr _interface.BookRepo) *BookUsecase {
	return &BookUsecase{em, cr}
}

func (b *BookUsecase) Insert(book *entity.BookEntity) error {
	embedding, err := b.em.EmbedSearchDocument(fmt.Sprintf("title: %s.\nsubtitle: %s.\ndescription: %s.\nby: %s.\ncategory: %s.", book.Title, book.Subtitle, book.Description, book.GetAuthorNames(), book.GetCategoryNames()))
	if err != nil {
		return err
	}

	book.Embedding = embedding

	return b.cr.InsertBookIfNotExists(book)
}
