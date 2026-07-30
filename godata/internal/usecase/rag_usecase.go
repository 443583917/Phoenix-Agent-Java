package usecase

import (
	"context"
	"strconv"

	"github.com/phoenix-agent-go/infra/id"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/repository"
)

// ──────────────────────────── RAG Errors ────────────────────────────

var (
	ErrRagCategoryNotFound = &AppError{Code: 701001, Msg: "RAG分类不存在"}
	ErrRagFileNotFound     = &AppError{Code: 702001, Msg: "RAG文件不存在"}
)

// ──────────────────────────── RagUsecase ────────────────────────────

type RagUsecase struct {
	categoryRepo repository.RagCategoryRepository
	fileRepo     repository.RagFileInfoRepository
}

func NewRagUsecase(
	categoryRepo repository.RagCategoryRepository,
	fileRepo repository.RagFileInfoRepository,
) *RagUsecase {
	return &RagUsecase{
		categoryRepo: categoryRepo,
		fileRepo:     fileRepo,
	}
}

func ragGenID() string {
	return strconv.FormatUint(id.MustGenerateID(), 10)
}

// ──────────────────────────── RagCategory ────────────────────────────

func (u *RagUsecase) CreateRagCategory(ctx context.Context, entity *model.RagCategory) (*model.RagCategory, error) {
	entity.ID = ragGenID()
	if err := u.categoryRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *RagUsecase) UpdateRagCategory(ctx context.Context, entity *model.RagCategory) error {
	existing, err := u.categoryRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrRagCategoryNotFound
	}
	return u.categoryRepo.Update(ctx, entity)
}

func (u *RagUsecase) DeleteRagCategory(ctx context.Context, id string) error {
	existing, err := u.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrRagCategoryNotFound
	}
	return u.categoryRepo.Delete(ctx, id)
}

func (u *RagUsecase) GetRagCategoryByID(ctx context.Context, id string) (*model.RagCategory, error) {
	entity, err := u.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrRagCategoryNotFound
	}
	return entity, nil
}

func (u *RagUsecase) PageRagCategory(ctx context.Context, page, size int, query *model.RagCategory) ([]*model.RagCategory, int64, error) {
	return u.categoryRepo.Page(ctx, page, size, query)
}

func (u *RagUsecase) ListRagCategory(ctx context.Context) ([]*model.RagCategory, error) {
	return u.categoryRepo.List(ctx)
}

// ──────────────────────────── RagFileInfo ────────────────────────────

func (u *RagUsecase) CreateRagFile(ctx context.Context, entity *model.RagFileInfo) (*model.RagFileInfo, error) {
	entity.ID = ragGenID()
	if err := u.fileRepo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *RagUsecase) UpdateRagFile(ctx context.Context, entity *model.RagFileInfo) error {
	existing, err := u.fileRepo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrRagFileNotFound
	}
	return u.fileRepo.Update(ctx, entity)
}

func (u *RagUsecase) DeleteRagFile(ctx context.Context, id string) error {
	existing, err := u.fileRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrRagFileNotFound
	}
	return u.fileRepo.Delete(ctx, id)
}

func (u *RagUsecase) GetRagFileByID(ctx context.Context, id string) (*model.RagFileInfo, error) {
	entity, err := u.fileRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, ErrRagFileNotFound
	}
	return entity, nil
}

func (u *RagUsecase) PageRagFile(ctx context.Context, page, size int, query *model.RagFileInfo) ([]*model.RagFileInfo, int64, error) {
	return u.fileRepo.Page(ctx, page, size, query)
}

func (u *RagUsecase) ListRagFile(ctx context.Context) ([]*model.RagFileInfo, error) {
	return u.fileRepo.List(ctx)
}
