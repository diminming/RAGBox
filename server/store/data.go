package store

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUser     = "kbadmin"
	dbPassword = "123qweASD"
	dbHost     = "172.17.0.3:3306"
)

var (
	Store *DataStore
)

type DataStore struct {
	store *sql.DB
}

func (s *DataStore) Close() error {
	return s.store.Close()
}

func (s *DataStore) GetStore() *sql.DB {
	return s.store
}

func (s *DataStore) Insert(ctx context.Context, query string, args ...any) (sql.Result, error) {
	stmt, err := s.store.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *DataStore) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	stmt, err := s.store.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func NewDataStore() (*DataStore, error) {
	store, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+")/knowledgebase?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}
	if err = store.Ping(); err != nil {
		return nil, err
	}

	Store = &DataStore{
		store: store,
	}
	return Store, nil
}
