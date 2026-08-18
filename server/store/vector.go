package store

import (
	"context"
	"errors"

	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

var (
	VectorStore *MilvusStore
)

func NewMilvusStore(addr string, database string) (*MilvusStore, error) {
	client, err := milvusclient.New(context.Background(), &milvusclient.ClientConfig{
		Address: addr,
	})

	if err != nil {
		return nil, err
	}

	store := &MilvusStore{
		client: client,
	}
	VectorStore = store
	return store, nil
}

type MilvusStore struct {
	client *milvusclient.Client
}

func (s *MilvusStore) ListDatabase(ctx context.Context) ([]string, error) {
	dbs, err := s.client.ListDatabase(ctx, milvusclient.NewListDatabaseOption())
	if err != nil {
		return nil, err
	}
	return dbs, nil
}

func (s *MilvusStore) CreateDatabase(ctx context.Context, database string) error {
	err := s.client.CreateDatabase(ctx, milvusclient.NewCreateDatabaseOption(database))
	if err != nil {
		return err
	}
	return nil
}

func (s *MilvusStore) DropDatabase(ctx context.Context, database string) error {
	err := s.client.DropDatabase(ctx, milvusclient.NewDropDatabaseOption(database))
	if err != nil {
		return err
	}
	return nil
}

func (s *MilvusStore) UseDatabase(ctx context.Context, database string) error {
	return s.client.UseDatabase(ctx, milvusclient.NewUseDatabaseOption(database))
}

func (s *MilvusStore) Close(ctx context.Context) error {
	return s.client.Close(ctx)
}

func (s *MilvusStore) ListCollection(ctx context.Context) ([]string, error) {
	collections, err := s.client.ListCollections(ctx, milvusclient.NewListCollectionOption())
	if err != nil {
		return nil, err
	}
	return collections, nil
}

func (s *MilvusStore) CreateCollection(ctx context.Context, name string, schema *entity.Schema) error {
	option := milvusclient.NewCreateCollectionOption(name, schema)
	return s.client.CreateCollection(ctx, option)
}

func (s *MilvusStore) DropCollection(ctx context.Context, name string) error {
	option := milvusclient.NewDropCollectionOption(name)
	return s.client.DropCollection(ctx, option)
}

func (s *MilvusStore) DescribeCollection(ctx context.Context, name string) (*entity.Collection, error) {
	option := milvusclient.NewDescribeCollectionOption(name)
	collection, err := s.client.DescribeCollection(ctx, option)
	if err != nil {
		return nil, err
	}
	return collection, nil
}

func (s *MilvusStore) CreateIndex(ctx context.Context, collectionName string, fieldName string) error {
	option := milvusclient.NewCreateIndexOption(collectionName, fieldName, index.NewHNSWIndex(entity.L2, 16, 200))
	_, err := s.client.CreateIndex(ctx, option)
	if err != nil {
		return err
	}
	return nil
}

func (s *MilvusStore) Insert(ctx context.Context, collectionName string, rows []map[string]any) error {
	rowItems := make([]any, 0, len(rows))
	for _, row := range rows {
		rowItems = append(rowItems, row)
	}

	option := milvusclient.NewRowBasedInsertOption(collectionName, rowItems...)
	_, err := s.client.Insert(ctx, option)
	if err != nil {
		return err
	}
	return nil
}

func (s *MilvusStore) LoadCollection(ctx context.Context, collectionName string) error {

	state, err := s.client.GetLoadState(ctx, milvusclient.NewGetLoadStateOption(collectionName))
	if err != nil {
		return err
	}
	if state.State == entity.LoadStateLoaded {
		return nil
	}

	task, err := s.client.LoadCollection(ctx, milvusclient.NewLoadCollectionOption(collectionName))
	if err != nil {
		return err
	}
	err = task.Await(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *MilvusStore) Search(ctx context.Context, collectionName string, vectors []float32, limit int) ([]milvusclient.ResultSet, error) {

	err := s.LoadCollection(ctx, collectionName)
	if err != nil {
		return nil, err
	}

	option := milvusclient.NewSearchOption(collectionName, limit, []entity.Vector{entity.FloatVector(vectors)}).WithOutputFields("id", "document_id", "content")
	results, err := s.client.Search(ctx, option)
	if err != nil {
		return nil, err
	}
	return results, nil

}

func (s *MilvusStore) QueryContentsByDocID(ctx context.Context, collectionName string, docID string) ([]string, error) {
	err := s.LoadCollection(ctx, collectionName)
	if err != nil {
		return nil, err
	}

	option := milvusclient.NewQueryOption(collectionName).
		WithFilter("document_id == {docID}").
		WithTemplateParam("docID", docID).
		WithOutputFields("content")
	result, err := s.client.Query(ctx, option)
	if err != nil {
		return nil, err
	}

	contentColumn := result.GetColumn("content")
	if contentColumn == nil {
		return nil, errors.New("content column not found in query result")
	}

	contents := make([]string, 0, result.ResultCount)
	for i := 0; i < result.ResultCount; i++ {
		content, err := contentColumn.GetAsString(i)
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}

	return contents, nil
}
