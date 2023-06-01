package data

import (
	"testing"
)

const (
	numTestLogs = 3
	testID      = "61799e857fb7ed9569be5ccf"
)

func Test_all(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Mongo integration test")
	}

	entries, err := testConfig.repository.All()
	if err != nil {
		t.Errorf("failed to retrieve all entries with error: %v", err)
	}
	if len(entries) != numTestLogs {
		t.Errorf("expected %d entries, got %d", numTestLogs, len(entries))
	}
}

func Test_get_one(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Mongo integration test")
	}

	tests := []struct {
		name     string
		id       string
		existing bool
	}{
		{"existing entry ID", testID, true},
		{"ID matches no entry", "000000000000000000000255", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := testConfig.repository.GetOne(tt.id)
			if (l == nil || err != nil) && tt.existing {
				t.Errorf("got an error while trying to retrieve entry with ID %s: %v", tt.id, err)
			}
		})
	}
}

func Test_insert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Mongo integration test")
	}

	l := LogEntry{
		Name: "insert_test",
		Data: "insertion integration test",
	}

	err := testConfig.repository.Insert(l)
	if err != nil {
		t.Errorf("failed while trying to insert: %v", err)
	}

	entries, _ := testConfig.repository.All()
	if len(entries) <= numTestLogs {
		t.Errorf("expected to have %d after insert, got %d", numTestLogs+1, len(entries))
	}
}

func Test_update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Mongo integration test")
	}

	tests := []struct {
		name     string
		entry    LogEntry
		existing bool
	}{
		{
			"update existing entry",
			LogEntry{
				ID:   testID,
				Name: "Log Entry 1",
				Data: "update integration test",
			},
			true,
		},
		{
			"entry does not exists",
			LogEntry{
				ID:   "000000000000000000000255",
				Name: "Log Entry 1",
				Data: "update integration test",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := testConfig.repository.Update(tt.entry)
			if tt.existing {
				if !ok || err != nil {
					t.Errorf("update of existing entry with ID %s failed: %v", tt.entry.ID, err)
				}
				l, _ := testConfig.repository.GetOne(tt.entry.ID)
				if l.Data != tt.entry.Data {
					t.Errorf("expected entry Data to be %s after update, got %s", tt.entry.Data, l.Data)
				}
			}
			if !tt.existing && err == nil {
				t.Errorf("expected update for entry with ID %s (not existing) to fail", tt.entry.ID)
			}
		})
	}
}

func Test_drop_collection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Mongo integration test")
	}

	err := testConfig.repository.DropCollection()
	if err != nil {
		t.Errorf("failed while trying to drop logs collection: %v", err)
	}

	entries, _ := testConfig.repository.All()
	if len(entries) > 0 {
		t.Errorf("expected no entries after drop collection, but found %d", len(entries))
	}
}
