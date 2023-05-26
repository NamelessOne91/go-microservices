package data

import (
	"context"
	"testing"
)

const (
	numTestUsers = 3
)

func Test_insert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Postgres integration test")
	}

	dbContainer, connPool := SetupTestDatabase()
	defer dbContainer.Terminate(context.Background())

	postgresRepo := NewPostgresRepository(connPool)

	newUser := User{
		FirstName: "new",
		LastName:  "user",
		Email:     "new-user@gmail.com",
		Active:    1,
	}

	id, err := postgresRepo.Insert(newUser)
	if err != nil {
		t.Errorf("failed to create new user with error: %v", err)
	}
	if id != numTestUsers+1 {
		t.Errorf("expected new user ID to be %d, got %d", numTestUsers+1, id)
	}
}
