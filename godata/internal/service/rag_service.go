package service

import (
	"context"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/usecase"
)

// RagService is a thin pass-through wrapper around RagUsecase.
type RagService struct {
	uc *usecase.RagUsecase
}

func NewRagService(uc *usecase.RagUsecase) *RagService {
	return &RagService{uc: uc}
}

// ──────────────────────────── RagCategory ────────────────────────────

func (s *RagService) CreateRagCategory(ctx context.Context, entity *model.RagCategory) (*model.RagCategory, error) {
	return s.uc.CreateRagCategory(ctx, entity)
}

func (s *RagService) UpdateRagCategory(ctx context.Context, entity *model.RagCategory) error {
	return s.uc.UpdateRagCategory(ctx, entity)
}

func (s *RagService) DeleteRagCategory(ctx context.Context, id string) error {
	return s.uc.DeleteRagCategory(ctx, id)
}

func (s *RagService) GetRagCategoryByID(ctx context.Context, id string) (*model.RagCategory, error) {
	return s.uc.GetRagCategoryByID(ctx, id)
}

func (s *RagService) PageRagCategory(ctx context.Context, page, size int, query *model.RagCategory) ([]*model.RagCategory, int64, error) {
	return s.uc.PageRagCategory(ctx, page, size, query)
}

func (s *RagService) ListRagCategory(ctx context.Context) ([]*model.RagCategory, error) {
	return s.uc.ListRagCategory(ctx)
}

// ──────────────────────────── RagFileInfo ────────────────────────────

func (s *RagService) CreateRagFile(ctx context.Context, entity *model.RagFileInfo) (*model.RagFileInfo, error) {
	return s.uc.CreateRagFile(ctx, entity)
}

func (s *RagService) UpdateRagFile(ctx context.Context, entity *model.RagFileInfo) error {
	return s.uc.UpdateRagFile(ctx, entity)
}

func (s *RagService) DeleteRagFile(ctx context.Context, id string) error {
	return s.uc.DeleteRagFile(ctx, id)
}

func (s *RagService) GetRagFileByID(ctx context.Context, id string) (*model.RagFileInfo, error) {
	return s.uc.GetRagFileByID(ctx, id)
}

func (s *RagService) PageRagFile(ctx context.Context, page, size int, query *model.RagFileInfo) ([]*model.RagFileInfo, int64, error) {
	return s.uc.PageRagFile(ctx, page, size, query)
}

func (s *RagService) ListRagFile(ctx context.Context) ([]*model.RagFileInfo, error) {
	return s.uc.ListRagFile(ctx)
}
