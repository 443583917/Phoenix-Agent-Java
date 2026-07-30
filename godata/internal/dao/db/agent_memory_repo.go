package db

import (
	"context"
	"errors"
	"time"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/repository"
	"gorm.io/gorm"
)

// ──────────────────────────── UserProfile ────────────────────────────

type userProfileRepo struct{ db *gorm.DB }

func NewUserProfileRepository(db *gorm.DB) repository.UserProfileRepository {
	return &userProfileRepo{db}
}

func (r *userProfileRepo) FindByUserIDAndAgentSN(ctx context.Context, userID, agentSN string) (*model.UserProfileInfo, error) {
	var entity model.UserProfileInfo
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND agent_sn = ? AND del_flag = 0", userID, agentSN).
		First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *userProfileRepo) SaveProfileUpdates(ctx context.Context, userID, agentSN string, profileData string) error {
	var existing model.UserProfileInfo
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND agent_sn = ? AND del_flag = 0", userID, agentSN).
		First(&existing)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	if result.RowsAffected > 0 {
		return r.db.WithContext(ctx).Model(&existing).Updates(map[string]interface{}{
			"profile_data": profileData,
			"update_time":  time.Now(),
		}).Error
	}

	entity := model.UserProfileInfo{
		UserID:      userID,
		AgentSn:     agentSN,
		ProfileData: profileData,
	}
	return r.db.WithContext(ctx).Create(&entity).Error
}

// ──────────────────────────── UserMemory ────────────────────────────

type userMemoryRepo struct{ db *gorm.DB }

func NewUserMemoryRepository(db *gorm.DB) repository.UserMemoryRepository {
	return &userMemoryRepo{db}
}

func (r *userMemoryRepo) Save(ctx context.Context, memory *model.UserMemoryInfo) error {
	return r.db.WithContext(ctx).Create(memory).Error
}

func (r *userMemoryRepo) FindByUserID(ctx context.Context, userID string) ([]*model.UserMemoryInfo, error) {
	var list []*model.UserMemoryInfo
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND del_flag = 0", userID).
		Order("create_time DESC").
		Find(&list).Error
	return list, err
}

// ──────────────────────────── UserAgentInfo ────────────────────────────

type userAgentInfoRepo struct{ db *gorm.DB }

func NewUserAgentInfoRepository(db *gorm.DB) repository.UserAgentInfoRepository {
	return &userAgentInfoRepo{db}
}

func (r *userAgentInfoRepo) FindByUserIDAndAgentSN(ctx context.Context, userID, agentSN string) (*model.UserAgentInfo, error) {
	var entity model.UserAgentInfo
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND agent_sn = ? AND del_flag = 0", userID, agentSN).
		First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *userAgentInfoRepo) RecordAction(ctx context.Context, userID, agentSN string) error {
	var existing model.UserAgentInfo
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND agent_sn = ? AND del_flag = 0", userID, agentSN).
		First(&existing)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	now := time.Now()
	if result.RowsAffected > 0 {
		return r.db.WithContext(ctx).Model(&existing).Updates(map[string]interface{}{
			"action_count": gorm.Expr("action_count + 1"),
			"last_date":    now,
			"update_time":  now,
		}).Error
	}

	count := int64(1)
	entity := model.UserAgentInfo{
		UserID:      userID,
		AgentSn:     agentSN,
		ActionCount: &count,
		LastDate:    &now,
	}
	return r.db.WithContext(ctx).Create(&entity).Error
}
