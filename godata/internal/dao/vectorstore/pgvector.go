package vectorstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/phoenix-agent-go/internal/config"
	"gorm.io/gorm"
)

type Document struct {
	ID       string                 `json:"id"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Score    float64                `json:"score"`
}

type PgVectorStore struct {
	db  *gorm.DB
	cfg config.VectorStoreConfig
}

func NewPgVectorStore(db *gorm.DB, cfg config.VectorStoreConfig) *PgVectorStore {
	return &PgVectorStore{db: db, cfg: cfg}
}

func (s *PgVectorStore) Search(ctx context.Context, query string, embedding []float64, topK int) ([]Document, error) {
	if topK <= 0 {
		topK = s.cfg.DefaultTopkLimit
	}
	vec := formatVector(embedding)
	var results []Document
	err := s.db.WithContext(ctx).Raw(
		`SELECT id, content, metadata::text, 1 - (embedding <=> ?) AS score
		 FROM tbl_vector_store_simple_data
		 WHERE 1 - (embedding <=> ?) > ?
		 ORDER BY embedding <=> ?
		 LIMIT ?`,
		vec, vec, s.cfg.SimilarityThreshold, vec, topK,
	).Scan(&results).Error
	return results, err
}

func (s *PgVectorStore) AddDocuments(ctx context.Context, docs []Document) error {
	for _, doc := range docs {
		if err := s.db.WithContext(ctx).Exec(
			`INSERT INTO tbl_vector_store_simple_data (id, content, metadata, embedding)
			 VALUES (?, ?, ?::jsonb, ?)
			 ON CONFLICT (id) DO UPDATE SET content = EXCLUDED.content, metadata = EXCLUDED.metadata`,
			doc.ID, doc.Content, "{}", "[]",
		).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *PgVectorStore) DeleteDocuments(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	if len(ids) > s.cfg.BatchDelTopkLimit {
		ids = ids[:s.cfg.BatchDelTopkLimit]
	}
	return s.db.WithContext(ctx).Where("id IN ?", ids).Delete(&Document{}).Error
}

func (s *PgVectorStore) DeleteByFilter(ctx context.Context, filter map[string]string) error {
	query := s.db.WithContext(ctx)
	for k, v := range filter {
		query = query.Where("metadata->>? = ?", k, v)
	}
	return query.Delete(&Document{}).Error
}

func formatVector(vec []float64) string {
	if len(vec) == 0 {
		return "[]"
	}
	parts := make([]string, len(vec))
	for i, v := range vec {
		parts[i] = fmt.Sprintf("%f", v)
	}
	return "[" + strings.Join(parts, ",") + "]"
}
