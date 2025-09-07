package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// ConnectAndInit connects to PostgreSQL and ensures tables exist
func ConnectAndInit() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("❌ One or more DB environment variables are missing")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

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

	CREATE TABLE IF NOT EXISTS parking_sessions (
		session_id VARCHAR(20) PRIMARY KEY,
		car_id VARCHAR(10),
		started_at TIMESTAMP,
		stopped_at TIMESTAMP,
		duration_minutes INT
	);

	CREATE TABLE IF NOT EXISTS sos_events (
		id SERIAL PRIMARY KEY,
		car_id VARCHAR(10),
		latitude DOUBLE PRECISION,
		longitude DOUBLE PRECISION,
		timestamp TIMESTAMP DEFAULT NOW()
	);
	`
	_, err = db.Exec(createTables)
	if err != nil {
		log.Fatal("❌ Failed to create tables:", err)
	}

	fmt.Println("✅ Database connected and tables ready!")
	return db
}