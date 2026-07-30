# Phase 5C+5D Combined: Chat/Session/Graph SSE + Config CRUD + NL2SQL Stubs

## Summary

Phase 5C+5D created all remaining handler files (9 handlers, 53+ endpoints), NL2SQL workflow node stubs (7 nodes + graph builder + shared types), and registered everything in the router. Build passes cleanly with zero errors.

## Files Created (20 total)

### 5C: Chat + Session + Graph (3 handlers)

| File | Handler | Endpoints |
|------|---------|-----------|
| `api/handler/chat/chat.go` | ChatHandler | 9: sessions CRUD, messages, pin/rename/reports |
| `api/handler/chat/graph.go` | GraphHandler | 1: SSE `/stream/search` for NL2SQL graph execution |
| `api/handler/chat/session_event.go` | SessionEventHandler | 1: SSE `/agent/:agentId/sessions/stream` |

### 5D: Config CRUD + Remaining Handlers (6 handlers)

| File | Handler | Endpoints |
|------|---------|-----------|
| `api/handler/datasource/datasource.go` | DatasourceHandler | 13: CRUD, tables/columns, test connection, logical relations |
| `api/handler/modelconfig/model_config.go` | ModelConfigHandler | 7: list/add/update/delete/activate/test/check-ready |
| `api/handler/prompt/prompt_config.go` | PromptConfigHandler | 14: save, list-by-type, active, batch, priority, display-order |
| `api/handler/semanticmodel/semantic_model.go` | SemanticModelHandler | 11: CRUD, batch, enable/disable, batch-import, excel import |
| `api/handler/knowledge/business_knowledge.go` | BusinessKnowledgeHandler | 8: CRUD, recall, refresh-vector-store, retry-embedding |
| `api/handler/common/file_upload.go` | FileUploadHandler | 2: avatar upload, static file serve |

### NL2SQL Workflow Stubs (9 files)

| File | Purpose |
|------|---------|
| `agent/workflows/nl2sql/types/types.go` | Shared types: NL2SQLState, Node interface, NodeOutput, EvidenceContext, SchemaContext, ExecutionPlan |
| `agent/workflows/nl2sql/graph.go` | StateGraph builder: wires 7 nodes in order |
| `agent/workflows/nl2sql/nodes/intent.go` | IntentRecognitionNode: classifies intent (stub returns "sql", 0.95 confidence) |
| `agent/workflows/nl2sql/nodes/evidence.go` | EvidenceRecallNode: recalls evidence/knowledge (stub) |
| `agent/workflows/nl2sql/nodes/schema.go` | SchemaRecallNode: recalls schema info (stub) |
| `agent/workflows/nl2sql/nodes/planner.go` | PlannerNode: generates execution plan (stub, handles rejected plans) |
| `agent/workflows/nl2sql/nodes/sql_gen.go` | SqlGenerateNode: generates SQL (stub, branch for NL2SQL-only mode) |
| `agent/workflows/nl2sql/nodes/python_exec.go` | PythonExecuteNode: executes Python (stub) |
| `agent/workflows/nl2sql/nodes/report.go` | ReportGeneratorNode: generates HTML report (stub) |

### Files Modified (2)

| File | Changes |
|------|---------|
| `api/router.go` | +4 imports (chat, modelconfig, prompt, semanticmodel); +~100 lines: 53 new endpoints registered under `/api` with Auth+RBAC |
| `.superpowers/sdd/progress.md` | Updated progress tracking |

## Architecture Notes

1. **Shared types in sub-package**: To avoid import cycles between `nl2sql` graph and `nodes` packages, the shared types (NL2SQLState, Node interface, NodeOutput, etc.) reside in `nl2sql/types`. The `nodes` package imports `nl2sql/types` aliased as `nl2sql` so the qualifier `nl2sql.NodeOutput` remains consistent.

2. **All stubs return placeholder data**: Every handler method returns `response.Success(c, ...)` with descriptive stub data. SSE endpoints write a single `data:` frame. Full implementation with LLM-powered logic deferred to Phase 5 detailed iteration.

3. **Route groups under `/api`**: All new endpoints are registered under the existing `dataAPI` Gin router group, which applies JWT auth and RBAC middleware. Routes follow the existing naming conventions:
   - `/api/agent/:id/sessions` and `/api/sessions/:id/...` for chat
   - `/api/stream/search` for graph SSE (distinct from `/api/admin/agent/stream/chatsql`)
   - `/api/datasource/...` for datasource CRUD (distinct from `/api/agent/:agentId/datasource`)
   - `/api/model-config/...`, `/api/prompt-config/...`, etc. for config CRUD

4. **Dependency injection**: All new handlers use `*service.DataService` which is already wired in `main.go` via `buildDataService()`. No changes to `main.go` signature needed.

## Build

`go build ./...` -- compiled successfully with zero errors.

## Commit

`80f23da` - feat: Phase 5C+5D - Chat/SSE handlers, Config CRUD handlers, NL2SQL workflow stubs
