// Package database creates and manages the MySQL connection pool I use across
// the whole application.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	// I import the driver for its side effect: it registers the "mysql" driver
	// with database/sql so my sql.Open("mysql", ...) call can find it.
	"github.com/go-sql-driver/mysql"
)

// getEnv reads an environment variable and falls back to a default when it is
// not set. I use this so the same binary runs both locally and in Docker.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Connection opens a MySQL connection pool, checks that the database is
// reachable and returns it. I read the settings from environment variables so I
// can run the same binary locally and inside a container.
func Connection() (*sql.DB, error) {
	cfg := mysql.Config{
		// I read the credentials from the environment, falling back to my local
		// defaults.
		User:   getEnv("DB_USER", "root"),
		Passwd: getEnv("DB_PASSWORD", "Root@123"),
		Net:    "tcp",
		Addr:   getEnv("DB_ADDRESS", "127.0.0.1:3306"),
		DBName: getEnv("DB_NAME", "products"),
	}

	// sql.Open only validates the DSN; it does not open a real connection yet.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Here I tune the pool: recycle connections periodically and cap how many
	// are open and idle at the same time.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(15)

	// I ping so the program fails fast if the database is unreachable.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
