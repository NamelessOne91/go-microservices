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

func Test_get_by_email(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Postgres integration test")
	}

	tests := []struct {
		name     string
		email    string
		existing bool
	}{
		{"existing email address", "test1@example.com", true},
		{"email address matches no user", "test255@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := testConfig.repository.GetByEmail(tt.email)
			if (user == nil) && tt.existing {
				t.Errorf("got an error while trying to retrieve user with email %s: %v", tt.email, err)
			}
		})
	}
}

func Test_get_one(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Postgres integration test")
	}

	tests := []struct {
		name     string
		id       int
		existing bool
	}{
		{"existing user ID", 3, true},
		{"ID matches no user", 255, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := testConfig.repository.GetOne(tt.id)
			if (user == nil) && tt.existing {
				t.Errorf("got an error while trying to retrieve user with ID %d: %v", tt.id, err)
			}
		})
	}
}

func Test_delete_by_ID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Postgres integration test")
	}
	userID := 1

	err := testConfig.repository.DeleteByID(1)
	if err != nil {
		t.Errorf("got an error while trying to delete user with ID %d: %v", userID, err)
	}

	user, _ := testConfig.repository.GetOne(userID)
	if user != nil {
		t.Errorf("user with ID %d should have been deleted but is still present in DB", userID)
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

func Test_update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Postgres integration test")
	}

	user := User{
		ID:        3,
		FirstName: "Updated",
		LastName:  "Johnson",
		Email:     "test3@example.com",
		Active:    1,
	}

	err := testConfig.repository.Update(user)
	if err != nil {
		t.Errorf("failed to update user with ID %d: %v", user.ID, err)
	}

	updated, _ := testConfig.repository.GetOne(user.ID)
	if updated.FirstName != user.FirstName {
		t.Errorf("expected user first name to have been updated to %s - got %s", user.FirstName, updated.FirstName)
	}
}

func Test_passwords_handling(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Postgres integration test")
	}

	tests := []struct {
		name        string
		userID      int
		newPswd     string
		pswdToCheck string
		match       bool
	}{
		{"password successfully changed and matched", 2, "new_test_password_1", "new_test_password_1", true},
		{"password does not match after change", 2, "new_test_password_2", "new_test_password_1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := User{ID: tt.userID}

			err := testConfig.repository.ResetPassword(tt.newPswd, user)
			if err != nil {
				t.Errorf("error while trying to update password for user with ID %d: %v", user.ID, err)
			}

			u, _ := testConfig.repository.GetOne(tt.userID)
			match, err := testConfig.repository.PasswordMatches(tt.pswdToCheck, *u)
			if err != nil {
				t.Errorf("got an error while checking user's password: %v", err)
			}
			if match != tt.match {
				t.Errorf("expected password check to be %t, got %t", tt.match, match)
			}
		})
	}
}
