package db

import (
	"context"
	"errors"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/repository"
	"gorm.io/gorm"
)

// ──────────────────────────── RagCategory ────────────────────────────

type ragCategoryRepo struct{ db *gorm.DB }

func NewRagCategoryRepository(db *gorm.DB) repository.RagCategoryRepository {
	return &ragCategoryRepo{db}
}

func (r *ragCategoryRepo) FindByID(ctx context.Context, id string) (*model.RagCategory, error) {
	var entity model.RagCategory
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *ragCategoryRepo) FindByCode(ctx context.Context, code string) (*model.RagCategory, error) {
	var entity model.RagCategory
	err := r.db.WithContext(ctx).Where("code = ? AND del_flag = 0", code).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *ragCategoryRepo) Page(ctx context.Context, page, size int, query *model.RagCategory) ([]*model.RagCategory, int64, error) {
	var list []*model.RagCategory
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.RagCategory{}).Where("del_flag = 0")
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

func (r *ragCategoryRepo) List(ctx context.Context) ([]*model.RagCategory, error) {
	var list []*model.RagCategory
	err := r.db.WithContext(ctx).Where("del_flag = 0").Order("create_time DESC").Find(&list).Error
	return list, err
}

func (r *ragCategoryRepo) Create(ctx context.Context, entity *model.RagCategory) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *ragCategoryRepo) Update(ctx context.Context, entity *model.RagCategory) error {
	return r.db.WithContext(ctx).Model(&model.RagCategory{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *ragCategoryRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.RagCategory{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}

// ──────────────────────────── RagFileInfo ────────────────────────────

type ragFileInfoRepo struct{ db *gorm.DB }

func NewRagFileInfoRepository(db *gorm.DB) repository.RagFileInfoRepository {
	return &ragFileInfoRepo{db}
}

func (r *ragFileInfoRepo) FindByID(ctx context.Context, id string) (*model.RagFileInfo, error) {
	var entity model.RagFileInfo
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *ragFileInfoRepo) FindByCategoryID(ctx context.Context, categoryID string) ([]*model.RagFileInfo, error) {
	var list []*model.RagFileInfo
	err := r.db.WithContext(ctx).Where("category_id = ? AND del_flag = 0", categoryID).Order("create_time DESC").Find(&list).Error
	return list, err
}

func (r *ragFileInfoRepo) Page(ctx context.Context, page, size int, query *model.RagFileInfo) ([]*model.RagFileInfo, int64, error) {
	var list []*model.RagFileInfo
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.RagFileInfo{}).Where("del_flag = 0")
	if query != nil {
		if query.Name != "" {
			dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
		}
		if query.FileType != "" {
			dbQuery = dbQuery.Where("file_type = ?", query.FileType)
		}
		if query.CategoryId != "" {
			dbQuery = dbQuery.Where("category_id = ?", query.CategoryId)
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

func (r *ragFileInfoRepo) List(ctx context.Context) ([]*model.RagFileInfo, error) {
	var list []*model.RagFileInfo
	err := r.db.WithContext(ctx).Where("del_flag = 0").Order("create_time DESC").Find(&list).Error
	return list, err
}

func (r *ragFileInfoRepo) Create(ctx context.Context, entity *model.RagFileInfo) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *ragFileInfoRepo) Update(ctx context.Context, entity *model.RagFileInfo) error {
	return r.db.WithContext(ctx).Model(&model.RagFileInfo{}).Where("id = ? AND del_flag = 0", entity.ID).Updates(entity).Error
}

func (r *ragFileInfoRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.RagFileInfo{}).Where("id = ? AND del_flag = 0", id).Update("del_flag", 1).Error
}
