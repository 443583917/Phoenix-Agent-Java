# Phoenix Go 重写 — 设计文档

> 将 Phoenix-Agent-Java（Spring AI Alibaba）后端全部迁移到 Go，前端 API 契约保持不变。
> Java 版本没有对应映射代码的模块先保留空目录，后续填充。

## 1. 概述

### 1.1 目标

| 维度 | 说明 |
|:---|:---|
| **语言** | Java 21 → Go 1.22+ |
| **框架** | Spring Boot 4.0 / Spring AI Alibaba → Gin + tRPC-Agent-Go |
| **范围** | 全部 12 个 Maven 模块 |
| **兼容** | REST + SSE 接口契约不变，前端零改动 |
| **策略** | 分层渐进式（Phase 1→6 串行，每层交付可验证） |

### 1.2 Java → Go 架构映射

```
Spring Boot App (8066)     →  cmd/api/main.go (Gin, 8066)
Spring WebFlux (Reactive)  →  Gin + goroutine (天然并发，无需 reactive)
MyBatis-Flex ORM           →  GORM + PostgreSQL driver
Sa-Token RBAC              →  golang-jwt + OAuth2 + Casbin
JetCache (Redis+Local)     →  go-redis + bigcache
PgVector 向量存储           →  Milvus (向量) + PostgreSQL (元数据)
Spring AI Graph            →  tRPC-Agent-Go workflows
ReactAgent / Tools         →  tRPC-Agent-Go agents / tools
MCP Server                 →  tRPC-Agent-Go MCP tool
RedisSaver Checkpoint      →  Redis hash (agent/runner checkpoint)
OpenTelemetry → Langfuse   →  otel-go + langfuse-go
Docker Python Executor     →  Go Docker SDK
```

### 1.3 技术栈最终选型

| 组件 | Go 库 | 对标 Java |
|:---|:---|:---|
| HTTP | `gin-gonic/gin` | WebFlux Router |
| ORM | `gorm.io/gorm` + `gorm.io/driver/postgres` | MyBatis-Flex |
| 认证 | `golang-jwt/jwt/v5` + OAuth2 | Sa-Token |
| 权限 | `casbin/casbin/v2` | Sa-Token RBAC |
| 缓存(L2) | `redis/go-redis/v9` | JetCache Redis |
| 缓存(L1) | `allegro/bigcache/v3` | JetCache Local |
| 向量 | `milvus-io/milvus-sdk-go/v2` | PgVector |
| 消息队列 | `rabbitmq/amqp091-go` + Redis Pub/Sub | Redis Pub/Sub |
| 可观测性 | `go.opentelemetry.io/otel` + `langfuse-go` | OTel → Langfuse |
| 配置 | `spf13/viper` | application.yml |
| 日志 | `uber-go/zap` | Logback |
| DB迁移 | `golang-migrate/migrate/v4` | SQL 脚本 |
| Python执行 | `github.com/docker/docker/client` | docker-java |
| ID生成 | `sony/sonyflake` | MyBatis-Flex ID |
| 定时任务 | `robfig/cron/v3` | Spring @Scheduled |
| Agent框架 | `trpc.group/trpc-go/trpc-agent-go` | Spring AI Alibaba Graph |
| 测试 | `stretchr/testify` + `DATA-DOG/go-sqlmock` | JUnit + Mockito |

---

## 2. 项目目录结构

