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
	CREATE TABLE IF NOT EXISTS jobs (
    job_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recruiter_id UUID NOT NULL, -- FK to recruiters table
    organisation_name VARCHAR(255),
    job_title VARCHAR(255) NOT NULL,
    job_description TEXT NOT NULL,
    job_location VARCHAR(100),
    salary INT,
    skills_required TEXT[],
    vacancy INT DEFAULT 1,
    status VARCHAR(50) DEFAULT 'active', -- active, closed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	_, err := DB.Exec(context.Background(), query)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Jobs table checked/created.")
}
