package database

import (
    "database/sql"
    "fmt"
    "os"

    _ "github.com/lib/pq"
)

var DB *sql.DB

// Connect функциясы дерекқорға қосылады және error қайтарады
func Connect() error {
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASS"),
        os.Getenv("DB_NAME"),
    )

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("❌ SQL Open error: %v", err)
    }

    if err = db.Ping(); err != nil {
        return fmt.Errorf("❌ Database ping error: %v", err)
    }

    DB = db
    return nil
}