```
godata/
├── cmd/                        # 多服务入口（每个服务独立进程）
│   ├── api/main.go             # HTTP API 服务
│   ├── rpc/main.go             # gRPC 服务
│   ├── agent/main.go           # Agent 服务（可独立部署）
│   └── job/main.go             # 定时任务服务
│
├── api/                        # API 层（Gin）
│   ├── router.go               # 总路由注册
│   ├── handler/                # 控制器层（按业务域拆分）
│   │   ├── agent/
│   │   ├── datasource/
│   │   ├── chat/
│   │   ├── knowledge/
│   │   ├── modelconfig/
│   │   ├── prompt/
│   │   ├── semanticmodel/
│   │   ├── privilege/
│   │   ├── platform/
│   │   ├── rag/
│   │   ├── kg/
│   │   └── common/
│   └── middleware/             # API 中间件
│       ├── auth.go
│       ├── rbac.go
│       ├── cors.go
│       ├── logger.go
│       ├── recovery.go
│       ├── tracing.go
│       └── ratelimit.go
│
├── rpc/                        # gRPC 层（Phase 6 实现，前期为空骨架）
│   ├── proto/
│   ├── server/
│   └── client/
│
├── agent/                      # AI Agent 引擎（核心）
│   ├── runtime/                # Agent 生命周期、注册、初始化
│   ├── agents/                 # 各类 Agent（React / Workflow / Assistant）
│   ├── tools/                  # 工具层（按工具类型拆分）
│   │   ├── function/           # FunctionTool
│   │   ├── rpc/                # RPC 工具
│   │   ├── datasource/         # 数据源工具
│   │   ├── privilege/          # 权限工具
│   │   ├── web/                # Web 搜索工具
│   │   └── external/           # 外部系统工具
│   ├── workflows/              # 工作流（按工作流拆分）
│   │   ├── nl2sql/
│   │   │   ├── nodes/
│   │   │   └── graph.go
│   │   ├── rag/
│   │   │   ├── nodes/
│   │   │   └── graph.go
│   │   └── kg/
│   │       ├── nodes/
│   │       └── graph.go
│   ├── memory/                 # 记忆系统（短期/长期/画像）
│   ├── knowledge/              # RAG 检索、Embedding、Splitter
│   └── runner/                 # 对话管道、SSE、HITL
│
├── internal/                   # 业务核心层（DDD）
│   ├── domain/                 # 领域层（纯业务规则）
│   │   ├── agent/
│   │   ├── datasource/
│   │   ├── chat/
│   │   ├── knowledge/
│   │   ├── modelconfig/
│   │   ├── prompt/
│   │   ├── semanticmodel/
│   │   ├── privilege/
│   │   ├── platform/
│   │   ├── rag/
│   │   └── common/
│   ├── usecase/                # 用例层（业务编排）
│   │   ├── agent_usecase.go
│   │   ├── rag_usecase.go
│   │   └── datasource_usecase.go
│   ├── service/                # 应用服务层（对接 handler）
│   ├── repository/             # 仓储层接口定义
│   │   ├── agent_repo.go
│   │   ├── datasource_repo.go
│   │   └── rag_repo.go
│   ├── dao/                    # 数据访问层（仓储实现）
│   │   ├── db/
│   │   ├── cache/
│   │   ├── queue/
│   │   └── external/
│   ├── event/                  # 领域事件
│   ├── job/                    # 定时任务
│   ├── model/                  # 数据模型（entity/dto/vo）
│   └── config/                 # 配置结构体（Viper 映射）
│
├── infra/                      # 基础设施层
│   ├── logger/
│   ├── config/
│   ├── monitoring/
│   ├── queue/
│   ├── cache/
│   ├── lock/
│   ├── id/
│   └── utils/
│
├── configs/                    # YAML 配置文件（按服务拆分）
│   ├── api/
│   ├── rpc/
│   ├── agent/
│   ├── job/
│   ├── db.yaml
│   ├── redis.yaml
│   ├── milvus.yaml
│   ├── rabbitmq.yaml
│   └── monitor.yaml
│
├── migrations/                 # 数据库迁移
├── scripts/
├── docs/
├── Dockerfile
├── docker-compose.yaml
├── Makefile
├── go.mod
└── go.sum
```

### 2.1 分层架构说明

```
┌─────────────────────────────────────────────┐
│  api/handler/   ←  HTTP 请求/响应处理         │
├─────────────────────────────────────────────┤
│  internal/service/  ← 应用服务（对接 handler）  │
├─────────────────────────────────────────────┤
│  internal/usecase/  ← 用例编排（跨领域协调）     │
├─────────────────────────────────────────────┤
│  internal/domain/   ← 领域模型 + 业务规则      │
├─────────────────────────────────────────────┤
│  internal/repository/  ← 仓储接口（依赖反转）  │
│  internal/dao/          ← 仓储实现（GORM等）   │
├─────────────────────────────────────────────┤
│  agent/   ← AI Agent 引擎（独立于业务层）       │
├─────────────────────────────────────────────┤
│  infra/   ← 基础设施（日志/缓存/队列/监控）     │
└─────────────────────────────────────────────┘
```

---

## 3. Java 模块 → Go 包映射

