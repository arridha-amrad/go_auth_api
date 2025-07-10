package repositories

import (
	"context"
	"database/sql"
	"testing"
)

func TestWithRealDB(t *testing.T) {
	ctx := context.Background()

	container, err := SetupTestDatabase(ctx) // jika di package repositories
	// atau testutils.SetupTestDatabase(ctx) jika di package terpisah
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)

	// Buat tabel
	if err := createTables(container.DB); err != nil {
		t.Fatal(err)
	}

	if _, err := container.DB.Exec("DELETE FROM users"); err != nil {
		t.Fatalf("failed to clean data: %v", err)
	}

	db := container.DB
	repo := NewUserRepository(db)

	t.Run("Create and get user", func(t *testing.T) {
		user, err := repo.CreateOne(ctx, CreateOneParams{
			Name:       "Test",
			Username:   "testuser",
			Email:      "test@example.com",
			Password:   "secure",
			JWTVersion: "1",
		})

		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		fetched, err := repo.GetById(ctx, user.ID)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if fetched.Username != "testuser" {
			t.Errorf("Expected testuser, got %s", fetched.Username)
		}
	})
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY,
            name TEXT NOT NULL,
            username TEXT NOT NULL UNIQUE,
            email TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL,
            jwt_version TEXT NOT NULL,
            provider TEXT NOT NULL DEFAULT 'local',
            is_verified BOOLEAN NOT NULL DEFAULT false,
            role TEXT NOT NULL DEFAULT 'user',
            created_at TIMESTAMP NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMP NOT NULL DEFAULT NOW()
        )
    `)
	return err
}
