package store_test

import (
	"context"
	"ragbox/restful"
	"ragbox/store"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	source, err := store.NewDataStore()
	if err != nil {
		t.Fatalf("failed to create data store: %v", err)
	}
	sqlstring := "insert into users(username, passwd, create_timestamp, update_timestamp) values(?, ?, ?, ?)"
	hashed, err := restful.HashPassword("admin1!")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	result, err := source.Insert(context.Background(), sqlstring, "admin", hashed, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}
	if result == nil {
		t.Fatalf("expected result, got nil")
	}
}
