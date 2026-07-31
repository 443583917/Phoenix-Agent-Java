# Phoenix Go 剩余功能实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 补全 Go 后端与 Java 后端之间的所有功能缺口，使 Go 版本功能等价于 Java 版本。

**Architecture:** 按优先级分为 4 个阶段：Phase 1 (纯 Go 快速补全) → Phase 2 (架构组件) → Phase 3 (外部集成) → Phase 4 (优化项)。每个 Task 独立可执行、独立可验证。

**Tech Stack:** Go 1.25, tRPC-Agent-Go v1.10.0, GORM, Gin, OpenTelemetry, Docker SDK

## Global Constraints

- `go build ./...` 和 `go vet ./...` 每个 Task 完成后必须 PASS
- 不引入新的外部服务依赖，除非 Task 明确说明
- 所有新增接口遵循现有 Handler → Service → Usecase → Repository 分层
- 配置文件使用 YAML，通过 Viper 加载

---

### Task 1: NL2SQL 状态键补全 + 条件终止边

**Files:**
- Modify: `agent/workflows/nl2sql/types/types.go`
- Modify: `agent/workflows/nl2sql/graph.go`

**说明:** 补全 Java 版 OverAllState 中缺失的 5 个状态键，以及 3 条缺失的条件终止边。

- [ ] **Step 1: 补全状态键**

在 `types/types.go` 的 `NL2SQLState` 结构体中添加：
```go
MultiTurnContext          string `json:"multiTurnContext,omitempty"`
TraceThreadID             string `json:"traceThreadId,omitempty"`
SQLSchemaMissingAdvice    string `json:"sqlSchemaMissingAdvice,omitempty"`
PythonFallbackMode        bool   `json:"pythonFallbackMode"`
Result                    string `json:"result,omitempty"`
```

- [ ] **Step 2: 添加缺失的条件终止边**

在 `graph.go` 的 `Build()` 方法中，将 `query_enhance → schema_recall` 改为条件边：空输出时 → END。将 `schema_recall → table_relation` 改为条件边：无表文档时 → END。将 `table_relation` 添加自环重试边（最多 3 次）。

- [ ] **Step 3: 构建验证**

Run: `go build ./... && go vet ./...`

---

### Task 2: 配置外部化

**Files:**
- Create: `internal/config/graph.go`
- Modify: `internal/config/app.go`
- Modify: `agent/workflows/nl2sql/graph.go`
- Modify: `agent/workflows/nl2sql/nodes/sql_gen.go`
- Modify: `agent/workflows/nl2sql/nodes/python_exec.go`

**说明:** 将 Java 的 `DataAgentProperties` 硬编码值提取为 Go 配置结构体。

- [ ] **Step 1: 创建 GraphConfig 结构体**

```go
// internal/config/graph.go
type GraphConfig struct {
    MaxSQLRetryCount       int     `mapstructure:"max_sql_retry_count"`
    MaxSQLOptimizeCount    int     `mapstructure:"max_sql_optimize_count"`
    SQLScoreThreshold      float64 `mapstructure:"sql_score_threshold"`
    MaxTurnHistory         int     `mapstructure:"max_turn_history"`
    MaxPlanLength          int     `mapstructure:"max_plan_length"`
    MaxColumnsPerTable     int     `mapstructure:"max_columns_per_table"`
    EnableSQLResultChart   bool    `mapstructure:"enable_sql_result_chart"`
    PythonMaxTriesCount    int     `mapstructure:"python_max_tries_count"`
    TableTopkLimit         int     `mapstructure:"table_topk_limit"`
    TableSimilarityThreshold float64 `mapstructure:"table_similarity_threshold"`
    DefaultSimilarityThreshold float64 `mapstructure:"default_similarity_threshold"`
    DefaultTopkLimit       int     `mapstructure:"default_topk_limit"`
}

func DefaultGraphConfig() GraphConfig {
    return GraphConfig{
        MaxSQLRetryCount: 10, MaxSQLOptimizeCount: 10,
        SQLScoreThreshold: 0.95, MaxTurnHistory: 5,
        MaxPlanLength: 2000, MaxColumnsPerTable: 150,
        EnableSQLResultChart: true, PythonMaxTriesCount: 5,
        TableTopkLimit: 10, TableSimilarityThreshold: 0.2,
        DefaultSimilarityThreshold: 0.4, DefaultTopkLimit: 8,
    }
}
```

