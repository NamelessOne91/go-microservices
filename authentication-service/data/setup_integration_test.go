package data

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

type dbIntegrationTestConfig struct {
	repository  Repository
	dbContainer testcontainers.Container
	db          *sql.DB
}

var testConfig dbIntegrationTestConfig

func TestMain(m *testing.M) {
	container, connPool := SetupTestDatabase()
	defer container.Terminate(context.Background())

	testConfig = dbIntegrationTestConfig{
		repository:  NewPostgresRepository(connPool),
		dbContainer: container,
		db:          connPool,
	}

	code := m.Run()
	os.Exit(code)
}
