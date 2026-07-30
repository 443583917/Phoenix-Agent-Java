package repository

import (
	"context"

	"github.com/phoenix-agent-go/internal/model"
)

// ──────────────────────────── RagCategory ────────────────────────────

type RagCategoryRepository interface {
	FindByID(ctx context.Context, id string) (*model.RagCategory, error)
	FindByCode(ctx context.Context, code string) (*model.RagCategory, error)
	Page(ctx context.Context, page, size int, query *model.RagCategory) ([]*model.RagCategory, int64, error)
	List(ctx context.Context) ([]*model.RagCategory, error)
	Create(ctx context.Context, category *model.RagCategory) error
	Update(ctx context.Context, category *model.RagCategory) error
	Delete(ctx context.Context, id string) error
}

// ──────────────────────────── RagFileInfo ────────────────────────────

type RagFileInfoRepository interface {
	FindByID(ctx context.Context, id string) (*model.RagFileInfo, error)
	FindByCategoryID(ctx context.Context, categoryID string) ([]*model.RagFileInfo, error)
	Page(ctx context.Context, page, size int, query *model.RagFileInfo) ([]*model.RagFileInfo, int64, error)
	List(ctx context.Context) ([]*model.RagFileInfo, error)
	Create(ctx context.Context, file *model.RagFileInfo) error
	Update(ctx context.Context, file *model.RagFileInfo) error
	Delete(ctx context.Context, id string) error
}