```
phoenix-common-api/core/rest     →  infra/ + internal/model/common_* + internal/domain/common/
phoenix-privilege-api/core/rest  →  internal/domain/privilege/ + api/handler/privilege/ + internal/dao/db/privilege_repo.go
phoenix-platform-api/core/rest   →  internal/domain/platform/ + api/handler/platform/ + internal/dao/db/platform_repo.go
phoenix-agent-api/core/rest      →  agent/ + internal/domain/agent/ + api/handler/agent/
phoenix-data-api/core/rest       →  agent/workflows/nl2sql/ + internal/domain/datasource/ + internal/domain/chat/ + ...
phoenix-rag-api/core/rest        →  internal/domain/rag/ + agent/knowledge/ + agent/workflows/rag/ + api/handler/rag/
phoenix-kg-api/core/rest         →  internal/domain/kg/ + agent/workflows/kg/ + api/handler/kg/
phoenix-tool                     →  infra/utils/
phoenix-codegen                  →  (用不着，Go 没有代码生成器需求)
phoenix-flink                    →  internal/dao/external/flink.go（降级为外部集成）
phoenix-admin-manager            →  cmd/api/main.go（统一入口）
```

---

## 4. 数据库策略

### 4.1 零 Schema 变更

复用现有 `sql/all_schema.sql` + `sql/all_data.sql`，**不改动任何表结构**。GORM 通过 tag 映射现有表。

```go
// internal/model/agent_entity.go
type Agent struct {
    ID           uint64     `gorm:"primaryKey;column:id" json:"id"`
    Type         string     `gorm:"column:type" json:"type"`
    CategoryID   *uint64    `gorm:"column:category_id" json:"categoryId"`
    SN           string     `gorm:"column:sn" json:"sn"`
    Name         string     `gorm:"column:name" json:"name"`
    Description  string     `gorm:"column:description" json:"description"`
    Status       int        `gorm:"column:status" json:"status"`
    ApiKey       string     `gorm:"column:api_key" json:"apiKey"`
    CreatedAt    time.Time  `gorm:"column:create_time" json:"createTime"`
    UpdatedAt    time.Time  `gorm:"column:update_time" json:"updateTime"`
}

func (Agent) TableName() string { return "tbl_data_agent" }
```

### 4.2 新增表使用 golang-migrate

```
migrations/
├── 000001_create_sessions_table.up.sql
├── 000001_create_sessions_table.down.sql
└── ...
```

### 4.3 Repository 模式（依赖反转）

```go
// internal/repository/agent_repo.go — 接口
type AgentRepository interface {
    FindByID(ctx context.Context, id uint64) (*model.Agent, error)
    FindBySN(ctx context.Context, sn string) (*model.Agent, error)
    List(ctx context.Context, query model.AgentPageQuery) ([]*model.Agent, int64, error)
    Create(ctx context.Context, agent *model.Agent) error
    Update(ctx context.Context, agent *model.Agent) error
    Delete(ctx context.Context, id uint64) error
}

// internal/dao/db/agent_repo.go — GORM 实现
type agentRepo struct {
    db *gorm.DB
}
func NewAgentRepository(db *gorm.DB) repository.AgentRepository {
    return &agentRepo{db: db}
}
```

---

## 5. 统一 API 规范

### 5.1 响应格式（对标 Java R 类）

```go
// infra/response/response.go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    TraceID string      `json:"traceId,omitempty"`
}

type PageResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
    Total   int64       `json:"total"`
    Page    int         `json:"page"`
    Size    int         `json:"size"`
    TraceID string      `json:"traceId,omitempty"`
}
```

### 5.2 路由注册（router.go）

```go
func SetupRouter() *gin.Engine {
    r := gin.New()
    r.Use(middleware.Recovery())
    r.Use(middleware.Logger())
    r.Use(middleware.Tracing())
    r.Use(middleware.CORS())

    // 健康检查
    r.GET("/echo", handler.Echo)

    api := r.Group("/api")
    api.Use(middleware.Auth())      // JWT 认证
    api.Use(middleware.RBAC())      // Casbin 权限

    {
        agentGroup := api.Group("/agent")
        agent.RegisterRoutes(agentGroup)

        datasourceGroup := api.Group("/datasource")
        datasource.RegisterRoutes(datasourceGroup)

        chatGroup := api.Group("")
        chat.RegisterRoutes(chatGroup)

        // ...其他模块
    }

    privilegeGroup := r.Group("/api/privilege")
    privilege.RegisterRoutes(privilegeGroup)

    platformGroup := r.Group("/platform")
    platformGroup.Use(middleware.Auth())
    platform.RegisterRoutes(platformGroup)

    return r
}
```

