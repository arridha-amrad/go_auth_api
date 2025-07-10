package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Driver PostgreSQL
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPostgresConnection(t *testing.T) {
	ctx := context.Background()

	// 1. Setup TestContainer PostgreSQL
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	defer postgresContainer.Terminate(ctx)

	// 2. Dapatkan koneksi database
	host, err := postgresContainer.Host(ctx)
	assert.NoError(t, err)

	port, err := postgresContainer.MappedPort(ctx, "5432")
	assert.NoError(t, err)

	connStr := fmt.Sprintf("host=%s port=%d user=postgres password=postgres dbname=testdb sslmode=disable", host, port.Int())
	db, err := sql.Open("postgres", connStr)
	assert.NoError(t, err)
	defer db.Close()

	// 3. Test koneksi dengan Ping
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = db.PingContext(ctxTimeout)
	assert.NoError(t, err)

	// 4. Test eksekusi query sederhana
	var result int
	err = db.QueryRowContext(ctx, "SELECT 1 + 1").Scan(&result)
	assert.NoError(t, err)
	assert.Equal(t, 2, result)

	t.Log("✅ Koneksi PostgreSQL berhasil di-test!")
}
