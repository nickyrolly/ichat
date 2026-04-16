// File: cmd/migrate/main.go
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/migrate/main.go [up|down|force]")
	}

	command := os.Args[1]

	m, err := migrate.New(
		"file://migrations",
		"postgres://"+getEnv("DB_USER", "nickyrolly")+":"+getEnv("DB_PASSWORD", "")+"@"+getEnv("DB_HOST", "localhost")+":"+getEnv("DB_PORT", "5432")+"/"+getEnv("DB_NAME", "ichat")+"?sslmode=disable",
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migration up completed successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migration down completed successfully")
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Usage: go run cmd/migrate/main.go force <version>")
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("Migration force failed: %v", err)
		}
		log.Printf("Migration forced to version %d", version)
	default:
		log.Fatal("Invalid command. Use [up|down|force]")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