### 5.3 SSE 兼容（对标 WebFlux SSE）

```go
// infra/sse/sse.go
func Stream(c *gin.Context, events <-chan Event) {
    c.Writer.Header().Set("Content-Type", "text/event-stream")
    c.Writer.Header().Set("Cache-Control", "no-cache")
    c.Writer.Header().Set("Connection", "keep-alive")
    c.Writer.WriteHeader(http.StatusOK)

    flusher, _ := c.Writer.(http.Flusher)
    for event := range events {
        data, _ := json.Marshal(event)
        fmt.Fprintf(c.Writer, "data: %s\n\n", data)
        flusher.Flush()
        if c.Request.Context().Err() != nil {
            return
        }
    }
}
```

---

## 6. 分阶段实施计划

### Phase 1 — 基础设施层

**目标**：项目骨架跑通，Gin 启动并响应 `/echo`

| 交付物 | 内容 |
|:---|:---|
| `go.mod` | `module github.com/phoenix-agent-go` |
| `cmd/api/main.go` | Gin 启动，Viper 加载配置，注册中间件 |
| `configs/` | api/, rpc/, agent/, job/ 子目录 + db.yaml, redis.yaml, milvus.yaml, rabbitmq.yaml, monitor.yaml |
| `infra/logger/` | Zap 封装，支持日志级别、文件轮转 |
| `infra/response/` | 统一 Response / PageResponse |
| `infra/errcode/` | 统一错误码定义 |
| `infra/jwt/` | JWT 生成/验证 |
| `infra/pagination/` | 分页参数解析 |
| `infra/sse/` | SSE 流式输出 |
| `infra/validate/` | 参数校验封装 |
| `infra/config/` | Viper 配置加载入口 |
| `infra/monitoring/` | OpenTelemetry 初始化 |
| `infra/queue/` | RabbitMQ + Redis 连接管理 |
| `infra/cache/` | go-redis + bigcache 初始化 |
| `infra/lock/` | Redis 分布式锁 |
| `infra/id/` | Sonyflake ID 生成器 |
| `internal/config/` | Viper 配置结构体映射（db/redis/milvus/rabbitmq/agent/monitor） |
| `api/middleware/` | cors, logger, recovery, tracing, auth (骨架) |
| `api/router.go` | 空路由骨架 |
| `Dockerfile` | 多阶段构建 |
| `docker-compose.yaml` | PostgreSQL + Redis + Milvus + RabbitMQ |
| `Makefile` | build, run, test, lint, migrate |

### Phase 2 — 权限认证模块（phoenix-privilege）

**目标**：完整 RBAC，前端登录流程可用

| 交付物 | 映射 Java |
|:---|:---|
| `internal/model/privilege_entity.go` | PrivilegeUser, PrivilegeRole, PrivilegeModule 等 12 张表 |
| `internal/model/privilege_dto.go` | LoginDTO, UserPageQuery 等 |
| `internal/domain/privilege/` | 领域规则（密码加密、角色校验） |
| `internal/repository/privilege_repo.go` | 仓储接口 |
| `internal/dao/db/privilege_repo.go` | GORM 实现 |
| `internal/dao/cache/privilege_cache.go` | Redis 权限缓存 + BigCache L1 |
| `internal/usecase/privilege_usecase.go` | 用例编排 |
| `internal/service/privilege_service.go` | 应用服务（对接 handler） |
| `api/handler/privilege/` | 12 个 handler，对齐所有 Java Controller |
| `api/middleware/auth.go` | JWT + OAuth2 认证中间件 |
| `api/middleware/rbac.go` | Casbin RBAC 权限中间件 |

**API 端点对齐：**

