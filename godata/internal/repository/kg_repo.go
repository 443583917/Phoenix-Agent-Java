package repository

import (
	"context"

	"github.com/phoenix-agent-go/internal/model"
)

// ──────────────────────────── KGEntity ────────────────────────────

type KGEntityRepository interface {
	FindByID(ctx context.Context, id string) (*model.KGEntity, error)
	FindByDomainID(ctx context.Context, domainID string) ([]*model.KGEntity, error)
	FindByType(ctx context.Context, typ string) ([]*model.KGEntity, error)
	Page(ctx context.Context, page, size int, query *model.KGEntity) ([]*model.KGEntity, int64, error)
	List(ctx context.Context) ([]*model.KGEntity, error)
	Create(ctx context.Context, entity *model.KGEntity) error
	Update(ctx context.Context, entity *model.KGEntity) error
	Delete(ctx context.Context, id string) error
}

// ──────────────────────────── KGRelation ────────────────────────────

type KGRelationRepository interface {
	FindByID(ctx context.Context, id string) (*model.KGRelation, error)
	FindBySourceEntityID(ctx context.Context, sourceEntityID string) ([]*model.KGRelation, error)
	FindByTargetEntityID(ctx context.Context, targetEntityID string) ([]*model.KGRelation, error)
	FindByEntityID(ctx context.Context, entityID string) ([]*model.KGRelation, error)
	Page(ctx context.Context, page, size int, query *model.KGRelation) ([]*model.KGRelation, int64, error)
	Create(ctx context.Context, relation *model.KGRelation) error
	Update(ctx context.Context, relation *model.KGRelation) error
	Delete(ctx context.Context, id string) error
}

// ──────────────────────────── KGDomain ────────────────────────────

type KGDomainRepository interface {
	FindByID(ctx context.Context, id string) (*model.KGDomain, error)
	FindByCode(ctx context.Context, code string) (*model.KGDomain, error)
	Page(ctx context.Context, page, size int, query *model.KGDomain) ([]*model.KGDomain, int64, error)
	List(ctx context.Context) ([]*model.KGDomain, error)
	Create(ctx context.Context, domain *model.KGDomain) error
	Update(ctx context.Context, domain *model.KGDomain) error
	Delete(ctx context.Context, id string) error
}
