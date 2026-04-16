// File: internal/database/connection.go
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// _ "github.com/lib/pq"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
)

type DB struct {
	*sql.DB
}

func NewConnection() (*DB, error) {
	// Setup Viper
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Set default values
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "nickyrolly")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "ichat")

	log.Printf("DB_HOST: %s", viper.GetString("DB_HOST"))
	log.Printf("DB_PORT: %s", viper.GetString("DB_PORT"))
	log.Printf("DB_USER: %s", viper.GetString("DB_USER"))
	log.Printf("DB_PASSWORD: %s", viper.GetString("DB_PASSWORD"))
	log.Printf("DB_NAME: %s", viper.GetString("DB_NAME"))

	// Read .env file
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning: .env file not found, using environment variables or defaults")
	}

	configFile := viper.ConfigFileUsed()
	log.Printf("Config file used: %s", configFile)

	allSettings := viper.AllSettings()
	log.Printf("All settings: %+v", allSettings)

	// Automatically read environment variables
	viper.AutomaticEnv()

	host := viper.GetString("DB_HOST")
	if host == "localhost" {
		host = "127.0.0.1" // Force IPv4
	}
	port := viper.GetString("DB_PORT")
	user := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	dbname := viper.GetString("DB_NAME")

	// Debug connection string
	// connStr := "host=127.0.0.1 port=5432 user=nickyrolly dbname=ichat sslmode=disable"

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	if password != "" {
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	}

	// connStr := "host=127.0.0.1 port=5432 user=nickyrolly dbname=ichat sslmode=disable"
	log.Printf("Connection string: %s", connStr)

	log.Printf("Connecting to database: %s as user %s", dbname, user)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return &DB{db}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
