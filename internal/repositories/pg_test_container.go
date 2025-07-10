package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestContainer struct {
	testcontainers.Container
	DB *sql.DB
}

func SetupTestDatabase(ctx context.Context) (*TestContainer, error) {
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

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=postgres password=postgres dbname=testdb sslmode=disable",
		host, port.Int()))
	if err != nil {
		return nil, err
	}

	return &TestContainer{Container: container, DB: db}, nil
}
