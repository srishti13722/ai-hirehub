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

func ConnectDataBase(){
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
	CREATE TABLE IF NOT EXISTS job_applications (
    application_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL,
    jobseeker_id UUID NOT NULL,
    recruiter_id UUID NOT NULL,
    status VARCHAR(50) DEFAULT 'applied', -- applied, shortlisted, rejected, offered
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	_, err := DB.Exec(context.Background(), query)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Job applications table checked/created.")
}
