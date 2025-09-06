package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// ConnectAndInit connects to PostgreSQL and ensures tables exist
func ConnectAndInit() *sql.DB {
	// Load connection string from .env (DB_DSN)
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("❌ DB_DSN not found in environment variables")
	}

	// Open connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// Check connection
	if err := db.Ping(); err != nil {
		log.Fatal("❌ Database not reachable:", err)
	}

	// Create tables if not exist
	createTables := `
	CREATE TABLE IF NOT EXISTS car_data (
		id SERIAL PRIMARY KEY,
		car_id VARCHAR(10),
		speed INT,
		rpm INT,
		latitude DOUBLE PRECISION,
		longitude DOUBLE PRECISION,
		location VARCHAR(50),
		updated_at TIMESTAMP DEFAULT NOW()
	);
	`
	_, err = db.Exec(createTables)
	if err != nil {
		log.Fatal("❌ Failed to create tables:", err)
	}

	fmt.Println("✅ Database connected and tables ready!")
	return db
}
