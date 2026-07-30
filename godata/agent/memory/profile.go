package memory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/phoenix-agent-go/internal/model"
	"gorm.io/gorm"
)

// UserProfile manages user profile data backed by the UserProfileInfo table.
//
// Profile data is serialised as JSON into the ProfileData column. This
// provides a flexible, schema-less key-value store for user preferences
// and behavioural signals used by agents to personalise responses.
//
// NOTE: The tRPC-Agent-Go framework does not provide a direct equivalent
// for user profile management. This is a business-layer concern that
// remains backed by GORM + PostgreSQL. The framework's memory.Service
// handles fact/episode memories, not structured user preference profiles.
type UserProfile struct {
	db *gorm.DB
}

// NewUserProfile creates a new UserProfile manager.
func NewUserProfile(db *gorm.DB) *UserProfile {
	return &UserProfile{db: db}
}

// GetProfile retrieves the user profile as a key-value map.
// Returns nil if no profile exists for the user.
func (p *UserProfile) GetProfile(ctx context.Context, userID string) (map[string]string, error) {
	var entity model.UserProfileInfo
	err := p.db.WithContext(ctx).
		Where("user_id = ? AND del_flag = 0", userID).
		First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}

	if entity.ProfileData == "" {
		return make(map[string]string), nil
	}

	var profile map[string]string
	if err := json.Unmarshal([]byte(entity.ProfileData), &profile); err != nil {
		return nil, fmt.Errorf("unmarshal profile: %w", err)
	}
	return profile, nil
}

// UpdateProfile merges the given key-value updates into the user's profile,
// creating the profile record if it does not yet exist.
func (p *UserProfile) UpdateProfile(ctx context.Context, userID string, updates map[string]string) error {
	// Fetch existing profile.
	existing, err := p.GetProfile(ctx, userID)
	if err != nil {
		return err
	}
	if existing == nil {
		existing = make(map[string]string)
	}

	// Merge updates.
	for k, v := range updates {
		existing[k] = v
	}

	data, err := json.Marshal(existing)
	if err != nil {
		return fmt.Errorf("marshal profile: %w", err)
	}

	// Upsert: try to update existing record first.
	var entity model.UserProfileInfo
	result := p.db.WithContext(ctx).
		Where("user_id = ? AND del_flag = 0", userID).
		First(&entity)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("find profile: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		entity.ProfileData = string(data)
		return p.db.WithContext(ctx).Model(&entity).Update("profile_data", string(data)).Error
	}

	// Not found — create new record.
	entity = model.UserProfileInfo{
		UserID:      userID,
		ProfileData: string(data),
	}
	return p.db.WithContext(ctx).Create(&entity).Error
}
