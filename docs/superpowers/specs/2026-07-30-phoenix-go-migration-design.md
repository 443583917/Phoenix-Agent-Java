# Phoenix Go 重写 — 设计文档

> 将 Phoenix-Agent-Java（Spring AI Alibaba）后端全部迁移到 Go，前端 API 契约保持不变。

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
├── cmd/                        # 多服务入口
│   ├── api/main.go             # HTTP API 服务（端口 8066）
│   ├── rpc/main.go             # gRPC 服务
│   ├── agent/main.go           # Agent 独立服务
│   └── job/main.go             # 定时任务服务
│
├── api/                        # HTTP API 层（Gin）
│   ├── router.go               # 总路由注册
│   ├── handler/
│   │   ├── agent/              # 智能体 API handler
│   │   ├── datasource/         # 数据源 API handler
│   │   ├── chat/               # 对话/SSE API handler
│   │   ├── knowledge/          # 知识库 API handler
│   │   ├── modelconfig/        # 模型配置 API handler
│   │   ├── prompt/             # Prompt 配置 API handler
│   │   ├── semanticmodel/      # 语义模型 API handler
│   │   ├── privilege/          # 权限认证 API handler
│   │   ├── platform/           # 平台管理 API handler
│   │   ├── rag/                # RAG API handler
│   │   ├── kg/                 # 知识图谱 API handler
│   │   └── common/             # 公共 API handler
│   └── middleware/
│       ├── auth.go             # OAuth2 JWT 认证
│       ├── rbac.go             # Casbin RBAC
│       ├── cors.go
│       ├── logger.go           # Zap 日志中间件
│       ├── recovery.go
│       ├── tracing.go          # OpenTelemetry
│       └── ratelimit.go
│
├── rpc/                        # gRPC 层（Phase 6 实现，前期为空骨架）
│   ├── proto/
│   │   ├── agent.proto
│   │   ├── privilege.proto
│   │   └── data.proto
│   ├── server/                 # gRPC 服务端实现
│   └── client/                 # gRPC 客户端（服务间调用）
│
├── agent/                      # tRPC-Agent-Go 层（核心 AI 引擎）
│   ├── runtime/
│   │   ├── manager.go          # AgentManager — 生命周期管理
│   │   └── registry.go         # Agent 注册表
│   ├── agents/
│   │   ├── react_agent.go      # React Agent（推理-行动-观察）
│   │   ├── workflow_agent.go   # Workflow Graph Agent
│   │   └── assistant_agent.go  # 通用助手 Agent
│   ├── tools/
│   │   ├── registry.go         # 工具注册中心（对标 @Tool 注解）
│   │   ├── agent/              # Agent 相关工具
│   │   ├── datasource/         # 数据源查询工具
│   │   └── privilege/          # 权限检查工具
│   ├── workflows/
│   │   ├── graph.go            # StateGraph 构建器
│   │   ├── nodes/              # NL2SQL 工作流节点
│   │   │   ├── intent.go       # 意图识别
│   │   │   ├── evidence.go     # 证据召回
│   │   │   ├── schema.go       # Schema 召回
│   │   │   ├── planner.go      # 计划生成
│   │   │   ├── sql_gen.go      # SQL 生成
│   │   │   ├── python_exec.go  # Python 执行
│   │   │   └── report.go       # 报告生成
│   │   └── checkpoint.go       # Redis 状态检查点
│   ├── memory/
│   │   ├── short_term.go       # 短期记忆（对话窗口）
│   │   ├── long_term.go        # 长期记忆（Milvus 向量）
│   │   └── profile.go          # 用户画像
│   ├── knowledge/
│   │   ├── retriever.go        # 混合检索（向量+关键词+RRF）
│   │   ├── embedding.go        # Embedding 模型代理
│   │   └── splitter.go         # 文档分割器
│   └── runner/
│       ├── runner.go           # Runner — 执行对话管道
│       ├── sse.go              # SSE 事件流
│       └── hitl.go             # Human-in-the-Loop
│
├── internal/                   # 业务核心层
│   ├── domain/                 # 领域层（纯业务规则，零框架依赖）
│   │   ├── agent/              # 智能体领域
│   │   ├── datasource/         # 数据源领域
│   │   ├── chat/               # 对话领域
│   │   ├── knowledge/          # 知识库领域
│   │   ├── modelconfig/        # 模型配置领域
│   │   ├── prompt/             # Prompt 领域
│   │   ├── semanticmodel/      # 语义模型领域
│   │   ├── privilege/          # 权限领域
│   │   ├── platform/           # 平台领域
│   │   ├── rag/                # RAG 领域
│   │   └── common/             # 公共领域
│   ├── service/                # 应用服务层（编排）
│   │   ├── agent_service.go
│   │   ├── datasource_service.go
│   │   ├── chat_service.go
│   │   ├── knowledge_service.go
│   │   ├── privilege_service.go
│   │   ├── platform_service.go
│   │   └── rag_service.go
│   ├── dao/                    # 数据访问层
│   │   ├── db/                 # PostgreSQL (GORM)
│   │   │   ├── agent_repo.go
│   │   │   ├── datasource_repo.go
│   │   │   ├── chat_repo.go
│   │   │   ├── privilege_repo.go
│   │   │   ├── platform_repo.go
│   │   │   └── rag_repo.go
│   │   ├── cache/              # Redis + BigCache
│   │   │   ├── agent_cache.go
│   │   │   ├── privilege_cache.go
│   │   │   └── bigcache.go
│   │   ├── queue/              # RabbitMQ + Redis
│   │   │   ├── producer.go
│   │   │   └── consumer.go
│   │   └── external/           # 外部集成
│   │       ├── milvus.go       # Milvus 向量操作
│   │       ├── docker_exec.go  # Docker Python 执行器
│   │       └── oss.go          # 对象存储
│   ├── event/                  # 领域事件
│   │   ├── types.go
│   │   └── handler/
│   │       ├── agent/
│   │       ├── chat/
│   │       └── privilege/
│   ├── job/                    # 定时任务
│   │   ├── scheduler.go
│   │   └── jobs/
│   ├── model/                  # 数据模型
│   │   ├── agent_entity.go
│   │   ├── datasource_entity.go
│   │   ├── chat_entity.go
│   │   ├── privilege_entity.go
│   │   ├── platform_entity.go
│   │   ├── rag_entity.go
│   │   ├── agent_dto.go
│   │   ├── chat_dto.go
│   │   ├── privilege_dto.go
│   │   ├── agent_vo.go
│   │   └── privilege_vo.go
│   └── config/                 # 配置加载（Viper）
│       ├── app.go
│       ├── db.go
│       ├── redis.go
│       ├── milvus.go
│       ├── rabbitmq.go
│       ├── agent.go
│       ├── rpc.go
│       └── monitor.go
│
├── pkg/                        # 基础设施包
│   ├── logger/                 # Zap 封装
│   ├── errcode/                # 统一错误码
│   ├── response/               # 统一响应格式
│   ├── jwt/                    # JWT 工具
│   ├── pagination/             # 分页工具
│   ├── sse/                    # SSE 工具
│   ├── utils/                  # 通用工具
│   └── validate/               # 参数校验（go-playground/validator）
│
├── configs/                    # YAML 配置文件
│   ├── api.yaml
│   ├── db.yaml
│   ├── redis.yaml
│   ├── milvus.yaml
│   ├── rabbitmq.yaml
│   ├── agent.yaml
│   ├── monitor.yaml
│   └── oss.yaml
│
├── migrations/                 # DDL 迁移（golang-migrate）
├── scripts/
├── docs/
├── Dockerfile
├── docker-compose.yaml
├── Makefile
├── go.mod
└── go.sum
```

---

## 3. Java 模块 → Go 包映射

```
phoenix-common-api/core/rest     →  pkg/ + internal/model/common_* + internal/domain/common/
phoenix-privilege-api/core/rest  →  internal/domain/privilege/ + api/handler/privilege/ + internal/dao/db/privilege_repo.go
phoenix-platform-api/core/rest   →  internal/domain/platform/ + api/handler/platform/ + internal/dao/db/platform_repo.go
phoenix-agent-api/core/rest      →  agent/ + internal/domain/agent/ + api/handler/agent/
phoenix-data-api/core/rest       →  agent/workflows/ + internal/domain/datasource/ + internal/domain/chat/ + ...
phoenix-rag-api/core/rest        →  internal/domain/rag/ + agent/knowledge/ + api/handler/rag/
phoenix-kg-api/core/rest         →  internal/domain/kg/ + api/handler/kg/
phoenix-tool                     →  pkg/utils/
phoenix-codegen                  →  (用不着，Go 没有代码生成器需求)
phoenix-flink                     →  internal/dao/external/flink.go（降级为外部集成）
phoenix-admin-manager             →  cmd/api/main.go（统一入口）
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

