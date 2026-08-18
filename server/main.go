package main

import (
	"context"
	"log/slog"
	"ragbox/config"
	_ "ragbox/logging"
	"ragbox/restful"
	"ragbox/store"
	"slices"
	"sync"

	"github.com/milvus-io/milvus/client/v2/entity"
)

const (
	milvusAddr   string = "172.17.0.2:19530"
	database     string = "knowledgebase"
	collection   string = "knowledgebase"
	dropIfExists bool   = false
)

func hasAutoIDPrimaryKeyID(schema *entity.Schema) bool {
	if schema == nil {
		return false
	}
	for _, f := range schema.Fields {
		if f.Name == "id" {
			return f.PrimaryKey && f.AutoID
		}
	}
	return false
}

func prepareVectorStore(ctx context.Context, s *store.MilvusStore) error {
	dbLst, err := s.ListDatabase(ctx)
	if err != nil {
		return err
	}
	if !slices.Contains(dbLst, database) {
		if err := s.CreateDatabase(ctx, database); err != nil {
			return err
		}
	}

	if err := s.UseDatabase(ctx, database); err != nil {
		return err
	}

	collections, err := s.ListCollection(ctx)
	if err != nil {
		return err
	}
	slog.Debug("collections", "collections", collections)

	exists := slices.Contains(collections, collection)

	if exists && dropIfExists {
		err = s.DropCollection(ctx, collection)
		if err != nil {
			return err
		}
		exists = false
	}

	if !exists {
		schema := entity.NewSchema().
			WithName("kb_schema").
			WithField(entity.NewField().
				WithName("id").
				WithDataType(entity.FieldTypeInt64).
				WithIsPrimaryKey(true).
				WithIsAutoID(true)).
			WithField(entity.NewField().
				WithName("document_id").
				WithDataType(entity.FieldTypeVarChar).
				WithMaxLength(36).
				WithIsPrimaryKey(false)).
			WithField(entity.NewField().
				WithName("content").
				WithDataType(entity.FieldTypeVarChar).
				WithMaxLength(65535).
				WithIsPrimaryKey(false)).
			WithField(entity.NewField().
				WithName("embedding").
				WithDataType(entity.FieldTypeFloatVector).
				WithDim(1536))
		err = s.CreateCollection(ctx, collection, schema)
		if err != nil {
			return err
		}
	}

	if err := s.CreateIndex(ctx, collection, "embedding"); err != nil {
		return err
	}

	return nil
}

func main() {

	if err := config.ReadCfgFromFile("/app/db_monitor/knowledgebase/config.yaml"); err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	dataStore, err := store.NewDataStore()
	if err != nil {
		panic(err)
	}
	defer dataStore.Close()

	vectorStore, err := store.NewMilvusStore(milvusAddr, database)
	if err != nil {
		panic(err)
	}
	defer vectorStore.Close(context.Background())

	if err := prepareVectorStore(context.Background(), vectorStore); err != nil {
		panic(err)
	}

	go func() {
		defer wg.Done()
		restful.NewRestfulServer(":8080").Run()
	}()

	wg.Wait()
}
