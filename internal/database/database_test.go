package database

import (
	"context"
	"log"
	"testing"

	"github.com/maeshinshin/go-multiapi/internal/util"
)

func TestMain(m *testing.M) {

	DbInfo.DB_DATABASE = "database"
	DbInfo.DB_USERNAME = "user"
	DbInfo.DB_PASSWORD = "password"

	teardown, err := util.MustStartMySQLContainer(DbInfo)

	if err != nil {
		log.Fatalf("could not start mysql container: %v", err)
	}

	if DbInfo.Db_HOST == "" || DbInfo.Db_PORT == "" {
		log.Fatalf("could not get mysql container Data: %v", DbInfo)
	}

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
