package db

import (
	"ai-agent/entity"
	"context"
	"fmt"
	"time"

	"ai-agent/repositories/db/repo_dto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

type (
	BookRepository struct {
		agent *pgxpool.Pool
	}
)

const (
	bookTable            = "domain_book"
	bookCategoryTable    = "domain_book_category"
	bookAuthorTable      = "domain_book_author"
	bookAuthorMapTable   = "domain_book_author_map"
	bookCategoryMapTable = "domain_book_category_map"
)

func NewBookRepository(db *PostgresRepo) *BookRepository {
	return &BookRepository{db.agent}
}

func (b *BookRepository) InsertBookIfNotExists(book *entity.BookEntity) error {
	ctx := context.Background()

	tx, err := b.agent.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var authors []repo_dto.BookAuthorDTO
	var categories []repo_dto.BookCategoryDTO

	for _, author := range book.Authors {
		authors = append(authors, repo_dto.BookAuthorDTO{
			ID:   getConsistentUUID(author.Name),
			Name: author.Name,
		})
	}

	for _, category := range book.Categories {
		categories = append(categories, repo_dto.BookCategoryDTO{
			ID:   getConsistentUUID(category.Name),
			Name: category.Name,
		})
	}

	if err := b.insertAuthors(ctx, tx, authors); err != nil {
		return err
	}

	if err := b.insertCategory(ctx, tx, categories); err != nil {
		return err
	}

	bookID := getConsistentUUID(book.Title)
	embedding := pgvector.NewVector(book.Embedding)

	insertQuery := fmt.Sprintf(`
		INSERT INTO %s (id, title, subtitle, description, thumbnail, published_year, rating_count, average_rating, num_pages, embedding, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, bookTable)

	_, err = tx.Exec(ctx, insertQuery, bookID, book.Title, book.Subtitle, book.Description, book.Thumb, book.Year, book.RatingCount, book.AverageRating, book.PageCount, embedding, time.Now())
	if err != nil {
		return fmt.Errorf("insert book: %w", err)
	}

	if err := b.insertBookAuthorMap(ctx, tx, bookID, authors); err != nil {
		return err
	}

	if err := b.insertBookCategoryMap(ctx, tx, bookID, categories); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (b *BookRepository) insertAuthors(ctx context.Context, tx pgx.Tx, authors []repo_dto.BookAuthorDTO) error {
	batch := &pgx.Batch{}

	query := fmt.Sprintf(`
        INSERT INTO %s (id, name) 
        VALUES ($1, $2)
        ON CONFLICT (id) DO NOTHING`, bookAuthorTable)

	for _, author := range authors {
		batch.Queue(query, author.ID, author.Name)
	}

	br := tx.SendBatch(ctx, batch)

	err := br.Close()
	if err != nil {
		return fmt.Errorf("insert authors: %w", err)
	}

	return nil
}
func (b *BookRepository) insertCategory(ctx context.Context, tx pgx.Tx, categories []repo_dto.BookCategoryDTO) error {
	insertCategoryQuery := fmt.Sprintf(`
		INSERT INTO %s (id, name) 
		VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING`, bookCategoryTable)

	batch := &pgx.Batch{}

	for _, category := range categories {
		batch.Queue(insertCategoryQuery, category.ID, category.Name)
	}

	br := tx.SendBatch(ctx, batch)

	err := br.Close()
	if err != nil {
		return fmt.Errorf("insert category: %w", err)
	}

	return nil
}

func (b *BookRepository) insertBookAuthorMap(ctx context.Context, tx pgx.Tx, bookID *pgtype.UUID, authors []repo_dto.BookAuthorDTO) error {
	insertBookAuthorMapQuery := fmt.Sprintf(`
		INSERT INTO %s (book_id, author_id) 
		VALUES ($1, $2)
		ON CONFLICT (book_id, author_id) DO NOTHING`, bookAuthorMapTable)

	batch := &pgx.Batch{}

	for _, author := range authors {
		batch.Queue(insertBookAuthorMapQuery, bookID, author.ID)
	}

	br := tx.SendBatch(ctx, batch)

	err := br.Close()
	if err != nil {
		return fmt.Errorf("insert book author map: %w", err)
	}

	return nil
}

func (b *BookRepository) insertBookCategoryMap(ctx context.Context, tx pgx.Tx, bookID *pgtype.UUID, categories []repo_dto.BookCategoryDTO) error {
	insertBookCategoryMapQuery := fmt.Sprintf(`
		INSERT INTO %s (book_id, category_id) 
		VALUES ($1, $2)
		ON CONFLICT (book_id, category_id) DO NOTHING`, bookCategoryMapTable)

	batch := &pgx.Batch{}

	for _, category := range categories {
		batch.Queue(insertBookCategoryMapQuery, bookID, category.ID)
	}

	br := tx.SendBatch(ctx, batch)

	err := br.Close()
	if err != nil {
		return fmt.Errorf("insert book category map: %w", err)
	}

	return nil
}

func getConsistentUUID(input string) *pgtype.UUID {
	namespace := uuid.NameSpaceDNS
	generated := uuid.NewSHA1(namespace, []byte(input))

	return &pgtype.UUID{
		Bytes: generated,
		Valid: true,
	}
}
