package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresTestRepository struct {
	Conn *sql.DB
}

var ErrMissingEmail = errors.New("user email has not been provided")
var ErrMissingPassword = errors.New("user password has not been provided")

func NewPostgresTestRepository(db *sql.DB) *PostgresTestRepository {
	return &PostgresTestRepository{
		Conn: db,
	}
}

func (u *PostgresTestRepository) GetAll() ([]*User, error) {
	users := []*User{}
	return users, nil
}

func (u *PostgresTestRepository) GetByEmail(email string) (*User, error) {
	if email == "" {
		return nil, ErrMissingEmail
	}

	user := User{
		ID:        1,
		FirstName: "First",
		LastName:  "Last",
		Email:     "me@here.com",
		Password:  "",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return &user, nil
}

func (u *PostgresTestRepository) GetOne(id int) (*User, error) {
	user := User{
		ID:        1,
		FirstName: "First",
		LastName:  "Last",
		Email:     "me@here.com",
		Password:  "",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

func (u *PostgresTestRepository) Update(user User) error {
	return nil
}

func (u *PostgresTestRepository) DeleteByID(id int) error {
	return nil
}

func (u *PostgresTestRepository) Insert(user User) (int, error) {
	return 2, nil
}

func (u *PostgresTestRepository) ResetPassword(password string, user User) error {
	return nil
}

func (u *PostgresTestRepository) PasswordMatches(plainText string, user User) (bool, error) {
	if plainText == "" {
		return false, ErrMissingPassword
	}
	return true, nil
}

func SetupTestDatabase() (testcontainers.Container, *sql.DB) {
	// create PostgreSQL container request
	container, err := postgres.RunContainer(
		context.Background(),
		testcontainers.WithImage("docker.io/postgres:14.8"),
		postgres.WithInitScripts(filepath.Join("testdata", "init-user-db.sql")),
		postgres.WithDatabase("users"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(10*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to init Postgres test container: %v", err)
	}

	// get host and port of PostgreSQL container
	host, _ := container.Host(context.Background())
	port, _ := container.MappedPort(context.Background(), "5432")

	// create db connection string and connect
	dbURI := fmt.Sprintf("postgres://postgres:postgres@%v:%v/users", host, port.Port())
	connPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		log.Fatalf("failed to connect to Postgres test container: %v", err)
	}

	return container, connPool
}
