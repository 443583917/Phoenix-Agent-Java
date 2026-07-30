# Phoenix Go 重写 — 交接文档

> 日期：2026-07-31  
> 分支：main  
> 状态：核心功能完成，增量优化阶段

---

## 1. 项目总览

| 指标 | 数值 |
|:---|:---|
| Go 文件 | 161 |
| 代码行数 | 21,551 |
| API 端点 | 300+ |
| Java 模块覆盖 | 12/12 |
| `go build` | PASS |
| `go vet` | PASS |

### 目录结构

```
godata/
├── cmd/           4 服务入口（api, rpc, agent, job）
├── api/           Gin HTTP 层（12 handler 模块 + 7 中间件）
├── rpc/           gRPC 层（骨架）
├── agent/         tRPC-Agent-Go AI 引擎
│   ├── runtime/    AgentManager + Registry
│   ├── agents/     ReactAgent, AssistantAgent
│   ├── tools/      工具注册 + 数据源 SQL 工具
│   ├── workflows/  NL2SQL (16节点) + RAG + KG
│   ├── memory/     短期/长期记忆 + 记忆提取管道
│   ├── knowledge/  嵌入 + 混合检索 + 文本分割
│   ├── hooks/      3 个 Hook（Profile/Limit/Summarization）
│   ├── interceptors/ 登录拦截器
│   └── runner/     SSE + HITL + Redis 确认缓存
├── internal/      业务核心（DDD）
│   ├── domain/     领域逻辑（privilege 已提取，其余在 usecase）
│   ├── usecase/    用例编排（privilege, platform, data, agent）
│   ├── service/    应用服务（薄层）
│   ├── repository/ 仓储接口
│   ├── dao/        数据访问（GORM, Redis, RabbitMQ, Milvus）
│   ├── model/      数据模型（Entity, DTO, VO）
│   └── config/     Viper 配置结构体
├── infra/          基础设施（logger, cache, queue, lock, monitoring...）
└── configs/        YAML 配置文件
```

---

## 2. 已完成功能

### 2.1 Java 模块 → Go 映射（全部完成）

| Java 模块 | Go 位置 | 核心功能 |
|:---|:---|:---|
| phoenix-privilege | Phase 2 | RBAC（12表），JWT+Casbin，85+ 端点 |
| phoenix-platform | Phase 3 | 平台管理（7表），50+ 端点 |
| phoenix-agent | Phase 4 | Agent 框架，SSE/HITL，tRPC-Agent-Go |
| phoenix-data | Phase 5 | NL2SQL 引擎（14表），100+ 端点 |
| phoenix-common | Phase 1/3 | BaseModel, PlatformInfo, 基础设施 |
| phoenix-tool | Phase 1 | 统一响应，错误码，SQL 校验 |
| phoenix-rag | Phase 6 | RAG 文件管理 + 工作流 |
| phoenix-kg | Phase 6 | 知识图谱 CRUD + 工作流 |
| phoenix-flink | Phase 6 | SDK stub |
| phoenix-admin | Phase 1 | cmd/api/main.go 统一入口 |
| phoenix-parent | Phase 1 | go.mod 依赖管理 |
| phoenix-codegen | — | 不需要（Go无代码生成器） |

### 2.2 NL2SQL StateGraph 工作流（完全同步）

16 个节点全部实现，与 Java `DataAgentConfiguration` 对齐：

```
START → IntentRecognition(LLM) → EvidenceRecall(LLM+Vector) → QueryEnhance(LLM)
→ SchemaRecall(Vector) → TableRelation(LLM) → FeasibilityAssessment(LLM)
→ Planner(LLM) → PlanExecutor → [SqlGenerate(LLM)→SemanticConsistency(LLM)→SqlExecute(DB)]
                        → [PythonGenerate(LLM)→PythonExecute(Docker)→PythonAnalyze(LLM)]
                        → [ReportGenerator(LLM)]
                        → [HumanFeedback(Interrupt/Resume)]
→ END
```

