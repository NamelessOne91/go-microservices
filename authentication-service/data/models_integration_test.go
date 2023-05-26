package data

import (
	"testing"
)

const (
	numTestUsers = 3
)

func Test_get_all(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Postgres integration test")
	}

	users, err := testConfig.repository.GetAll()
	if err != nil {
		t.Errorf("failed to retrieve all users with error: %v", err)
	}
	if len(users) != numTestUsers {
		t.Errorf("expected %d users, got %d", numTestUsers, len(users))
	}
}

func Test_insert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Postgres integration test")
	}

	newUser := User{
		FirstName: "new",
		LastName:  "user",
		Email:     "new-user@gmail.com",
		Active:    1,
	}

	id, err := testConfig.repository.Insert(newUser)
	if err != nil {
		t.Errorf("failed to create new user with error: %v", err)
	}
	if id != numTestUsers+1 {
		t.Errorf("expected new user ID to be %d, got %d", numTestUsers+1, id)
	}
}
