package hooks

import (
	"context"
	"strings"
	"sync"

	"go.uber.org/zap"
	tmodel "trpc.group/trpc-go/trpc-agent-go/model"
)

// SummarizationHook monitors the message count in a conversation and, when
// it exceeds a configurable threshold, calls the LLM to produce a summary
// that replaces the older messages.
//
// This prevents the context window from overflowing and keeps the most
// relevant conversation context available for the agent.
type SummarizationHook struct {
	mu            sync.Mutex
	model         tmodel.Model
	maxMessages   int
	summarizedIDs map[string]bool // sessionIDs that have been summarized
}

// NewSummarizationHook creates a new SummarizationHook.
// maxMessages is the threshold at which summarization triggers (default 20).
func NewSummarizationHook(mdl tmodel.Model, maxMessages int) *SummarizationHook {
	if maxMessages <= 0 {
		maxMessages = 20
	}
	return &SummarizationHook{
		model:         mdl,
		maxMessages:   maxMessages,
		summarizedIDs: make(map[string]bool),
	}
}

// BeforeModel is the tRPC-Agent-Go BeforeModelCallbackStructured callback.
// It checks the message count and, if the threshold is exceeded, summarises
// the conversation history by calling the LLM to produce a compressed summary
// that replaces older messages.
func (h *SummarizationHook) BeforeModel(
	ctx context.Context,
	args *tmodel.BeforeModelArgs,
) (*tmodel.BeforeModelResult, error) {
	if h.model == nil || len(args.Request.Messages) <= h.maxMessages {
		return nil, nil
	}

	// Extract session identifier from context to avoid double-summarization.
	sessionID, _ := ctx.Value("sessionID").(string)
	if sessionID != "" {
		h.mu.Lock()
		if h.summarizedIDs[sessionID] {
			h.mu.Unlock()
			return nil, nil
		}
		h.mu.Unlock()
	}

	zap.L().Info("summarization hook: message threshold exceeded, summarizing",
		zap.String("sessionID", sessionID),
		zap.Int("messages", len(args.Request.Messages)),
		zap.Int("threshold", h.maxMessages),
	)

	// Keep the last few messages for recent context; summarize the rest.
	keepRecent := h.maxMessages / 2
	if keepRecent < 2 {
		keepRecent = 2
	}

	messagesToSummarize := args.Request.Messages[:len(args.Request.Messages)-keepRecent]
	recentMessages := args.Request.Messages[len(args.Request.Messages)-keepRecent:]

	// Build the summarization prompt.
	var dialogBuilder strings.Builder
	for _, msg := range messagesToSummarize {
		dialogBuilder.WriteString(string(msg.Role))
		dialogBuilder.WriteString(": ")
		dialogBuilder.WriteString(msg.Content)
		dialogBuilder.WriteString("\n")
	}

	summaryPrompt := `请将以下对话历史压缩为一段简短的摘要，保留关键信息和上下文。只输出摘要文本，不要包含任何前缀或说明。

对话历史：
` + dialogBuilder.String()

	summaryReq := &tmodel.Request{
		Messages: []tmodel.Message{
			tmodel.NewSystemMessage("你是一个对话摘要助手。你的任务是将对话历史压缩为简洁的摘要。"),
			tmodel.NewUserMessage(summaryPrompt),
		},
	}

	responseCh, err := h.model.GenerateContent(ctx, summaryReq)
	if err != nil {
		zap.L().Warn("summarization hook: LLM summarization call failed",
			zap.Error(err),
		)
		return nil, nil
	}

	var summaryText strings.Builder
	for resp := range responseCh {
		if resp.Error != nil {
			zap.L().Warn("summarization hook: LLM response error",
				zap.String("message", resp.Error.Message),
			)
			return nil, nil
		}
		if len(resp.Choices) > 0 && resp.Choices[0].Delta.Content != "" {
			summaryText.WriteString(resp.Choices[0].Delta.Content)
		}
	}

	summary := strings.TrimSpace(summaryText.String())
	if summary == "" {
		return nil, nil
	}

	// Replace old messages with a summary system message plus recent messages.
	newMessages := make([]tmodel.Message, 0, len(recentMessages)+1)
	newMessages = append(newMessages, tmodel.NewSystemMessage(
		"【对话历史摘要】\n"+summary,
	))
	newMessages = append(newMessages, recentMessages...)

	args.Request.Messages = newMessages

	if sessionID != "" {
		h.mu.Lock()
		h.summarizedIDs[sessionID] = true
		h.mu.Unlock()
	}

	zap.L().Info("summarization hook: conversation summarized",
		zap.String("sessionID", sessionID),
		zap.Int("originalMsgs", len(messagesToSummarize)+len(recentMessages)),
		zap.Int("summarizedMsgs", len(newMessages)),
	)

	return nil, nil
}

// Reset clears the summarization state for a session so it can be
// re-summarized on a subsequent run.
func (h *SummarizationHook) Reset(sessionID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.summarizedIDs, sessionID)
}
