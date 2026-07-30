package usecase

import (
	"context"
	"strconv"

	"github.com/phoenix-agent-go/infra/id"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/repository"
)

// ──────────────────────────── KG Errors ────────────────────────────

var (
	ErrKGEntityNotFound   = &AppError{Code: 711001, Msg: "KG实体不存在"}
	ErrKGRelationNotFound = &AppError{Code: 712001, Msg: "KG关系不存在"}
	ErrKGDomainNotFound   = &AppError{Code: 713001, Msg: "KG领域不存在"}
)

// ──────────────────────────── KgUsecase ────────────────────────────

type KgUsecase struct {
	entityRepo   repository.KGEntityRepository
	relationRepo repository.KGRelationRepository
	domainRepo   repository.KGDomainRepository
}

func NewKgUsecase(
	entityRepo repository.KGEntityRepository,
	relationRepo repository.KGRelationRepository,
	domainRepo repository.KGDomainRepository,
) *KgUsecase {
	return &KgUsecase{
		entityRepo:   entityRepo,
		relationRepo: relationRepo,
		domainRepo:   domainRepo,
	}
}

func kgGenID() string {
	return strconv.FormatUint(id.MustGenerateID(), 10)
}

// ──────────────────────────── KGEntity ────────────────────────────

func (u *KgUsecase) CreateKGEntity(ctx context.Context, entity *model.KGEntity) (*model.KGEntity, error) {
	entity.ID = kgGenID()
	if err := u.entityRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *KgUsecase) UpdateKGEntity(ctx context.Context, entity *model.KGEntity) error {
	existing, err := u.entityRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrKGEntityNotFound
	}
	return u.entityRepo.Update(ctx, entity)
}

func (u *KgUsecase) DeleteKGEntity(ctx context.Context, id string) error {
	existing, err := u.entityRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrKGEntityNotFound
	}
	return u.entityRepo.Delete(ctx, id)
}

func (u *KgUsecase) GetKGEntityByID(ctx context.Context, id string) (*model.KGEntity, error) {
	entity, err := u.entityRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrKGEntityNotFound
	}
	return entity, nil
}

func (u *KgUsecase) PageKGEntity(ctx context.Context, page, size int, query *model.KGEntity) ([]*model.KGEntity, int64, error) {
	return u.entityRepo.Page(ctx, page, size, query)
}

func (u *KgUsecase) ListKGEntity(ctx context.Context) ([]*model.KGEntity, error) {
	return u.entityRepo.List(ctx)
}

// ──────────────────────────── KGRelation ────────────────────────────

func (u *KgUsecase) CreateKGRelation(ctx context.Context, relation *model.KGRelation) (*model.KGRelation, error) {
	relation.ID = kgGenID()
	if err := u.relationRepo.Create(ctx, relation); err != nil {
		return nil, err
	}
	return relation, nil
}

func (u *KgUsecase) UpdateKGRelation(ctx context.Context, relation *model.KGRelation) error {
	existing, err := u.relationRepo.FindByID(ctx, relation.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrKGRelationNotFound
	}
	return u.relationRepo.Update(ctx, relation)
}

func (u *KgUsecase) DeleteKGRelation(ctx context.Context, id string) error {
	existing, err := u.relationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrKGRelationNotFound
	}
	return u.relationRepo.Delete(ctx, id)
}

func (u *KgUsecase) GetKGRelationByID(ctx context.Context, id string) (*model.KGRelation, error) {
	relation, err := u.relationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if relation == nil {
		return nil, ErrKGRelationNotFound
	}
	return relation, nil
}

func (u *KgUsecase) PageKGRelation(ctx context.Context, page, size int, query *model.KGRelation) ([]*model.KGRelation, int64, error) {
	return u.relationRepo.Page(ctx, page, size, query)
}

func (u *KgUsecase) FindRelationsByEntityID(ctx context.Context, entityID string) ([]*model.KGRelation, error) {
	return u.relationRepo.FindByEntityID(ctx, entityID)
}

// ──────────────────────────── KGDomain ────────────────────────────

func (u *KgUsecase) CreateKGDomain(ctx context.Context, domain *model.KGDomain) (*model.KGDomain, error) {
	domain.ID = kgGenID()
	if err := u.domainRepo.Create(ctx, domain); err != nil {
		return nil, err
	}
	return domain, nil
}

func (u *KgUsecase) UpdateKGDomain(ctx context.Context, domain *model.KGDomain) error {
	existing, err := u.domainRepo.FindByID(ctx, domain.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrKGDomainNotFound
	}
	return u.domainRepo.Update(ctx, domain)
}

func (u *KgUsecase) DeleteKGDomain(ctx context.Context, id string) error {
	existing, err := u.domainRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrKGDomainNotFound
	}
	return u.domainRepo.Delete(ctx, id)
}

func (u *KgUsecase) GetKGDomainByID(ctx context.Context, id string) (*model.KGDomain, error) {
	domain, err := u.domainRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if domain == nil {
		return nil, ErrKGDomainNotFound
	}
	return domain, nil
}

func (u *KgUsecase) PageKGDomain(ctx context.Context, page, size int, query *model.KGDomain) ([]*model.KGDomain, int64, error) {
	return u.domainRepo.Page(ctx, page, size, query)
}

func (u *KgUsecase) ListKGDomain(ctx context.Context) ([]*model.KGDomain, error) {
	return u.domainRepo.List(ctx)
}
