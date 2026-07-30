package db

import (
	"context"
	"errors"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/repository"
	"gorm.io/gorm"
)

// ──────────────────────────── KGEntity ────────────────────────────

type kgEntityRepo struct{ db *gorm.DB }

func NewKGEntityRepository(db *gorm.DB) repository.KGEntityRepository {
	return &kgEntityRepo{db}
}

func (r *kgEntityRepo) FindByID(ctx context.Context, id string) (*model.KGEntity, error) {
	var entity model.KGEntity
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *kgEntityRepo) FindByDomainID(ctx context.Context, domainID string) ([]*model.KGEntity, error) {
	var list []*model.KGEntity
	err := r.db.WithContext(ctx).Where("domain_id = ? AND del_flag = 0", domainID).Order("create_time DESC").Find(&list).Error
	return list, err
}

func (r *kgEntityRepo) FindByType(ctx context.Context, typ string) ([]*model.KGEntity, error) {
	var list []*model.KGEntity
	err := r.db.WithContext(ctx).Where("type = ? AND del_flag = 0", typ).Order("create_time DESC").Find(&list).Error
	return list, err
}

func (r *kgEntityRepo) Page(ctx context.Context, page, size int, query *model.KGEntity) ([]*model.KGEntity, int64, error) {
	var list []*model.KGEntity
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.KGEntity{}).Where("del_flag = 0")
	if query != nil {
		if query.Name != "" {
			dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
		}
		if query.Type != "" {
			dbQuery = dbQuery.Where("type = ?", query.Type)
		}
		if query.DomainId != "" {
			dbQuery = dbQuery.Where("domain_id = ?", query.DomainId)
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Order("create_time DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *kgEntityRepo) List(ctx context.Context) ([]*model.KGEntity, error) {
	var list []*model.KGEntity
	err := r.db.WithContext(ctx).Where("del_flag = 0").Order("create_time DESC").Find(&list).Error
	return list, err
}

func (r *kgEntityRepo) Create(ctx context.Context, entity *model.KGEntity) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *kgEntityRepo) Update(ctx context.Context, entity *model.KGEntity) error {
	return r.db.WithContext(ctx).Model(&model.KGEntity{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *kgEntityRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.KGEntity{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── KGRelation ────────────────────────────

type kgRelationRepo struct{ db *gorm.DB }

func NewKGRelationRepository(db *gorm.DB) repository.KGRelationRepository {
	return &kgRelationRepo{db}
}

func (r *kgRelationRepo) FindByID(ctx context.Context, id string) (*model.KGRelation, error) {
	var entity model.KGRelation
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *kgRelationRepo) FindBySourceEntityID(ctx context.Context, sourceEntityID string) ([]*model.KGRelation, error) {
	var list []*model.KGRelation
	err := r.db.WithContext(ctx).Where("source_entity_id = ? AND del_flag = 0", sourceEntityID).Order("create_time DESC").Find(&list).Error
	return list, err
}

func (r *kgRelationRepo) FindByTargetEntityID(ctx context.Context, targetEntityID string) ([]*model.KGRelation, error) {
	var list []*model.KGRelation
	err := r.db.WithContext(ctx).Where("target_entity_id = ? AND del_flag = 0", targetEntityID).Order("create_time DESC").Find(&list).Error
	return list, err
}

func (r *kgRelationRepo) FindByEntityID(ctx context.Context, entityID string) ([]*model.KGRelation, error) {
	var list []*model.KGRelation
	err := r.db.WithContext(ctx).
		Where("(source_entity_id = ? OR target_entity_id = ?) AND del_flag = 0", entityID, entityID).
		Order("create_time DESC").Find(&list).Error
	return list, err
}

func (r *kgRelationRepo) Page(ctx context.Context, page, size int, query *model.KGRelation) ([]*model.KGRelation, int64, error) {
	var list []*model.KGRelation
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.KGRelation{}).Where("del_flag = 0")
	if query != nil {
		if query.RelationType != "" {
			dbQuery = dbQuery.Where("relation_type = ?", query.RelationType)
		}
		if query.SourceEntityId != "" {
			dbQuery = dbQuery.Where("source_entity_id = ?", query.SourceEntityId)
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Order("create_time DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *kgRelationRepo) Create(ctx context.Context, entity *model.KGRelation) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *kgRelationRepo) Update(ctx context.Context, entity *model.KGRelation) error {
	return r.db.WithContext(ctx).Model(&model.KGRelation{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *kgRelationRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.KGRelation{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── KGDomain ────────────────────────────

type kgDomainRepo struct{ db *gorm.DB }

func NewKGDomainRepository(db *gorm.DB) repository.KGDomainRepository {
	return &kgDomainRepo{db}
}

func (r *kgDomainRepo) FindByID(ctx context.Context, id string) (*model.KGDomain, error) {
	var entity model.KGDomain
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *kgDomainRepo) FindByCode(ctx context.Context, code string) (*model.KGDomain, error) {
	var entity model.KGDomain
	err := r.db.WithContext(ctx).Where("code = ? AND del_flag = 0", code).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *kgDomainRepo) Page(ctx context.Context, page, size int, query *model.KGDomain) ([]*model.KGDomain, int64, error) {
	var list []*model.KGDomain
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.KGDomain{}).Where("del_flag = 0")
	if query != nil {
		if query.Name != "" {
			dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
		}
		if query.Code != "" {
			dbQuery = dbQuery.Where("code LIKE ?", "%"+query.Code+"%")
		}
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	if err := dbQuery.Offset(offset).Limit(size).Order("create_time DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *kgDomainRepo) List(ctx context.Context) ([]*model.KGDomain, error) {
	var list []*model.KGDomain
	err := r.db.WithContext(ctx).Where("del_flag = 0").Order("create_time DESC").Find(&list).Error
	return list, err
}

func (r *kgDomainRepo) Create(ctx context.Context, entity *model.KGDomain) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *kgDomainRepo) Update(ctx context.Context, entity *model.KGDomain) error {
	return r.db.WithContext(ctx).Model(&model.KGDomain{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *kgDomainRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.KGDomain{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}
