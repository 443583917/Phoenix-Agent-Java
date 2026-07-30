# Phoenix Go Phase 4 — Agent Framework Plan

**Goal:** Build the AI Agent engine using tRPC-Agent-Go + 2 SSE streaming controllers

**Architecture:** tRPC-Agent-Go (agents/tools/runner) → agent/ layer (manager, memory, knowledge) → internal/ (domain, repos) → SSE handlers

**Java controllers to match:**
- ReactAgentController: `POST /api/admin/agent/chat` (SSE), `GET /api/admin/agent/stream/chatsql` (SSE)
- HarnessController: `POST /api/admin/harness/confirm` (SSE), `POST /api/admin/harness/chat` (SSE with tool confirm)

---

### Task 1: Agent Entities + DTOs

Files:
- `internal/model/agent_entity.go` — GORM entities:
  - UserAgentInfo (tbl_agent_user_agent_info): id, userId, agentId, usageCount
  - UserMemoryInfo (tbl_agent_user_memory_info): id, userId, memoryType, content, embedding
  - UserProfileInfo (tbl_agent_user_profile_info): id, userId, profileData
  - CombinedStore (tbl_agent_combined_store): id, content, embedding, metadata
  All embed PrivilegeBaseModel from Phase 2.

- `internal/model/agent_dto.go` — DTOs:
  - ChatModelRequest: sessionId, content, agentSn
  - ConfirmRequest: agentSn, userId, sessionId, confirmed
  - HarnessRequest: harnessSn, userId, sessionId, message
  - SSE event map structure

Build `go build ./internal/model/`. Commit.

### Task 2: Agent Runtime + Agents + Tools

Files:
- `agent/runtime/manager.go` — AgentManager: manages agent instances, streamCall(agentSn, message) → SSE channel
- `agent/runtime/registry.go` — Registry: register/list agents by SN
- `agent/agents/react_agent.go` — ReactAgent wrapper over tRPC-Agent-Go llmagent
- `agent/agents/workflow_agent.go` — WorkflowAgent stub (Phase 5 fills in)
- `agent/agents/assistant_agent.go` — Simple LLM assistant
- `agent/tools/registry.go` — ToolRegistry: register tools, get by name
- `agent/tools/function/calculator.go` — Calculator tool (from existing main.go example)

AgentManager.StreamCall pseudo-code:
```go
func (m *AgentManager) StreamCall(ctx context.Context, req StreamRequest) (<-chan SSEEvent, error) {
    agent := m.registry.Get(req.AgentSN)
    events, err := agent.Run(ctx, req.UserID, req.SessionID, req.Message)
    // wrap runner events into SSE channel
    return ch, err
}
```

Use `trpc.group/trpc-go/trpc-agent-go` (already in go.mod).

Build. Commit.

### Task 3: Memory + Knowledge + Runner

Files:
- `agent/memory/short_term.go` — Conversation window (in-memory map + Redis)
- `agent/memory/long_term.go` — Long-term memory via Milvus vector search
- `agent/memory/profile.go` — User profile management
- `agent/knowledge/retriever.go` — Hybrid search (vector via Milvus + keyword)
- `agent/knowledge/embedding.go` — Embedding model proxy (OpenAI-compatible)
- `agent/knowledge/splitter.go` — Text splitter for documents
- `agent/runner/runner.go` — Conversation runner wrapping tRPC-Agent-Go runner
- `agent/runner/sse.go` — SSE event stream writer
- `agent/runner/hitl.go` — Human-in-the-Loop event handling (confirm/reject)

Build. Commit.

### Task 4: API Handlers + Routes + Integration

Files:
- `api/handler/agent/react_agent.go` — ReactAgentHandler:
  - `POST /api/admin/agent/chat` — SSE streaming, accepts ChatModelRequest, returns SSE stream {"content":"...","end":false/true}
  - `GET /api/admin/agent/stream/chatsql` — SSE graph search (stub, full impl in Phase 5)

- `api/handler/agent/harness.go` — HarnessHandler:
  - `POST /api/admin/harness/chat` — SSE streaming with tool confirmation buttons
  - `POST /api/admin/harness/confirm` — SSE confirm handler

Register routes in router.go under `/api/admin/agent` and `/api/admin/harness` groups with Auth middleware.

Wire up main.go: init AgentManager, pass to handlers.

Build, vet, test. Start server, verify SSE streaming.