| 原 Java Controller | 端点 | Go Handler |
|:---|:---|:---|
| LoginController | `/api/privilege/auth/*` | `handler/privilege/auth.go` |
| PrivilegeUserController | `/api/privilege/user` | `handler/privilege/user.go` |
| PrivilegeRoleController | `/api/privilege/role` | `handler/privilege/role.go` |
| PrivilegeUserRoleController | `/api/privilege/user-role` | `handler/privilege/user_role.go` |
| PrivilegeModuleController | `/api/privilege/module` | `handler/privilege/module.go` |
| PrivilegeAclController | `/api/privilege/acl` | `handler/privilege/acl.go` |
| PrivilegeDepartmentController | `/api/privilege/department` | `handler/privilege/department.go` |
| PrivilegeCompanyController | `/api/privilege/company` | `handler/privilege/company.go` |
| PrivilegeEmployeeController | `/api/privilege/employee` | `handler/privilege/employee.go` |
| PrivilegeDictionaryController | `/api/privilege/dictionary` | `handler/privilege/dictionary.go` |
| PrivilegePvalueController | `/api/privilege/pvalue` | `handler/privilege/pvalue.go` |
| PrivilegeLoginLogController | `/api/privilege/login-log` | `handler/privilege/login_log.go` |

### Phase 3 — 平台管理模块（phoenix-platform + phoenix-common）

**目标**：多租户、三方集成（钉钉/飞书/企微）

| 交付物 | 映射 Java |
|:---|:---|
| `internal/model/platform_entity.go` | GroupInfo, AccountInfo, TenantInfo 等 |
| `internal/model/common_entity.go` | PlatformInfo |
| `internal/domain/platform/` | 平台领域规则 |
| `internal/dao/db/platform_repo.go` | GORM CRUD |
| `internal/dao/external/dingtalk.go` | 钉钉 SDK 封装 |
| `internal/dao/external/feishu.go` | 飞书 SDK 封装 |
| `internal/dao/external/wecom.go` | 企业微信 SDK 封装 |
| `internal/usecase/platform_usecase.go` | 用例编排 |
| `internal/service/platform_service.go` | 应用服务 |
| `api/handler/platform/` | GroupAgentInfo, AccountLogin 等 handler |
| `api/handler/common/` | PlatformInfo, PlatformSync handler |

### Phase 4 — Agent 框架（phoenix-agent）

**目标**：React Agent + 会话管理 + 记忆 + MCP，核心 AI 能力

| 交付物 | 映射 Java |
|:---|:---|
| `agent/runtime/manager.go` | AgentManager |
| `agent/runtime/registry.go` | Agent 注册表 |
| `agent/agents/react_agent.go` | ReactAgent |
| `agent/agents/workflow_agent.go` | Workflow Graph Agent |
| `agent/tools/function/` | FunctionTool 注册（对标 @Tool 注解） |
| `agent/tools/datasource/` | SQL 查询工具 |
| `agent/tools/rpc/` | RPC 调用工具 |
| `agent/tools/web/` | Web 搜索工具 |
| `agent/tools/privilege/` | 权限检查工具 |
| `agent/tools/external/mcp.go` | MCP 协议工具 |
| `agent/memory/short_term.go` | 对话窗口管理（Redis checkpoint） |
| `agent/memory/long_term.go` | Milvus 向量检索记忆 |
| `agent/memory/profile.go` | 用户画像 |
| `agent/knowledge/retriever.go` | 混合检索（向量 + 关键词 + RRF 融合） |
| `agent/runner/runner.go` | 对话执行管道 |
| `agent/runner/sse.go` | SSE 流式输出 |
| `agent/runner/hitl.go` | Human-in-the-Loop |
| `internal/model/agent_entity.go` | UserAgentInfo, UserMemoryInfo, CombinedStore |
| `internal/domain/agent/` | Agent 领域规则 |
| `internal/dao/db/agent_repo.go` | GORM CRUD |
| `internal/dao/external/milvus.go` | Milvus 操作封装 |
| `internal/usecase/agent_usecase.go` | Agent 用例编排 |
| `internal/service/agent_service.go` | Agent 应用服务 |
| `api/handler/agent/` | ReactAgentController, HarnessController |

### Phase 5 — NL2SQL 数据引擎（phoenix-data）

**目标**：StateGraph 工作流，自然语言转 SQL/Python 分析，最复杂的模块