### 4.2 新增模块使用 golang-migrate

```
migrations/
├── 000001_create_sessions_table.up.sql
├── 000001_create_sessions_table.down.sql
└── ...
```

---

## 5. 统一 API 规范

### 5.1 响应格式（对标 Java R 类）

```go
// pkg/response/response.go
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
        // 智能体管理
        agentGroup := api.Group("/agent")
        agent.RegisterRoutes(agentGroup)

        // 数据源管理
        datasourceGroup := api.Group("/datasource")
        datasource.RegisterRoutes(datasourceGroup)

        // 对话
        chatGroup := api.Group("")
        chat.RegisterRoutes(chatGroup)

        // ...其他模块
    }

    // 权限认证（无需 Auth，但需要独立 JWT 校验）
    privilegeGroup := r.Group("/api/privilege")
    privilege.RegisterRoutes(privilegeGroup)

    // 平台管理
    platformGroup := r.Group("/platform")
    platformGroup.Use(middleware.Auth())
    platform.RegisterRoutes(platformGroup)

    return r
}
```

### 5.3 SSE 兼容（对标 WebFlux SSE）

```go
// pkg/sse/sse.go
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

### Phase 1 — 基础设施层（pkg/ + configs/ + internal/config/ + cmd/api/）

**目标**：项目骨架跑通，Gin 启动并响应 `/echo`

| 交付物 | 内容 |
|:---|:---|
| `go.mod` | `module github.com/phoenix-agent-go` |
| `cmd/api/main.go` | Gin 启动，Viper 加载配置，注册中间件 |
| `configs/*.yaml` | 全套配置文件 |
| `pkg/logger/` | Zap 封装，支持日志级别、文件轮转 |
| `pkg/response/` | 统一 Response / PageResponse |
| `pkg/errcode/` | 统一错误码定义 |
| `pkg/jwt/` | JWT 生成/验证 |
| `pkg/pagination/` | 分页参数解析 |
| `pkg/sse/` | SSE 流式输出 |
| `pkg/validate/` | 参数校验封装 |
| `internal/config/` | Viper 配置加载（db/redis/milvus/rabbitmq/agent/monitor） |
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
| `internal/dao/db/privilege_repo.go` | GORM CRUD |
| `internal/dao/cache/privilege_cache.go` | Redis 权限缓存 + BigCache L1 |
| `internal/service/privilege_service.go` | LoginService, UserService, RoleService... |
| `api/handler/privilege/` | LoginController → 12 个 handler |
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
| `internal/model/common_entity.go` | PlatformInfo (third-party config) |
| `internal/domain/platform/` | 平台领域规则 |
| `internal/dao/db/platform_repo.go` | GORM CRUD |
| `internal/dao/external/dingtalk.go` | 钉钉 SDK |
| `internal/dao/external/feishu.go` | 飞书 SDK |
| `internal/dao/external/wecom.go` | 企业微信 SDK |
| `internal/service/platform_service.go` | 应用服务 |
| `api/handler/platform/` | GroupAgentInfo, AccountLogin 等 handler |
| `api/handler/common/` | PlatformInfo, PlatformSync handler |

### Phase 4 — Agent 框架（phoenix-agent）

**目标**：React Agent + 会话管理 + 记忆 + MCP，核心 AI 能力

| 交付物 | 映射 Java |
|:---|:---|
| `agent/runtime/manager.go` | AgentManager |
| `agent/agents/react_agent.go` | ReactAgent |
| `agent/tools/registry.go` | @Tool 注解扫描 → 显式注册 |
| `agent/tools/sql_tool.go` | SQL 执行工具 |
| `agent/tools/knowledge_tool.go` | 知识检索工具 |
| `agent/tools/mcp_tool.go` | MCP 协议工具 |
| `agent/memory/short_term.go` | 对话窗口管理（Redis checkpoint） |
| `agent/memory/long_term.go` | Milvus 向量检索记忆 |
| `agent/memory/profile.go` | 用户画像 |
| `agent/knowledge/` | 混合检索（向量 + 关键词 + RRF 融合） |
| `agent/runner/runner.go` | 对话执行管道 |
| `agent/runner/sse.go` | SSE 流式输出 |
| `agent/runner/hitl.go` | Human-in-the-Loop |
| `internal/model/agent_entity.go` | UserAgentInfo, UserMemoryInfo, CombinedStore |
| `internal/dao/db/agent_repo.go` | GORM CRUD |
| `internal/dao/external/milvus.go` | Milvus 操作封装 |
| `internal/service/agent_service.go` | Agent 应用服务 |
| `api/handler/agent/` | ReactAgentController, HarnessController |

### Phase 5 — NL2SQL 数据引擎（phoenix-data）

**目标**：StateGraph 工作流，自然语言转 SQL/Python 分析，最复杂的模块

| 交付物 | 映射 Java |
|:---|:---|
| `agent/workflows/graph.go` | DataAgentConfiguration StateGraph |
| `agent/workflows/nodes/intent.go` | IntentRecognitionNode |
| `agent/workflows/nodes/evidence.go` | EvidenceRecallNode |
| `agent/workflows/nodes/schema.go` | SchemaRecallNode + QueryEnhanceNode + TableRelationNode |
| `agent/workflows/nodes/planner.go` | PlannerNode + FeasibilityAssessmentNode + PlanExecutorNode |
| `agent/workflows/nodes/sql_gen.go` | SqlGenerateNode + SemanticConsistencyNode + SqlExecuteNode |
| `agent/workflows/nodes/python_exec.go` | PythonGenerateNode + PythonExecuteNode + PythonAnalyzeNode |
| `agent/workflows/nodes/report.go` | ReportGeneratorNode |
| `agent/workflows/checkpoint.go` | Redis 状态检查点 |
| `internal/domain/datasource/` | 数据源连接管理 |
| `internal/dao/db/datasource_repo.go` | 多数据库连接器 (MySQL/PG/Oracle/MSSQL/H2/Hive/DM) |
| `internal/dao/db/ddl_repo.go` | DDL 元数据提取 |
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

**目标**：RAG, KG, Flink, 全量测试, 生产部署

| 交付物 | 映射 Java |
|:---|:---|
| `internal/domain/rag/` | RAG 文件管理 |
| `api/handler/rag/` | RagFileInfoController, RagCategoryController |
| `internal/domain/kg/` | 知识图谱 |
| `api/handler/kg/` | KG handler |
| `internal/dao/external/flink.go` | Flink 集成 |
| `cmd/rpc/main.go` | gRPC 服务启动 |
| `cmd/agent/main.go` | Agent 独立服务启动 |
| `cmd/job/main.go` | 定时任务服务启动 |
| `docker-compose.yaml` | 完整生产环境编排 |

---

## 7. 关键设计决策

### 7.1 为什么用按模块垂直切分

```
internal/
├── domain/privilege/     ← privilege 模块的领域逻辑
├── dao/db/privilege_repo.go   ← privilege 模块的数据库访问
├── dao/cache/privilege_cache.go  ← privilege 模块的缓存
api/
├── handler/privilege/    ← privilege 模块的 HTTP handler
```

每个模块在每一层都有自己的文件/目录，**水平分层 + 垂直模块**。新增模块只需在每层加一个目录。

### 7.2 并发模型：Handler → Service → DAO 全链路 ctx 传递

```go
func (h *AgentHandler) GetAgent(c *gin.Context) {
    ctx := c.Request.Context()  // 来自 Gin 的 context
    agent, err := h.service.GetAgent(ctx, id)
}

func (s *AgentService) GetAgent(ctx context.Context, id uint64) (*model.Agent, error) {
    // 先查 BigCache L1
    if cached := s.cache.Get(ctx, id); cached != nil { return cached, nil }
    // 再查 GORM
    agent, err := s.repo.FindByID(ctx, id)
    // 写回 L1
    s.cache.Set(ctx, id, agent)
    return agent, err
}
```

Gin 天然每请求一个 goroutine，无需 WebFlux reactive，代码直观。

### 7.3 事务管理

Go 没有 `@Transactional` 注解，通过 `db.Transaction()` 显式管理：

```go
func (s *AgentService) CreateAgent(ctx context.Context, dto CreateAgentDTO) (*model.Agent, error) {
    var agent *model.Agent
    err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        agent = dto.ToEntity()
        if err := tx.Create(agent).Error; err != nil { return err }
        // 同时创建默认数据源
        ds := dto.ToDefaultDatasource(agent.ID)
        return tx.Create(ds).Error
    })
    return agent, err
}
```

### 7.4 错误处理

统一错误码，不使用 exception。对标 Java 的 BizException / GlobalExceptionHandler：

```go
// pkg/errcode/errcode.go
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

YAML 文件分离，Viper 合并加载：

```yaml
# configs/db.yaml
database:
  host: "127.0.0.1"
  port: 5432
  user: "phoenix"
  password: "phoenix"
  name: "phoenix"
  sslmode: "disable"
  max_open_conns: 25
  max_idle_conns: 10

# configs/redis.yaml
redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0

# configs/milvus.yaml
milvus:
  addr: "127.0.0.1:19530"
  collection: "phoenix_vectors"
  dim: 512

# configs/agent.yaml
agent:
  model:
    provider: "deepseek"
    model: "deepseek-chat"
    api_key: "${AI_API_KEY}"
    base_url: "https://api.deepseek.com"
  stream: true
  max_tokens: 4096
```

---

## 8. 部署架构

```
docker-compose.yaml
├── phoenix-go (Gin API, :8066)
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
COPY configs/ /etc/phoenix/
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
