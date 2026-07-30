package service

import (
	"context"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/usecase"
)

// KgService is a thin pass-through wrapper around KgUsecase.
type KgService struct {
	uc *usecase.KgUsecase
}

func NewKgService(uc *usecase.KgUsecase) *KgService {
	return &KgService{uc: uc}
}

// ──────────────────────────── KGEntity ────────────────────────────

func (s *KgService) CreateKGEntity(ctx context.Context, entity *model.KGEntity) (*model.KGEntity, error) {
	return s.uc.CreateKGEntity(ctx, entity)
}
func (s *KgService) UpdateKGEntity(ctx context.Context, entity *model.KGEntity) error {
	return s.uc.UpdateKGEntity(ctx, entity)
}
func (s *KgService) DeleteKGEntity(ctx context.Context, id string) error {
	return s.uc.DeleteKGEntity(ctx, id)
}
func (s *KgService) GetKGEntityByID(ctx context.Context, id string) (*model.KGEntity, error) {
	return s.uc.GetKGEntityByID(ctx, id)
}
func (s *KgService) PageKGEntity(ctx context.Context, page, size int, query *model.KGEntity) ([]*model.KGEntity, int64, error) {
	return s.uc.PageKGEntity(ctx, page, size, query)
}
func (s *KgService) ListKGEntity(ctx context.Context) ([]*model.KGEntity, error) {
	return s.uc.ListKGEntity(ctx)
}

// ──────────────────────────── KGRelation ────────────────────────────

func (s *KgService) CreateKGRelation(ctx context.Context, relation *model.KGRelation) (*model.KGRelation, error) {
	return s.uc.CreateKGRelation(ctx, relation)
}
func (s *KgService) UpdateKGRelation(ctx context.Context, relation *model.KGRelation) error {
	return s.uc.UpdateKGRelation(ctx, relation)
}
func (s *KgService) DeleteKGRelation(ctx context.Context, id string) error {
	return s.uc.DeleteKGRelation(ctx, id)
}
func (s *KgService) GetKGRelationByID(ctx context.Context, id string) (*model.KGRelation, error) {
	return s.uc.GetKGRelationByID(ctx, id)
}
func (s *KgService) PageKGRelation(ctx context.Context, page, size int, query *model.KGRelation) ([]*model.KGRelation, int64, error) {
	return s.uc.PageKGRelation(ctx, page, size, query)
}
func (s *KgService) FindRelationsByEntityID(ctx context.Context, entityID string) ([]*model.KGRelation, error) {
	return s.uc.FindRelationsByEntityID(ctx, entityID)
}

// ──────────────────────────── KGDomain ────────────────────────────

func (s *KgService) CreateKGDomain(ctx context.Context, domain *model.KGDomain) (*model.KGDomain, error) {
	return s.uc.CreateKGDomain(ctx, domain)
}
func (s *KgService) UpdateKGDomain(ctx context.Context, domain *model.KGDomain) error {
	return s.uc.UpdateKGDomain(ctx, domain)
}
func (s *KgService) DeleteKGDomain(ctx context.Context, id string) error {
	return s.uc.DeleteKGDomain(ctx, id)
}
func (s *KgService) GetKGDomainByID(ctx context.Context, id string) (*model.KGDomain, error) {
	return s.uc.GetKGDomainByID(ctx, id)
}
func (s *KgService) PageKGDomain(ctx context.Context, page, size int, query *model.KGDomain) ([]*model.KGDomain, int64, error) {
	return s.uc.PageKGDomain(ctx, page, size, query)
}
func (s *KgService) ListKGDomain(ctx context.Context) ([]*model.KGDomain, error) {
	return s.uc.ListKGDomain(ctx)
}