- 9 个 Prompt 模板（Go text/template）
- 11 个条件分发器（重试/回退/修复循环）
- SQL 语义重试 + 执行重试（最多3次）
- Python 重试 + 回退模式
- 计划修复循环（最多3次）
- 人工干预中断/恢复

### 2.3 Agent 框架（完全同步）

| Java 组件 | Go 实现 |
|:---|:---|
| AgentManager（类型路由） | `agent/runtime/manager.go` |
| CombinedDbHook（画像注入） | `agent/hooks/profile_hook.go` |
| ModelCallLimitHook（调用限制） | `agent/hooks/limit_hook.go` |
| SummarizationHook（对话摘要） | `agent/hooks/summarization_hook.go` |
| LoginUserAgentInterceptor（追踪+记忆注入） | `agent/interceptors/login_interceptor.go` |
| MemoryPipelineService（异步记忆提取） | `agent/memory/pipeline.go` |
| HitlCacheService（Redis 确认缓存） | `agent/runner/hitl_cache.go` |

### 2.4 基础设施

| 组件 | 状态 |
|:---|:---|
| Gin HTTP + 中间件 | ✅ |
| GORM + PostgreSQL | ✅ |
| JWT + Casbin RBAC | ✅ |
| Redis 缓存 + 分布式锁 | ✅ |
| RabbitMQ Producer + Consumer | ✅ |
| Milvus 向量存储 | ✅ |
| BigCache L1 缓存 | ✅ |
| OpenTelemetry 追踪 | ✅ |
| Viper 配置管理 | ✅ |
| Zap 日志 | ✅ |
| Sonyflake ID 生成 | ✅ |
| Docker Python 执行器 | ✅ |
| golang-migrate | ✅ 骨架 |

---

## 3. 未完成 / 待填充

### 3.1 功能缺口（需实现）

| 优先级 | 内容 | 当前状态 | Java 对应 |
|:---|:---|:---|:---|
| 🔴 高 | **KG 工作流节点** | `agent/workflows/kg/nodes/` 空目录 | phoenix-kg 有完整 KG 图 |
| 🔴 高 | **gRPC proto 定义 + 实现** | `rpc/` 空骨架 | 设计规范预留 |
| 🔴 高 | **数据库迁移文件** | `migrations/` 空 | `sql/all_schema.sql` |
| 🟡 中 | **MCP 外部工具** | `agent/tools/external/` 空 | Java 有 MCP 协议支持 |
| 🟡 中 | **DDD 领域提取** | 9 个 `internal/domain/` 目录空 | 逻辑在 usecase 层 |
| 🟡 中 | **事件驱动架构** | `internal/event/` 空 | Java 有 Redis Pub/Sub 事件 |
| 🟡 中 | **定时任务** | `internal/job/jobs/` 空 | Java 有 @Scheduled |
| 🟡 中 | **钉钉/飞书/企微 SDK** | `internal/dao/external/` 仅 stub | Java 有完整实现 |
| 🟢 低 | **Web 搜索工具** | `agent/tools/web/` 空 | 设计规范预留 |
| 🟢 低 | **RPC 工具** | `agent/tools/rpc/` 空 | 设计规范预留 |

### 3.2 代码质量

| 项目 | 状态 | 说明 |
|:---|:---|:---|
| 单元测试 | ⚠️ 仅 2 个包有测试 | infra/config, infra/logger |
| 集成测试 | ❌ 无 | 需补充 DB/Cache/Queue 集成测试 |
| API 契约测试 | ❌ 无 | 需验证 Go 端点与 Java 返回格式一致 |
| 性能测试 | ❌ 无 | — |

### 3.3 部署

| 项目 | 状态 |
|:---|:---|
| Dockerfile | ✅ 多阶段构建 |
| docker-compose | ✅ 8 个服务 |
| Makefile | ✅ 14 个目标 |
| CI/CD | ❌ 无 |
| Kubernetes | ❌ 无 |

---

## 4. 详细重构计划（未完成部分）

### 4.1 KG 工作流节点（高优先，预估 3h）

