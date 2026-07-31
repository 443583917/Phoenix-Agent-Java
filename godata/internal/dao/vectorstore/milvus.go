package vectorstore

import (
	"context"
	"fmt"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/phoenix-agent-go/internal/config"
	"go.uber.org/zap"
)

const (
	defaultCollectionName = "phoenix_vectors"
	fieldID               = "id"
	fieldContent          = "content"
	fieldMetadata         = "metadata"
	fieldEmbedding        = "embedding"
)

type Document struct {
	ID       string                 `json:"id"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Score    float64                `json:"score"`
}

type MilvusStore struct {
	client     client.Client
	cfg        config.VectorStoreConfig
	collection string
	logger     *zap.Logger
}

func NewMilvusStore(addr string, cfg config.VectorStoreConfig) (*MilvusStore, error) {
	c, err := client.NewClient(context.Background(), client.Config{
		Address: addr,
	})
	if err != nil {
		return nil, fmt.Errorf("milvus connect: %w", err)
	}

	s := &MilvusStore{
		client:     c,
		cfg:        cfg,
		collection: defaultCollectionName,
		logger:     zap.L().Named("vectorstore.milvus"),
	}

	if err := s.ensureCollection(); err != nil {
		c.Close()
		return nil, err
	}

	return s, nil
}

func (s *MilvusStore) ensureCollection() error {
	ctx := context.Background()
	exists, err := s.client.HasCollection(ctx, s.collection)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	schema := &entity.Schema{
		CollectionName: s.collection,
		Description:    "Phoenix Agent vector store",
		AutoID:         false,
		Fields: []*entity.Field{
			{
				Name:       fieldID,
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: true,
				TypeParams: map[string]string{"max_length": "128"},
			},
			{
				Name:       fieldContent,
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "65535"},
			},
			{
				Name:       fieldMetadata,
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "65535"},
			},
			{
				Name:     fieldEmbedding,
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": fmt.Sprintf("%d", s.cfg.Dimensions),
				},
			},
		},
	}

	if err := s.client.CreateCollection(ctx, schema, entity.DefaultShardNumber); err != nil {
		return fmt.Errorf("create collection: %w", err)
	}

	idx, err := entity.NewIndexHNSW(entity.COSINE, 128, 256)
	if err != nil {
		return err
	}
	if err := s.client.CreateIndex(ctx, s.collection, fieldEmbedding, idx, false); err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	s.logger.Info("milvus collection created", zap.String("collection", s.collection))
	return nil
}

func (s *MilvusStore) Search(ctx context.Context, query string, embedding []float64, topK int) ([]Document, error) {
	if topK <= 0 {
		topK = s.cfg.DefaultTopkLimit
	}
	if len(embedding) == 0 {
		return nil, nil
	}

	vec := toFloat32Vector(embedding)
	sp, _ := entity.NewIndexHNSWSearchParam(64)

	results, err := s.client.Search(
		ctx, s.collection, nil, "",
		[]string{fieldID, fieldContent, fieldMetadata},
		[]entity.Vector{entity.FloatVector(vec)},
		fieldEmbedding, entity.COSINE, topK, sp,
	)
	if err != nil {
		return nil, fmt.Errorf("milvus search: %w", err)
	}

	var docs []Document
	for _, result := range results {
		idCol, _ := result.Fields.GetColumn(fieldID).(*entity.ColumnVarChar)
		contentCol, _ := result.Fields.GetColumn(fieldContent).(*entity.ColumnVarChar)
		for i := 0; i < result.ResultCount; i++ {
			doc := Document{Score: float64(result.Scores[i])}
			if idCol != nil {
				doc.ID = idCol.Data()[i]
			}
			if contentCol != nil {
				doc.Content = contentCol.Data()[i]
			}
			docs = append(docs, doc)
		}
	}
	return docs, nil
}

func (s *MilvusStore) AddDocuments(ctx context.Context, docs []Document) error {
	if len(docs) == 0 {
		return nil
	}

	ids := make([]string, len(docs))
	contents := make([]string, len(docs))
	metadatas := make([]string, len(docs))
	embeddings := make([][]float32, len(docs))

	for i, doc := range docs {
		ids[i] = doc.ID
		contents[i] = doc.Content
		metadatas[i] = "{}"
		embeddings[i] = make([]float32, s.cfg.Dimensions)
	}

	idCol := entity.NewColumnVarChar(fieldID, ids)
	contentCol := entity.NewColumnVarChar(fieldContent, contents)
	metadataCol := entity.NewColumnVarChar(fieldMetadata, metadatas)
	embeddingCol := entity.NewColumnFloatVector(fieldEmbedding, s.cfg.Dimensions, embeddings)

	_, err := s.client.Insert(ctx, s.collection, "", idCol, contentCol, metadataCol, embeddingCol)
	if err != nil {
		return fmt.Errorf("milvus insert: %w", err)
	}

	return s.client.Flush(ctx, s.collection, false)
}

func (s *MilvusStore) DeleteDocuments(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	expr := fmt.Sprintf("id in [%s]", formatIDList(ids))
	return s.client.Delete(ctx, s.collection, "", expr)
}

func (s *MilvusStore) DeleteByFilter(ctx context.Context, filter map[string]string) error {
	for k, v := range filter {
		expr := fmt.Sprintf(`metadata like '%%"%s":"%s"%%'`, k, v)
		if err := s.client.Delete(ctx, s.collection, "", expr); err != nil {
			return err
		}
	}
	return nil
}

func (s *MilvusStore) Close() error {
	return s.client.Close()
}

func formatIDList(ids []string) string {
	result := ""
	for i, id := range ids {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf(`"%s"`, id)
	}
	return result
}

func toFloat32Vector(v []float64) []float32 {
	out := make([]float32, len(v))
	for i, f := range v {
		out[i] = float32(f)
	}
	return out
}
