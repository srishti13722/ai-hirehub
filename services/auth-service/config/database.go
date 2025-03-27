package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// Global DB
var DB *pgx.Conn

func ConnectDataBase() {
	_ = godotenv.Load()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	log.Println("Connected to PostgreSQL!")
	DB = conn
}

func RunMigrations() {
	query := `
	CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firstname VARCHAR(20) NOT NULL,
    lastname VARCHAR(20) NOT NULL,
    email VARCHAR(50) NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) CHECK (role in('job_seeker', 'recruiter', 'admin')) NOT NULL DEFAULT 'job_seeker',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	_, err := DB.Exec(context.Background(), query)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Users table checked/created.")
}
