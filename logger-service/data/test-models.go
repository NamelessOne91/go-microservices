package data

import (
	"context"
	"errors"
	"log"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoTestRepository struct {
	Conn *mongo.Client
}

var ErrNoName = errors.New("empty or missing name field")
var ErrNoData = errors.New("empty or missing data field")

func NewMongoTestRepository(conn *mongo.Client) *MongoTestRepository {
	return &MongoTestRepository{
		Conn: conn,
	}
}

func (r *MongoTestRepository) Insert(entry LogEntry) error {
	if entry.Name == "" {
		return ErrNoName
	}
	if entry.Data == "" {
		return ErrNoData
	}
	return nil
}

func (r *MongoTestRepository) All() ([]*LogEntry, error) {
	logs := []*LogEntry{}
	return logs, nil
}

func (r *MongoTestRepository) GetOne(id string) (*LogEntry, error) {
	entry := LogEntry{
		ID:        "a1b2c3",
		Name:      "test",
		Data:      "test log",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return &entry, nil
}

func (r *MongoTestRepository) DropCollection() error {
	return nil
}

func (r *MongoTestRepository) Update(l LogEntry) (bool, error) {
	return true, nil
}

func SetupTestDatabase() (testcontainers.Container, *mongo.Client) {
	ctx := context.Background()
	initPath := filepath.Join("testdata", "init-logs-db.js")
	initAbsPath, err := filepath.Abs(initPath)
	if err != nil {
		log.Fatalf("could not determine abs path for Mongo test container init script")
	}

	req := testcontainers.ContainerRequest{
		Image:        "mongo:6",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Waiting for connections"),
			wait.ForListeningPort("27017/tcp"),
		),
		Env: map[string]string{
			"MONGO_INITDB_DATABASE":      "logs",
			"MONGO_INITDB_ROOT_USERNAME": "admin",
			"MONGO_INITDB_ROOT_PASSWORD": "password",
		},
		Mounts: []testcontainers.ContainerMount{
			{
				Source:   testcontainers.GenericBindMountSource{HostPath: initAbsPath},
				Target:   "/docker-entrypoint-initdb.d/init-logs-db.js",
				ReadOnly: false,
			},
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("failed to init Mongo test container: %v", err)
	}

	endpoint, err := container.Endpoint(ctx, "mongodb")
	if err != nil {
		log.Fatalf("could not retrieve Mongo test container endpoint: %v", err)
	}

	clientOpts := options.Client().ApplyURI(endpoint)
	clientOpts.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		log.Fatalf("failed to connect to Mongo test container: %v", err)
	}
	return container, c
}
