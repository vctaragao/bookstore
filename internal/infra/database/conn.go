package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func NewConn() *sql.DB {
	dbHost := getEnv("DB_HOST")
	dbName := getEnv("DB_NAME")
	dbPort := getEnv("DB_PORT")
	dbUsername := getEnv("DB_USER")
	dbPassword := getEnv("DB_PASSWORD")

	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost,
		dbPort,
		dbUsername,
		dbPassword,
		dbName,
	)

	dbConn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("unable to open connection: %v\n", err)
	}

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("unable to ping database: %v\n", err)
	}

	return dbConn
}

func getEnv(key string) string {
	envValue := os.Getenv(key)
	if envValue == "" {
		log.Fatalf("missing required env %v", key)
	}

	return envValue
}
