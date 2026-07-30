package hooks

import (
	"context"
	"encoding/json"

	"github.com/phoenix-agent-go/internal/repository"

	"go.uber.org/zap"
	tmodel "trpc.group/trpc-go/trpc-agent-go/model"
)

// ProfileInjectionHook is a BEFORE_MODEL hook that reads the user's profile
// from the database and injects it as a system message, allowing the agent
// to personalise responses based on stored user preferences.
//
// The injected message format is:
//
//	"登录人画像信息：姓名{name}, 工号：{code}, 偏好：{preferences}"
type ProfileInjectionHook struct {
	profileRepo repository.UserProfileRepository
}

// NewProfileInjectionHook creates a new ProfileInjectionHook.
func NewProfileInjectionHook(profileRepo repository.UserProfileRepository) *ProfileInjectionHook {
	return &ProfileInjectionHook{profileRepo: profileRepo}
}

// BeforeModel is the tRPC-Agent-Go BeforeModelCallbackStructured callback invoked
// before each LLM call. It reads the user's profile from the store and injects
// it as a system message at the beginning of the message list.
func (h *ProfileInjectionHook) BeforeModel(
	ctx context.Context,
	args *tmodel.BeforeModelArgs,
) (*tmodel.BeforeModelResult, error) {
	// Extract userID and agentSN from context values set by the caller.
	userID, _ := ctx.Value("userID").(string)
	agentSN, _ := ctx.Value("agentSN").(string)

	if userID == "" || agentSN == "" || h.profileRepo == nil {
		return nil, nil
	}

	profile, err := h.profileRepo.FindByUserIDAndAgentSN(ctx, userID, agentSN)
	if err != nil {
		zap.L().Warn("profile hook: failed to read user profile",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, nil
	}

	if profile == nil || profile.ProfileData == "" {
		return nil, nil
	}

	// Parse profile data (stored as JSON).
	var profileMap map[string]string
	if err := json.Unmarshal([]byte(profile.ProfileData), &profileMap); err != nil {
		zap.L().Warn("profile hook: failed to parse profile data",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, nil
	}

	// Build the profile injection string.
	profileText := h.buildProfileText(profileMap)
	if profileText == "" {
		return nil, nil
	}

	// Prepend the profile message to the existing messages.
	systemMsg := tmodel.NewSystemMessage("登录人画像信息：" + profileText)
	args.Request.Messages = append([]tmodel.Message{systemMsg}, args.Request.Messages...)

	return nil, nil
}

// buildProfileText constructs the human-readable profile string from profile data.
func (h *ProfileInjectionHook) buildProfileText(profile map[string]string) string {
	name := profile["name"]
	code := profile["code"]
	preferences := profile["preferences"]

	parts := make([]string, 0, 3)
	if name != "" {
		parts = append(parts, "姓名"+name)
	}
	if code != "" {
		parts = append(parts, "工号："+code)
	}
	if preferences != "" {
		parts = append(parts, "偏好："+preferences)
	}

	if len(parts) == 0 {
		return ""
	}

	// Join with commas.
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += ", " + parts[i]
	}
	return result
}
