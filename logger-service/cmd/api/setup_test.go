package main

import (
	"os"
	"testing"

	"github.com/NamelessOne91/logger/data"
)

var testApp Config

func TestMain(m *testing.M) {
	repo := data.NewMongoTestRepository(nil)
	testApp.Repo = repo

	os.Exit(m.Run())
}
