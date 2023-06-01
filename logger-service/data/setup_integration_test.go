package data

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type dbIntegrationTestConfig struct {
	repository  Repository
	dbContainer testcontainers.Container
	db          *mongo.Client
}

var testConfig dbIntegrationTestConfig

func TestMain(m *testing.M) {
	container, client := SetupTestDatabase()
	defer func() {
		fmt.Println("terminating Mongo test container")
		container.Terminate(context.Background())
	}()

	testConfig = dbIntegrationTestConfig{
		repository:  NewMongoRepository(client),
		dbContainer: container,
		db:          client,
	}

	code := m.Run()
	os.Exit(code)
}
