package database

import (
	"context"
	"log"
	"testing"

	"github.com/maeshinshin/go-multiapi/internal/testutil"
)

func TestMain(m *testing.M) {

	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbname = dbName
	password = dbPwd
	username = dbUser

	teardown, containerData, err := testutil.MustStartMySQLContainer(dbname, password, dbUser)

	if err != nil {
		log.Fatalf("could not start mysql container: %v", err)
	}

	if containerData == nil {
		log.Fatalf("could not get mysql container Data: %v", err)
	}

	host = containerData.Host()
	port = containerData.Port()

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown mysql container: %v", err)
	}
}

func TestNew(t *testing.T) {
	srv := New()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := New()

	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestClose(t *testing.T) {
	srv := New()

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}