| 交付物 | 映射 Java |
|:---|:---|
| `agent/workflows/nl2sql/graph.go` | DataAgentConfiguration StateGraph |
| `agent/workflows/nl2sql/nodes/intent.go` | IntentRecognitionNode |
| `agent/workflows/nl2sql/nodes/evidence.go` | EvidenceRecallNode |
| `agent/workflows/nl2sql/nodes/schema.go` | SchemaRecallNode + QueryEnhanceNode + TableRelationNode |
| `agent/workflows/nl2sql/nodes/planner.go` | PlannerNode + FeasibilityAssessmentNode + PlanExecutorNode |
| `agent/workflows/nl2sql/nodes/sql_gen.go` | SqlGenerateNode + SemanticConsistencyNode + SqlExecuteNode |
| `agent/workflows/nl2sql/nodes/python_exec.go` | PythonGenerateNode + PythonExecuteNode + PythonAnalyzeNode |
| `agent/workflows/nl2sql/nodes/report.go` | ReportGeneratorNode |
| `agent/workflows/nl2sql/nodes/checkpoint.go` | Redis 状态检查点 |
| `internal/domain/datasource/` | 数据源连接管理 |
| `internal/dao/db/ddl_repo.go` | 多数据库 DDL 元数据提取 (MySQL/PG/Oracle/MSSQL/H2/Hive/DM) |
| `internal/dao/external/docker_exec.go` | Docker Python 执行器 |
| `internal/domain/chat/` | 对话会话领域 |
| `internal/model/data_entity.go` | Agent, Datasource, ChatSession, ChatMessage 等 |
| `api/handler/chat/` | ChatController — SSE 对话接口 |
| `api/handler/datasource/` | DatasourceController |
| `api/handler/knowledge/` | AgentKnowledgeController, BusinessKnowledgeController |
| `api/handler/modelconfig/` | ModelConfigController |
| `api/handler/prompt/` | PromptConfigController |
| `api/handler/semanticmodel/` | SemanticModelController |

### Phase 6 — 扩展模块 + 部署

**目标**：RAG, KG, Flink, gRPC, 全量测试, 生产部署

| 交付物 | 映射 Java |
|:---|:---|
| `internal/domain/rag/` | RAG 文件管理 |
| `agent/workflows/rag/` | RAG 工作流 |
| `api/handler/rag/` | RagFileInfoController, RagCategoryController |
| `internal/domain/kg/` | 知识图谱 |
| `agent/workflows/kg/` | KG 工作流 |
| `api/handler/kg/` | KG handler |
| `internal/dao/external/flink.go` | Flink 集成 |
| `cmd/rpc/main.go` | gRPC 服务启动 |
| `cmd/agent/main.go` | Agent 独立服务启动 |
| `cmd/job/main.go` | 定时任务服务启动 |
| `rpc/proto/` | agent.proto, privilege.proto, data.proto |
| `rpc/server/` + `rpc/client/` | gRPC 服务端 + 客户端 |

---

## 7. 关键设计决策

### 7.1 DDD 分层：domain → usecase → repository

```
api/handler/          ← HTTP 适配（Gin ctx → DTO → usecase）
     │
internal/service/     ← 应用服务（薄层，对接 handler，格式转换）
     │
internal/usecase/     ← 用例编排（跨 domain 协调、事务边界、事件发布）
     │
internal/domain/      ← 领域核心（纯业务规则、不依赖框架）
     │
internal/repository/  ← 接口定义（依赖反转，domain 依赖接口而非实现）
     │
internal/dao/         ← 仓储实现（GORM, Redis, Milvus, RabbitMQ）
```

**为什么不直接用 service 调 dao？** 加上 usecase 和 repository 两层之后：
- `repository` 接口让 domain 不依赖 GORM，可单独测试
- `usecase` 处理跨领域编排，service 只管格式转换，职责清晰
- 对标 Java 中 Service 的 `@Transactional` 编排逻辑，usecase 是事务边界

### 7.2 并发模型：Handler → Usecase → Repository 全链路 ctx 传递

```go
func (h *AgentHandler) GetAgent(c *gin.Context) {
    ctx := c.Request.Context()
    agent, err := h.usecase.GetAgent(ctx, id)
}

func (u *AgentUsecase) GetAgent(ctx context.Context, id uint64) (*model.Agent, error) {
    // 先查 BigCache L1
    if cached := u.cache.Get(ctx, id); cached != nil { return cached, nil }
    // 再查 Repository（GORM）
    agent, err := u.repo.FindByID(ctx, id)
    // 写回 L1
    u.cache.Set(ctx, id, agent)
    return agent, err
}
```

Gin 天然每请求一个 goroutine，无需 WebFlux reactive。

### 7.3 事务管理

事务边界放在 usecase 层，通过 `db.Transaction()` 显式管理：