- [ ] **Step 2: 替换硬编码值**

将 `sql_gen.go` 中的 `3` 替换为 `cfg.MaxSQLRetryCount`，`python_exec.go` 中的 `3` 替换为 `cfg.PythonMaxTriesCount`，`schema.go` 中的 `10/15` 替换为 `cfg.TableTopkLimit`。

- [ ] **Step 3: 添加 YAML 配置**

在 `configs/api/app.yaml` 中添加 `graph:` 节，包含所有默认值。

- [ ] **Step 4: 构建验证**

Run: `go build ./... && go vet ./...`

---

### Task 3: MultiTurnContextManager

**Files:**
- Create: `agent/runtime/multi_turn.go`
- Modify: `api/handler/chat/graph.go`
- Modify: `agent/workflows/nl2sql/types/types.go`

**说明:** 实现 Java 版 `MultiTurnContextManager` 的多轮对话上下文管理。

- [ ] **Step 1: 创建 MultiTurnContextManager**

```go
// agent/runtime/multi_turn.go
type MultiTurnContextManager struct {
    mu       sync.RWMutex
    contexts map[string]*TurnContext
    maxTurns int
}

type TurnContext struct {
    ThreadID       string
    History        []TurnEntry
    PendingChunks  strings.Builder
}

type TurnEntry struct {
    Query    string
    Response string
}
```

实现方法：`BeginTurn(threadID)`, `FinishTurn(threadID, query, response)`, `DiscardPending(threadID)`, `RestartLastTurn(threadID)`, `BuildContext(threadID) string`, `AppendPlannerChunk(threadID, chunk string)`。

- [ ] **Step 2: 接入 GraphHandler**

在 `GraphHandler` 中添加 `mtcm *runtime.MultiTurnContextManager` 字段。在 `StreamSearch` 中调用 `BuildContext()` 获取上下文并注入初始状态。完成后调用 `FinishTurn()`。

- [ ] **Step 3: 构建验证**

Run: `go build ./... && go vet ./...`

---

### Task 4: HTML 报告二进制下载

**Files:**
- Modify: `api/handler/chat/chat.go`
- Create: `internal/service/session/report_template.go`（已有，需添加下载方法）

**说明:** 将 `GenerateReportHTML` 从返回 JSON 改为返回二进制 HTML 文件下载。

- [ ] **Step 1: 修改 handler**

```go
func (h *ChatHandler) GenerateReportHTML(c *gin.Context) {
    sessionID := c.Param("id")
    messages, _ := h.svc.GetSessionMessages(c.Request.Context(), sessionID)
    html := session.BuildReportHTML(messages, sessionID)
    
    c.Header("Content-Type", "text/html; charset=utf-8")
    c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="report_%s.html"`, sessionID))
    c.String(http.StatusOK, html)
}
```

- [ ] **Step 2: 创建 BuildReportHTML 函数**

在 `report_template.go` 中实现完整的 HTML 模板生成，包含 CSS 样式、消息列表、时间戳等。

- [ ] **Step 3: 构建验证**

Run: `go build ./... && go vet ./...`

---

### Task 5: 缺失端点补全

**Files:**
- Modify: `api/handler/platform/account_tenant_info.go`
- Modify: `api/handler/platform/account_group_info.go`
- Modify: `api/handler/privilege/module.go`
- Modify: `api/router.go`

**说明:** 补全 Java 有但 Go 缺失的 4 个查询/管理端点。

- [ ] **Step 1: AccountTenantInfo — GetByAccountID**

```go
func (h *AccountTenantInfoHandler) GetByAccountID(c *gin.Context) {
    // GET /platform/account-tenant-info/account/:accountId
}
```

- [ ] **Step 2: AccountGroupInfo — GetByGroupID + GetByAccountID**

```go
func (h *AccountGroupInfoHandler) GetByGroupID(c *gin.Context) {}
func (h *AccountGroupInfoHandler) GetByAccountID(c *gin.Context) {}
```

- [ ] **Step 3: Module — GetPvalues + UpdatePvalue**

```go
func (h *ModuleHandler) GetPvalues(c *gin.Context) {}
func (h *ModuleHandler) UpdatePvalue(c *gin.Context) {}
```

- [ ] **Step 4: 注册路由**

在 `router.go` 中注册上述端点。

- [ ] **Step 5: 构建验证**

Run: `go build ./... && go vet ./...`

---

### Task 6: AgentKnowledge 向量化事件系统

**Files:**
- Create: `internal/service/knowledge_embedding.go`
- Modify: `api/handler/knowledge/agent_knowledge.go`
- Modify: `internal/job/jobs/jobs.go`

**说明:** 实现 Java 版 ApplicationEvent 模式的异步向量化事件。

- [ ] **Step 1: 创建 EmbeddingService**

```go
// internal/service/knowledge_embedding.go
type EmbeddingService struct {
    repo    repository.AgentKnowledgeRepository
    logger  *zap.Logger
}

