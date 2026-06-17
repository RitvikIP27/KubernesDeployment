package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	var connStr string

	// Support DATABASE_URL directly (standard for production Supabase/hosted DBs)
	if envConn := os.Getenv("DATABASE_URL"); envConn != "" {
		connStr = envConn
	} else {
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "postgres")
		dbname := getEnv("DB_NAME", "postgres")
		sslmode := getEnv("DB_SSLMODE", "disable")

		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	}

	var err error
	for i := 0; i < 30; i++ {
		DB, err = sql.Open("postgres", connStr)
		if err == nil {
			err = DB.Ping()
			if err == nil {
				log.Println("Connected to PostgreSQL/Supabase database")
				DB.SetMaxOpenConns(10)
				DB.SetMaxIdleConns(5)
				DB.SetConnMaxLifetime(5 * time.Minute)
				return
			}
		}
		log.Printf("Waiting for PostgreSQL database... attempt %d/30 (error: %v)", i+1, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("Could not connect to database: %v", err)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
