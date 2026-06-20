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
				ensureSchema()
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

func ensureSchema() {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT,
			oauth_provider TEXT,
			oauth_id TEXT,
			name TEXT,
			avatar_url TEXT,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS skills (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL,
			category VARCHAR(50) DEFAULT '',
			target_hours INT DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS learning_logs (
			id SERIAL PRIMARY KEY,
			skill_id INT NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			hours DECIMAL(4,1) NOT NULL,
			notes TEXT,
			log_date DATE NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS settings (
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			key VARCHAR(50) NOT NULL,
			value TEXT NOT NULL,
			PRIMARY KEY (user_id, key)
		)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS oauth_provider TEXT`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS oauth_id TEXT`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS name TEXT`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url TEXT`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP`,
		`INSERT INTO users (email, password_hash, name)
		 VALUES ('admin@example.com', '$2b$12$XAxdZVaBMFdWah3WnS92l.hSsl9Nion996q8YX4eAx2ogiUOfH1be', 'Admin User')
		 ON CONFLICT (email) DO NOTHING`,
		`ALTER TABLE skills ADD COLUMN IF NOT EXISTS user_id INT`,
		`UPDATE skills
		 SET user_id = (SELECT id FROM users WHERE email = 'admin@example.com')
		 WHERE user_id IS NULL`,
		`ALTER TABLE skills ALTER COLUMN user_id SET NOT NULL`,
		`ALTER TABLE learning_logs ADD COLUMN IF NOT EXISTS user_id INT`,
		`UPDATE learning_logs
		 SET user_id = (SELECT id FROM users WHERE email = 'admin@example.com')
		 WHERE user_id IS NULL`,
		`ALTER TABLE learning_logs ALTER COLUMN user_id SET NOT NULL`,
		`ALTER TABLE settings ADD COLUMN IF NOT EXISTS user_id INT`,
		`UPDATE settings
		 SET user_id = (SELECT id FROM users WHERE email = 'admin@example.com')
		 WHERE user_id IS NULL`,
		`ALTER TABLE settings ALTER COLUMN user_id SET NOT NULL`,
		`CREATE UNIQUE INDEX IF NOT EXISTS settings_user_id_key_idx ON settings (user_id, key)`,
		`INSERT INTO settings (user_id, key, value)
		 SELECT id, 'target_role', 'Site Reliability Engineer (SRE)' FROM users WHERE email = 'admin@example.com'
		 ON CONFLICT (user_id, key) DO NOTHING`,
	}

	for _, statement := range statements {
		if _, err := DB.Exec(statement); err != nil {
			log.Fatalf("Could not initialize database schema: %v", err)
		}
	}

	log.Println("Database schema is ready")
}
