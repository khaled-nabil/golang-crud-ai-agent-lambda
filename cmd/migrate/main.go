package main

import (
	"fmt"
	"log"
	"os"

	"ai-agent/pkg/secretspkg"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	log.Println("Loading configuration...")
	appConfig, err := secretspkg.New()
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

	log.Printf("Connecting to database at %s:%d/%s...", appConfig.DBHost, appConfig.DBPort, appConfig.DBName)
	log.Printf("Using migrations from: %s", migrationsPath)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dsn,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	err = m.Up()
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Migration successful!")

	if sErr, err := m.Close(); err != nil {
		log.Printf("Error closing migration instance: %v", err)
	} else if sErr != nil {
		log.Printf("Migration close error: %v", sErr)
	}
}