**文件：** `agent/workflows/kg/nodes/`

创建 KG 工作流节点（对标 Java `phoenix-kg` 模块）：

```go
// agent/workflows/kg/nodes/entity_extract.go
// 实体提取节点：LLM 从文本提取知识图谱实体

// agent/workflows/kg/nodes/relation_extract.go
// 关系提取节点：LLM 识别实体间关系

// agent/workflows/kg/nodes/graph_merge.go
// 图合并节点：融合新提取的三元组到已有图

// agent/workflows/kg/graph.go
// KG 工作流图：entity_extract → relation_extract → graph_merge → END
```

### 4.2 gRPC 服务（高优先，预估 4h）

**文件：** `rpc/proto/`, `rpc/server/`, `rpc/client/`, `cmd/rpc/main.go`

1. 定义 proto 文件：
   - `rpc/proto/agent.proto` — AgentService (StreamCall, ListAgents)
   - `rpc/proto/privilege.proto` — PrivilegeService (Login, GetUser, ListRoles)
   - `rpc/proto/data.proto` — DataService (StreamNL2SQL, ChatSession)

2. 生成 Go 代码：`protoc --go_out=. --go-grpc_out=.`

3. 实现 gRPC server：
   ```go
   // cmd/rpc/main.go
   lis, _ := net.Listen("tcp", ":9090")
   s := grpc.NewServer()
   pb.RegisterAgentServiceServer(s, agentServer)
   pb.RegisterPrivilegeServiceServer(s, privilegeServer)
   ```

4. gRPC 客户端（供服务间调用）

### 4.3 数据库迁移（高优先，预估 1h）

```bash
# 从现有 SQL 生成 migration
cd godata
migrate create -ext sql -dir migrations -seq init_schema
# 复制 sql/all_schema.sql 内容到 up 文件
# 复制 DROP 语句到 down 文件
```

### 4.4 MCP 外部工具（中优先，预估 2h）

**文件：** `agent/tools/external/mcp.go`

对标 Java `McpServerService`：
- MCP 协议客户端注册
- 发现远程工具并注册到 ToolRegistry
- SSE 传输层支持

### 4.5 DDD 领域提取（中优先，预估 3h）

将 usecase 中的业务规则提取到 domain 层：

```
internal/usecase/privilege_usecase.go
  → internal/domain/privilege/
    - user_rules.go（密码规则、用户校验）
    - role_rules.go（角色校验）
    - acl_rules.go（权限位掩码合并）

internal/usecase/platform_usecase.go
  → internal/domain/platform/
    - account_rules.go
    - group_rules.go

... 其他 7 个 domain 目录同理
```

### 4.6 事件驱动架构（中优先，预估 2h）

**文件：** `internal/event/`

```go
// internal/event/types.go
type Event struct {
    ID        string
    Type      string
    Payload   json.RawMessage
    Timestamp time.Time
}

// internal/event/bus.go
type EventBus struct {
    queue *queue.Consumer
    handlers map[string][]EventHandler
}

// internal/event/handler/privilege/user_created.go
// 用户创建后：初始化默认角色、记录审计日志
```

### 4.7 定时任务（中优先，预估 1h）

**文件：** `internal/job/jobs/`

```go
// internal/job/jobs/daily_report.go
// 每日报表生成

// internal/job/jobs/audit_sync.go
// 审计日志同步

// internal/job/scheduler.go
// 使用 robfig/cron 注册定时任务
```

### 4.8 单元测试补充（中优先，预估 4h）

优先补充测试的包：
1. `internal/usecase/` — 业务逻辑（Mock Repository）
2. `api/handler/privilege/` — HTTP handler（httptest）
3. `agent/memory/` — 记忆管道
4. `agent/workflows/nl2sql/nodes/` — 节点逻辑

### 4.9 三方 SDK（中优先，预估 3h）

**文件：** `internal/dao/external/`

对标 Java 实现：
- `dingtalk.go` — `GetUserIDByMobile()`, `SendMessage()`
- `feishu.go` — `GetUserAccessToken()`, `GetUserInfo()`
- `wecom.go` — `GetAccessToken()`, `SendMessage()`

