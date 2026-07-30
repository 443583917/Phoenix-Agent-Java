package memory

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/repository"

	"go.uber.org/zap"
	tmodel "trpc.group/trpc-go/trpc-agent-go/model"
)

// MemoryExtractionResult is the structured output from the memory extraction LLM call.
type MemoryExtractionResult struct {
	HasMemory           bool              `json:"hasMemory"`
	UserProfileUpdates  map[string]string `json:"userProfileUpdates"`
	Facts               []string          `json:"facts"`
	Summary             string            `json:"summary"`
}

// VectorStore abstracts the vector storage backend for saving memory embeddings
// and performing semantic search.
type VectorStore interface {
	// SaveMemory stores a memory entry with its embedding for later semantic search.
	SaveMemory(ctx context.Context, userID, content string, topics []string) error
	// SearchMemories searches for relevant memories by semantic similarity.
	SearchMemories(ctx context.Context, userID, query string, topK int) ([]Memory, error)
}

// MemoryPipeline implements async memory extraction that runs after every
// agent conversation turn. It calls an LLM to extract user profile updates,
// factual memories, and conversation summaries, then persists them.
//
// Errors are logged and never propagated to the caller.
type MemoryPipeline struct {
	model       tmodel.Model
	profileRepo repository.UserProfileRepository
	memoryRepo  repository.UserMemoryRepository
	vectorStore VectorStore
}

// NewMemoryPipeline creates a new MemoryPipeline.
// model is the tRPC-Agent-Go model used for the memory extraction LLM call.
func NewMemoryPipeline(
	mdl tmodel.Model,
	profileRepo repository.UserProfileRepository,
	memoryRepo repository.UserMemoryRepository,
	vectorStore VectorStore,
) *MemoryPipeline {
	return &MemoryPipeline{
		model:       mdl,
		profileRepo: profileRepo,
		memoryRepo:  memoryRepo,
		vectorStore: vectorStore,
	}
}

// ProcessAndExtractMemory runs asynchronously after each agent conversation turn.
// It builds an extraction prompt, calls the LLM to extract memory insights,
// and persists profile updates, facts, and summaries.
//
// This method is designed to be called via "go" (fire-and-forget). All errors
// are logged at warn/error level and never returned to the caller.
func (p *MemoryPipeline) ProcessAndExtractMemory(
	ctx context.Context,
	userID, agentSN, userQuery string,
) {
	if p.model == nil {
		zap.L().Warn("memory pipeline: no model configured, skipping extraction")
		return
	}

	// 1. Build the extraction prompt.
	systemPrompt := `你是一个智能体的记忆提取引擎。你的任务是从用户的输入中提取有价值的信息，用于构建用户画像和长期记忆。

请分析用户的输入，提取以下信息并以 JSON 格式返回：

1. hasMemory (bool): 是否包含有价值的信息需要记忆
2. userProfileUpdates (map[string]string): 用户画像更新，例如姓名、偏好、角色等
3. facts ([]string): 用户陈述的事实信息列表
4. summary (string): 对用户输入的一句话摘要

请只返回 JSON，不要包含任何其他文本。如果用户输入中没有可提取的信息，将 hasMemory 设为 false。`

	userPrompt := "用户输入：" + userQuery

	// 2. Call the LLM for extraction.
	request := &tmodel.Request{
		Messages: []tmodel.Message{
			tmodel.NewSystemMessage(systemPrompt),
			tmodel.NewUserMessage(userPrompt),
		},
	}

	responseCh, err := p.model.GenerateContent(ctx, request)
	if err != nil {
		zap.L().Warn("memory pipeline: LLM extraction call failed",
			zap.String("userID", userID),
			zap.String("agentSN", agentSN),
			zap.Error(err),
		)
		return
	}

	// Collect the full response text.
	var fullResponse strings.Builder
	for resp := range responseCh {
		if resp.Error != nil {
			zap.L().Warn("memory pipeline: LLM response error",
				zap.String("userID", userID),
				zap.String("message", resp.Error.Message),
			)
			return
		}
		if len(resp.Choices) > 0 && resp.Choices[0].Delta.Content != "" {
			fullResponse.WriteString(resp.Choices[0].Delta.Content)
		}
	}

	responseText := strings.TrimSpace(fullResponse.String())
	if responseText == "" {
		return
	}

	// 3. Parse the extraction result.
	var result MemoryExtractionResult
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		// Try to extract JSON from within markdown code blocks.
		cleaned := extractJSON(responseText)
		if cleaned == "" || json.Unmarshal([]byte(cleaned), &result) != nil {
			zap.L().Warn("memory pipeline: failed to parse LLM response",
				zap.String("userID", userID),
				zap.String("response", responseText),
				zap.Error(err),
			)
			return
		}
	}

	// 4. If no memory to extract, return early.
	if !result.HasMemory {
		return
	}

	// 5. Save profile updates.
	if len(result.UserProfileUpdates) > 0 {
		profileData, err := json.Marshal(result.UserProfileUpdates)
		if err != nil {
			zap.L().Warn("memory pipeline: failed to marshal profile updates", zap.Error(err))
		} else {
			if err := p.profileRepo.SaveProfileUpdates(ctx, userID, agentSN, string(profileData)); err != nil {
				zap.L().Warn("memory pipeline: failed to save profile updates",
					zap.String("userID", userID),
					zap.Error(err),
				)
			}
		}
	}

	// 6. Save facts as user memory entries.
	for _, fact := range result.Facts {
		if strings.TrimSpace(fact) == "" {
			continue
		}
		memoryEntry := &model.UserMemoryInfo{
			UserID:     userID,
			AgentSn:    agentSN,
			MemoryType: "FACT",
			Content:    fact,
		}
		if err := p.memoryRepo.Save(ctx, memoryEntry); err != nil {
			zap.L().Warn("memory pipeline: failed to save memory fact",
				zap.String("userID", userID),
				zap.Error(err),
			)
		}
	}

	// 7. Save summary to vector store for semantic search.
	if result.Summary != "" && p.vectorStore != nil {
		if err := p.vectorStore.SaveMemory(ctx, userID, result.Summary, nil); err != nil {
			zap.L().Warn("memory pipeline: failed to save summary to vector store",
				zap.String("userID", userID),
				zap.Error(err),
			)
		}
	}

	zap.L().Info("memory pipeline: extraction complete",
		zap.String("userID", userID),
		zap.Int("facts", len(result.Facts)),
		zap.Int("profileUpdates", len(result.UserProfileUpdates)),
	)
}

// extractJSON attempts to find a JSON object within text, e.g. inside
// markdown code fences (```json ... ```).
func extractJSON(text string) string {
	// Try to find JSON between ```json and ``` markers.
	start := strings.Index(text, "```json")
	if start == -1 {
		start = strings.Index(text, "```")
	}
	if start != -1 {
		start = strings.Index(text[start:], "\n") + start + 1
		end := strings.Index(text[start:], "```")
		if end != -1 {
			return strings.TrimSpace(text[start : start+end])
		}
	}

	// Fallback: find the outermost { } pair.
	braceStart := strings.Index(text, "{")
	braceEnd := strings.LastIndex(text, "}")
	if braceStart != -1 && braceEnd > braceStart {
		return text[braceStart : braceEnd+1]
	}

	return ""
}
