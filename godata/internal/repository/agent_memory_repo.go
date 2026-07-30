package repository

import (
	"context"

	"github.com/phoenix-agent-go/internal/model"
)

// UserProfileRepository manages user profile data (tbl_agent_user_profile_info).
type UserProfileRepository interface {
	// FindByUserIDAndAgentSN retrieves the profile for a user-agent pair.
	FindByUserIDAndAgentSN(ctx context.Context, userID, agentSN string) (*model.UserProfileInfo, error)
	// SaveProfileUpdates creates or updates profile data for a user-agent pair.
	SaveProfileUpdates(ctx context.Context, userID, agentSN string, profileData string) error
}

// UserMemoryRepository manages user memory fact entries (tbl_agent_user_memory_info).
type UserMemoryRepository interface {
	// Save stores a new user memory fact entry.
	Save(ctx context.Context, memory *model.UserMemoryInfo) error
	// FindByUserID retrieves all memory entries for a user.
	FindByUserID(ctx context.Context, userID string) ([]*model.UserMemoryInfo, error)
}

// UserAgentInfoRepository manages user-agent usage records (tbl_agent_user_agent_info).
type UserAgentInfoRepository interface {
	// FindByUserIDAndAgentSN retrieves the usage record for a user-agent pair.
	FindByUserIDAndAgentSN(ctx context.Context, userID, agentSN string) (*model.UserAgentInfo, error)
	// RecordAction creates or updates the usage record, incrementing the action count.
	RecordAction(ctx context.Context, userID, agentSN string) error
}