### 4.10 Web 搜索 + RPC 工具（低优先，预估 2h）

**文件：** `agent/tools/web/`, `agent/tools/rpc/`

- Web 搜索工具：集成 DuckDuckGo API（tRPC-Agent-Go 已有 `tool/duckduckgo/`）
- RPC 工具：跨服务调用工具

---

## 5. Superpowers 插件使用指引

### 已使用的 Skills

| Skill | 用途 | 产出 |
|:---|:---|:---|
| `superpowers:brainstorming` | 需求分析、方案设计 | `specs/2026-07-30-phoenix-go-migration-design.md` |
| `superpowers:writing-plans` | 编写实施计划 | 6 个 Phase 计划文件 |
| `superpowers:subagent-driven-development` | 任务执行 + 审查 | 全部 161 个 Go 文件 |

### 推荐继续使用的 Skills

```bash
# 1. 继续执行未完成计划
/superpowers:writing-plans          # 为 KG 工作流 / gRPC / 迁移编写计划
/superpowers:subagent-driven-development  # 逐个任务执行 + 审查

# 2. 代码质量
/superpowers:requesting-code-review  # 全量代码审查
/superpowers:systematic-debugging    # 排查问题

# 3. 部署
/superpowers:finishing-a-development-branch  # 合并/PR/推送

# 4. 如果切换工作目录
/superpowers:using-git-worktrees     # 创建隔离工作树
```

### 关键命令

```bash
cd d:/GitHub/Phoenix-Agent-Java/godata

# 编译
go build ./...

# 测试
go test ./... -short

# 静态检查
go vet ./...

# 依赖管理
go mod tidy

# 启动服务
go run ./cmd/api

# Docker
docker build -t phoenix-api .
docker-compose up -d
```

---

## 6. 配置文件说明

| 文件 | 内容 |
|:---|:---|
| `configs/db.yaml` | PostgreSQL 连接 |
| `configs/redis.yaml` | Redis 连接 |
| `configs/milvus.yaml` | Milvus 向量库 |
| `configs/rabbitmq.yaml` | RabbitMQ 消息队列 |
| `configs/monitor.yaml` | OpenTelemetry + 日志 |
| `configs/api/app.yaml` | HTTP 端口 + CORS + JWT |
| `configs/agent/app.yaml` | AI 模型配置 |

环境变量覆盖：前缀 `PHOENIX_`，如 `PHOENIX_DATABASE_HOST=prod-db`

---

## 7. 快速恢复指引

```bash
# 1. 克隆项目
cd d:/GitHub/Phoenix-Agent-Java/godata

# 2. 安装依赖
go mod download

# 3. 验证编译
go build ./... && go vet ./...

# 4. 启动基础设施
docker-compose up -d postgres redis milvus rabbitmq

# 5. 导入数据库
psql -h localhost -U phoenix -d phoenix < ../sql/all_schema.sql
psql -h localhost -U phoenix -d phoenix < ../sql/all_data.sql

# 6. 启动服务
go run ./cmd/api

# 7. 验证
curl http://localhost:8066/echo
# → {"code":0,"message":"success","data":"ok"}
```

---

## 8. 架构决策记录

| 决策 | 原因 |
|:---|:---|
| 使用 tRPC-Agent-Go 而非自建 Agent | 框架提供 llmagent/runner/tool/session/memory/graph 全套 API |
| 保持 API 契约不变 | 前端（Vben5）零修改 |
| DDD 分层但领域提取后置 | 优先功能完整性，领域逻辑先在 usecase 中 |
| 全部使用 tRPC-Agent-Go API | 嵌入/记忆/知识/会话全部走框架，零自定义垫片 |
| 内存优先 + Redis 持久化 | 短期记忆用 inmemory session，长期记忆用 Redis |
| 不修改数据库 Schema | 复用 Java 现有表结构 |

---

*此文档由 Superpowers 插件辅助生成。最后更新：2026-07-31*
