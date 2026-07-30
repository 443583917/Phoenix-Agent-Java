package privilege

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/phoenix-agent-go/internal/service"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

type CheckTool struct {
	svc *service.PrivilegeService
}

func NewCheckTool(svc *service.PrivilegeService) *CheckTool {
	return &CheckTool{svc: svc}
}

func (t *CheckTool) Declaration() *tool.Declaration {
	return &tool.Declaration{
		Name:        "check_privilege",
		Description: "Check if a user has a specific privilege. Returns user status, active state, and role membership.",
		InputSchema: &tool.Schema{
			Type:     "object",
			Required: []string{"userId"},
			Properties: map[string]*tool.Schema{
				"userId": {
					Type:        "string",
					Description: "User ID to check",
				},
				"roleCode": {
					Type:        "string",
					Description: "Optional role name/code to verify against the user's assigned roles",
				},
			},
		},
	}
}

type checkParams struct {
	UserID   string `json:"userId"`
	RoleCode string `json:"roleCode"`
}

func (t *CheckTool) Call(ctx context.Context, args any) (any, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("marshal args: %w", err)
	}

	var params checkParams
	if err := json.Unmarshal(argsJSON, &params); err != nil {
		return nil, fmt.Errorf("unmarshal args: %w", err)
	}

	if t.svc == nil {
		return nil, fmt.Errorf("privilege service not initialized")
	}

	user, err := t.svc.GetUserByID(ctx, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	isActive := user.Status == 0
	result := map[string]any{
		"userId":    params.UserID,
		"username":  user.Username,
		"hasAccess": isActive,
		"status":    user.Status,
	}

	if params.RoleCode != "" {
		userRoles, err := t.svc.GetUserRolesByUserID(ctx, params.UserID)
		if err == nil {
			hasRole := false
			for _, ur := range userRoles {
				if strings.EqualFold(ur.RoleName, params.RoleCode) {
					hasRole = true
					break
				}
			}
			result["hasRole"] = hasRole
			result["roleCode"] = params.RoleCode
		} else {
			result["hasRole"] = false
			result["roleCode"] = params.RoleCode
		}
	}

	return result, nil
}
