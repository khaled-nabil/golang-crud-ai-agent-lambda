package main

import (
	"fmt"
	"log"
	"os"

	"ai-agent/adapters/secrets"
	"ai-agent/entity"
	"ai-agent/usecase"

	"ai-agent/repositories/db"
	"encoding/csv"
	"io"

	"strings"

	"strconv"

	"ai-agent/adapters/ollama"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	log.Println("Loading configuration...")
	appConfig, err := secrets.NewSecretsManager()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	migrationsPath := os.Getenv("DB_MIGRATIONS_PATH")
	if migrationsPath == "" {
		log.Fatalf("DB_MIGRATIONS_PATH is not set")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		appConfig.DBUser,
		appConfig.DBPassword,
		appConfig.DBHost,
		appConfig.DBPort,
		appConfig.DBName,
	)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dsn,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	err = m.Up()
	if err != nil {
		if err.Error() != "no change" {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		log.Print("No changes to migrate")
	}

	log.Println("Migration successful!")

	if sErr, err := m.Close(); err != nil {
		log.Printf("Error closing migration instance: %v", err)
	} else if sErr != nil {
		log.Printf("Migration close error: %v", sErr)
	}

	insertionDomain := os.Getenv("DB_INSERTION_DOMAIN")
	filepath := os.Getenv("DB_INSERTION_FILEPATH")

	if insertionDomain != "" && filepath == "" {
		log.Fatalf("DB_INSERTION_FILEPATH is not set")
	}

	switch insertionDomain {
	case "book":
		createBooks(appConfig, filepath)
	}
}

func createBooks(cfg *secrets.AppConfig, filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Reading file: %v", err)
	}
	defer f.Close()

	repository, err := db.NewPostgresRepo(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}
	bookRepo := db.NewBookRepository(repository)
	ollamaClient, errr := ollama.NewOllama()
	if errr != nil {
		log.Fatalf("Failed to initialize Ollama: %v", errr)
	}

	bookUsecase := usecase.NewBookUsecase(ollamaClient, bookRepo)

	csvr := csv.NewReader(f)

	_, err = csvr.Read()
	if err != nil {
		log.Fatalf("Failed to read CSV header: %v", err)
	}

	for {
		var authors []entity.BookAuthorEntity
		var categories []entity.BookCategoryEntity

		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			log.Fatalf("Failed to Read CSV %s %v", filepath, err)
		}

		y, err := strconv.Atoi(getValueOrDefault(row[8], "0"))
		if err != nil {
			log.Fatalf("Failed to convert year to int: %v", err)
		}
		ar, err := strconv.ParseFloat(getValueOrDefault(row[9], "0"), 32)
		if err != nil {
			log.Fatalf("Failed to convert average rating to float: %v", err)
		}
		pc, err := strconv.Atoi(getValueOrDefault(row[10], "0"))
		if err != nil {
			log.Fatalf("Failed to convert page count to int: %v", err)
		}
		rc, err := strconv.Atoi(getValueOrDefault(row[11], "0"))
		if err != nil {
			log.Fatalf("Failed to convert rating count to int: %v", err)
		}

		a := strings.SplitSeq(row[4], ";")
		for author := range a {
			authors = append(authors, entity.BookAuthorEntity{
				Name: author,
			})
		}

		c := strings.SplitSeq(row[5], ";")
		for category := range c {
			categories = append(categories, entity.BookCategoryEntity{
				Name: category,
			})
		}

		book := entity.BookEntity{
			Title:         row[2],
			Subtitle:      row[3],
			Thumb:         row[6],
			Description:   row[7],
			Year:          int16(y),
			AverageRating: float32(ar),
			PageCount:     int(pc),
			RatingCount:   int(rc),
			Authors:       authors,
			Categories:    categories,
		}

		if err := bookUsecase.Insert(&book); err != nil {
			log.Printf("Failed to insert book: %s %d || %v", book.Title, book.Year, err)
		}
	}
}

func getValueOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}