```go
func (u *AgentUsecase) CreateAgent(ctx context.Context, dto CreateAgentDTO) (*model.Agent, error) {
    var agent *model.Agent
    err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        agent = dto.ToEntity()
        if err := tx.Create(agent).Error; err != nil { return err }
        ds := dto.ToDefaultDatasource(agent.ID)
        return tx.Create(ds).Error
    })
    return agent, err
}
```

### 7.4 错误处理

统一错误码，对标 Java 的 BizException / GlobalExceptionHandler：

```go
// infra/errcode/errcode.go
type ErrCode struct {
    Code int
    Msg  string
}

var (
    Success         = ErrCode{0, "success"}
    Unauthorized    = ErrCode{401, "未认证"}
    Forbidden       = ErrCode{403, "无权限"}
    NotFound        = ErrCode{404, "资源不存在"}
    InternalError   = ErrCode{500, "服务器内部错误"}
    InvalidParams   = ErrCode{1001, "参数校验失败"}
    AgentOffline    = ErrCode{2001, "智能体已下线"}
    DatasourceError = ErrCode{3001, "数据源连接失败"}
)

// middleware/recovery.go 统一捕获 panic + 转换为 Response
```

### 7.5 配置管理

YAML 按服务拆分目录，共享配置平级放置，Viper 合并加载：

```yaml
# configs/db.yaml (所有服务共享)
database:
  host: "127.0.0.1"
  port: 5432
  user: "phoenix"
  password: "phoenix"
  name: "phoenix"
  sslmode: "disable"
  max_open_conns: 25
  max_idle_conns: 10

# configs/redis.yaml (所有服务共享)
redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0

# configs/api/app.yaml (API 服务专用)
server:
  port: 8066
  mode: debug

# configs/agent/app.yaml (Agent 服务专用)
agent:
  model:
    provider: "deepseek"
    model: "deepseek-chat"
    api_key: "${AI_API_KEY}"
    base_url: "https://api.deepseek.com"
  stream: true
  max_tokens: 4096
```

### 7.6 Agent 层独立于 business 层

```
agent/                ← 纯 AI 引擎，依赖 tRPC-Agent-Go
    可独立编译、部署为 cmd/agent
    不 import internal/ 的任何包

internal/             ← 业务逻辑，通过接口调用 agent 层
    通过 agent 层的接口（AgentManager, ToolRegistry）
    注入给 ai handler
```

依赖方向：`api → internal → agent ← infra`（internal 和 agent 都依赖 infra，但互相不依赖）

---

## 8. 部署架构

```
docker-compose.yaml
├── phoenix-api (Gin, :8066)    — 来自 cmd/api
├── phoenix-agent (:8090)       — 来自 cmd/agent（可选独立部署）
├── postgres (:5432) + pgvector
├── redis (:6379)
├── milvus-standalone (:19530, :9091)
├── rabbitmq (:5672, :15672)
└── (optional) python-executor (Docker-in-Docker)
```

### Dockerfile（多阶段构建）

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /phoenix-api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata python3 py3-pip
COPY --from=builder /phoenix-api /usr/local/bin/
COPY configs/ /etc/phoenix/configs/
EXPOSE 8066
CMD ["phoenix-api"]
```

---

## 9. 风险与缓解

| 风险 | 缓解 |
|:---|:---|
| tRPC-Agent-Go 功能不足 | Phase 4 前做 PoC 验证 StateGraph + SSE + MCP 能力 |
| PgVector → Milvus 迁移 | PgVector 仅存元数据+embedding，迁移策略：执行 migration 脚本将历史向量批量写入 Milvus，Go 端 embedding 全量走 Milvus。过渡期间 Java 继续用 PgVector，两端共存不冲突 |
| Java 遗留功能 Go 不好实现 | Spring AI 的某些特性（如 Graph Checkpoint）可能需自研，Phase 5 做详细评估 |
| 前端 API 差异 | 每个 Phase 结束后，用 Java 后端的 API 测试集验证 Go 后端 |
| 团队 Go 熟练度 | Phase 1-2 是纯粹的 CRUD + 基础设施，难度低，适合磨合 |

---

## 10. 完成标准

每个 Phase 的完成标准：

- [ ] 所有 handler 对齐 Java Controller 的 REST 契约
- [ ] 单元测试覆盖核心领域逻辑（`go test ./...`）
- [ ] 集成测试覆盖 DB + Cache + Queue
- [ ] docker-compose 可一键启动验证
- [ ] 前端对该模块的页面功能正常