func (s *EmbeddingService) TriggerEmbedding(ctx context.Context, knowledgeID string) error
func (s *EmbeddingService) TriggerDeletion(ctx context.Context, knowledgeID string) error
func (s *EmbeddingService) ProcessPendingBatch(ctx context.Context, batchSize int) error
```

- [ ] **Step 2: 在 handler 中触发事件**

`AgentKnowledgeHandler.Create` 创建后调用 `TriggerEmbedding`。`Delete` 时调用 `TriggerDeletion`。

- [ ] **Step 3: 添加定时清理 Job**

在 `jobs.go` 中添加 `KnowledgeResourceCleanupJob`，每小时清理孤儿向量/文件资源。

- [ ] **Step 4: 构建验证**

Run: `go build ./... && go vet ./...`

---

### Task 7: Langfuse 链路追踪接入

**Files:**
- Create: `internal/service/tracing/langfuse.go`
- Modify: `api/handler/chat/graph.go`
- Modify: `agent/workflows/nl2sql/nodes/*.go`（各节点添加 span）

**说明:** 接入 OpenTelemetry 实现 LLM 调用链路追踪。

- [ ] **Step 1: 创建 TracingService**

```go
// internal/service/tracing/langfuse.go
type TracingService struct {
    tracer trace.Tracer
}

func (s *TracingService) StartGraphSpan(ctx context.Context, threadID, query string) (context.Context, trace.Span)
func (s *TracingService) StartNodeSpan(ctx context.Context, nodeName string) (context.Context, trace.Span)
func (s *TracingService) RecordTokens(span trace.Span, prompt, completion int)
func (s *TracingService) EndSpan(span trace.Span, output string)
```

- [ ] **Step 2: 接入 GraphHandler**

在 `StreamSearch` 开始处调用 `StartGraphSpan`，结束处调用 `EndSpan`。

- [ ] **Step 3: 各节点添加 span**

在每个节点的 `Execute` 方法中调用 `StartNodeSpan`/`EndSpan`，记录 token 用量。

- [ ] **Step 4: 构建验证**

Run: `go build ./... && go vet ./...`

---

### Task 8: PgVectorStore 向量存储集成

**Files:**
- Create: `internal/dao/vectorstore/pgvector.go`
- Modify: `agent/knowledge/retriever.go`
- Modify: `internal/config/vectorstore.go`

**说明:** 实现 Java 版 PgVectorStore，512 维/余弦距离/HNSW 索引。

- [ ] **Step 1: 创建 VectorStoreConfig**

```go
// internal/config/vectorstore.go
type VectorStoreConfig struct {
    Dimensions         int     `mapstructure:"dimensions"`
    SimilarityThreshold float64 `mapstructure:"similarity_threshold"`
    TableTopkLimit     int     `mapstructure:"table_topk_limit"`
    DefaultTopkLimit   int     `mapstructure:"default_topk_limit"`
    BatchDelTopkLimit  int     `mapstructure:"batch_del_topk_limit"`
    EnableHybridSearch bool    `mapstructure:"enable_hybrid_search"`
}
```

- [ ] **Step 2: 实现 PgVectorStore**

```go
// internal/dao/vectorstore/pgvector.go
type PgVectorStore struct {
    db     *gorm.DB
    cfg    VectorStoreConfig
    embed  knowledge.Embedder
}

func (s *PgVectorStore) Search(ctx context.Context, query string, topK int) ([]Document, error)
func (s *PgVectorStore) AddDocuments(ctx context.Context, docs []Document) error
func (s *PgVectorStore) DeleteDocuments(ctx context.Context, ids []string) error
func (s *PgVectorStore) DeleteByFilter(ctx context.Context, filter map[string]string) error
```

使用 `pgvector` GORM 扩展和 `<=>` 余弦距离运算符。

- [ ] **Step 3: 接入 Retriever**

更新 `agent/knowledge/retriever.go` 使用 PgVectorStore 替代空实现。

- [ ] **Step 4: 构建验证**

Run: `go build ./... && go vet ./...`

---

### Task 9: 文本分割策略

**Files:**
- Create: `agent/knowledge/splitter/strategies.go`
- Modify: `agent/knowledge/splitter.go`

**说明:** 实现 Java 版 5 种文本分割器。

- [ ] **Step 1: 定义接口和策略**

```go
type TextSplitter interface {
    Split(text string) []Chunk
}

type Chunk struct {
    Content string
    Index   int
}
```

实现 5 种策略：
- `TokenSplitter` — 按 token 数切分（chunkSize=1000, minChunkSizeChars=400）
- `RecursiveSplitter` — 递归按分隔符切分（separators=["\n\n", "\n", " ", ""]）
- `SentenceSplitter` — 按句子切分
- `SemanticSplitter` — 基于 embedding 相似度切分
- `ParagraphSplitter` — 按段落切分

- [ ] **Step 2: 构建验证**

Run: `go build ./... && go vet ./...`

---

### Task 10: Python Docker 沙箱

**Files:**
- Create: `internal/service/code/docker_executor.go`
- Modify: `internal/service/code/executor.go`
- Modify: `internal/config/code.go`

**说明:** 实现 Java 版 Docker 容器隔离的 Python 代码执行。

- [ ] **Step 1: 创建 CodeExecutorConfig**

```go
type CodeExecutorConfig struct {
    Type                string `mapstructure:"type"` // docker/local/simulation
    ImageName           string `mapstructure:"image_name"`
    ContainerPrefix     string `mapstructure:"container_prefix"`
    LimitMemoryMB       int    `mapstructure:"limit_memory_mb"`
    CPUCores            int    `mapstructure:"cpu_cores"`
    CodeTimeout         int    `mapstructure:"code_timeout_seconds"`
    ContainerTimeout    int    `mapstructure:"container_timeout_seconds"`
    NetworkMode         string `mapstructure:"network_mode"`
}
```

- [ ] **Step 2: 实现 DockerExecutor**

使用 Docker Engine API (`github.com/docker/docker/client`) 创建容器、挂载代码文件、执行、获取结果、清理容器。

- [ ] **Step 3: 构建验证**

Run: `go build ./... && go vet ./...`

---

### Task 11: Graph Checkpoint/Saver

**Files:**
- Create: `internal/dao/checkpoint/checkpoint.go`
- Modify: `agent/workflows/nl2sql/graph.go`

**说明:** 实现图状态持久化，支持 interrupt/resume。

- [ ] **Step 1: 实现 CheckpointStore**

```go
type CheckpointStore struct {
    db *gorm.DB
}

func (s *CheckpointStore) Save(threadID string, state graph.State) error
func (s *CheckpointStore) Load(threadID string) (graph.State, error)
func (s *CheckpointStore) Delete(threadID string) error
```

- [ ] **Step 2: 配置 interruptBefore**

在 `graph.go` 的 `Build()` 中使用 `graph.WithInterruptBefore(NodeHumanFeedback)` 配置中断点。

- [ ] **Step 3: 构建验证**

Run: `go build ./... && go vet ./...`

---

### Task 12: 实体字段补全

**Files:**
- Modify: `internal/model/data_entity.go`
- Modify: `internal/model/privilege_entity.go`
- Modify: `internal/model/rag_entity.go`

**说明:** 补全 Java 有但 Go 缺失的实体字段。

- [ ] **Step 1: AgentKnowledge 添加 isResourceCleaned**

```go
IsResourceCleaned int `gorm:"column:is_resource_cleaned;default:0" json:"isResourceCleaned"`
```

- [ ] **Step 2: PrivilegeUser 添加缺失字段**

添加 `AclTimestamp`, `PwdFtime`, `PwdInit` 字段。

- [ ] **Step 3: RagFileInfo 添加缺失字段**

添加 `PdfType`, `PageTopMargin`, `PagesPerDocument`, `TextSplitter` 字段。

- [ ] **Step 4: 构建验证**

Run: `go build ./... && go vet ./...`

---

### Task 13: 动态模型代理/热切换

**Files:**
- Create: `internal/service/model/registry.go`
- Create: `internal/service/model/proxy.go`
- Modify: `agent/workflows/nl2sql/graph.go`

**说明:** 实现 Java 版 AiModelRegistry + 动态 EmbeddingModel 代理。

- [ ] **Step 1: 创建 ModelRegistry**

```go
type ModelRegistry struct {
    mu     sync.RWMutex
    models map[string]model.Model
    active map[string]string // type → modelID
}

func (r *ModelRegistry) Register(id string, m model.Model)
func (r *ModelRegistry) GetActive(modelType string) model.Model
func (r *ModelRegistry) SetActive(modelType, id string)
func (r *ModelRegistry) SwitchModel(modelType, newID string) error
```

- [ ] **Step 2: 创建 ModelProxy**

```go
type ModelProxy struct {
    registry *ModelRegistry
    modelType string
}

func (p *ModelProxy) GenerateContent(ctx, req) (<-chan *model.Response, error) {
    m := p.registry.GetActive(p.modelType)
    return m.GenerateContent(ctx, req)
}
```

- [ ] **Step 3: 构建验证**

Run: `go build ./... && go vet ./...`

---

### Task 14: StreamContext 管理

**Files:**
- Create: `internal/service/stream/context.go`
- Modify: `api/handler/chat/graph.go`

**说明:** 实现 Java 版 `StreamContext` 管理，支持并发流管理和清理。

- [ ] **Step 1: 创建 StreamContextManager**

```go
type StreamContextManager struct {
    mu       sync.RWMutex
    contexts map[string]*StreamContext
}

type StreamContext struct {
    ThreadID       string
    AgentID        string
    StartTime      time.Time
    CollectedOutput strings.Builder
    TextType       string
    Span           trace.Span
    CancelFunc     context.CancelFunc
}
```

- [ ] **Step 2: 接入 GraphHandler**

在 `StreamSearch` 开始时注册 StreamContext，结束时清理。客户端断连时自动清理。

- [ ] **Step 3: 构建验证**

Run: `go build ./... && go vet ./...`

---

## 执行顺序

```
Phase 1 — 快速补全 (无需外部依赖)
  Task 1: NL2SQL 状态键 + 条件边     (1h)
  Task 2: 配置外部化                 (2h)
  Task 4: HTML 报告二进制下载         (0.5h)
  Task 5: 缺失端点补全               (1h)
  Task 12: 实体字段补全              (0.5h)

Phase 2 — 架构组件
  Task 3: MultiTurnContextManager    (2h)
  Task 6: 向量化事件系统             (2h)
  Task 7: Langfuse 链路追踪          (3h)
  Task 14: StreamContext 管理        (1h)
  Task 11: Graph Checkpoint/Saver    (2h)

Phase 3 — 外部集成
  Task 8: PgVectorStore 向量存储     (4h)
  Task 9: 文本分割策略               (3h)
  Task 10: Python Docker 沙箱        (4h)

Phase 4 — 优化项
  Task 13: 动态模型代理              (2h)
```

**总预估: ~28h**
