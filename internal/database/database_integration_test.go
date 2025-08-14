package database

import (
	"os"
	"testing"
)

func TestDatabase_Connect(t *testing.T) {
	host := getenv("POSTGRES_HOST", "localhost")
	port := getenv("POSTGRES_PORT", "5432")
	user := getenv("POSTGRES_USER", "postgres")
	pass := getenv("POSTGRES_PASSWORD", "password")
	dbn := getenv("POSTGRES_DB", "wb_orders")

	db, err := NewDatabase(host, port, user, pass, dbn)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
