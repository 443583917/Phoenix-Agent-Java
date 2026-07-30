# Phoenix Go 重写 — 交接与可执行重构计划

> 日期：2026-07-31 | 分支：main | `go build ./...` PASS | `go vet ./...` PASS

## 当前状态

161 个 Go 文件，21,551 行代码。所有 12 个 Java 模块已映射，核心架构齐全。发现 **16 个功能缺口**（7 个 stub handler + 9 个空目录/占位），需要补全。

---

## 目录架构与代码归属规则

```
godata/
├── cmd/                            # ====== 服务入口（每个独立进程） ======
│   ├── api/main.go                 # HTTP API 服务，主入口（端口 8066）
│   ├── rpc/main.go                 # gRPC 服务（骨架，Task 5 补全）
│   ├── agent/main.go               # Agent 独立服务（骨架）
│   └── job/main.go                 # 定时任务服务（骨架）
│
├── api/                            # ====== HTTP 层 ======
│   ├── router.go                   # 总路由注册 — 新增 handler 必须在此注册
│   └── handler/                    # 控制器（按业务域拆分，每域一个子目录）
│       ├── privilege/              # → 权限认证（12 个 handler 文件）
│       ├── platform/               # → 平台管理（7 个 handler 文件）
│       ├── agent/                  # → 智能体 API（SSE 流式对话）
│       ├── chat/                   # → 对话/SSE/Graph 流式搜索
│       ├── datasource/             # → 数据源管理 + Agent 数据源关联
│       ├── knowledge/              # → Agent 知识库 + 业务知识
│       ├── modelconfig/            # → AI 模型配置
│       ├── prompt/                 # → Prompt 模板配置
│       ├── semanticmodel/          # → 语义模型
│       ├── rag/                    # → RAG 文件管理
│       ├── kg/                     # → 知识图谱
│       └── common/                 # → 公共（上传/平台信息/同步）
│   └── middleware/                 # Gin 中间件（auth/rbac/cors/logger/tracing/recovery/ratelimit）
│
├── rpc/                            # ====== gRPC 层（骨架，Task 5 补全） ======
│   ├── proto/                      # .proto 定义文件
│   ├── server/                     # gRPC server 实现
│   └── client/                     # gRPC 客户端
│
├── agent/                          # ====== AI Agent 引擎 ======
│   ├── runtime/                    # Agent 生命周期管理
│   │   ├── manager.go              # AgentManager — 创建、路由、流式调用
│   │   └── registry.go             # Agent 注册表（SN 索引）
│   ├── agents/                     # Agent 类型
│   │   ├── react_agent.go          # ReactAgent 构建器
│   │   └── assistant_agent.go      # 无工具通用助手
│   ├── tools/                      # 工具注册 + 实现
│   │   ├── registry.go             # ToolRegistry — 全局工具注册表
│   │   ├── function/calculator.go  # 计算器 FunctionTool 示例
│   │   ├── datasource/sql_tool.go  # ⭐ 数据源 SQL 查询工具（NL2SQL 关键）
│   │   ├── external/               # 🟡 MCP 外部工具（待 Task 7）  ← 新工具往这加
│   │   ├── privilege/              # 🟢 权限检查工具（待 Task 12） ← 新工具往这加
│   │   ├── rpc/                    # 🟢 RPC 调用工具（待 Task 12）  ← 新工具往这加
│   │   └── web/                    # 🟢 Web 搜索工具（待 Task 12）  ← 新工具往这加
│   ├── workflows/                  # 工作流图（StateGraph）
│   │   ├── nl2sql/                 # ⭐ NL2SQL 工作流（16 节点，完整 LLM 实现）
│   │   │   ├── graph.go            # StateGraph 构建器
│   │   │   ├── nodes/              # 16 个节点文件
│   │   │   ├── prompts/            # 9 个 LLM Prompt 模板
│   │   │   └── types/              # 类型定义 + 状态常量
│   │   ├── rag/                    # RAG 工作流（stub + retrieve 节点）
│   │   │   ├── graph.go
│   │   │   └── nodes/
│   │   └── kg/                     # 🟡 KG 工作流（待 Task 4）   ← KG 节点写这里
│   │       ├── graph.go
│   │       └── nodes/              # entity_extract / relation_extract / graph_merge
│   ├── memory/                     # 记忆系统
│   │   ├── short_term.go           # 短期记忆（tRPC session/inmemory）
│   │   ├── long_term.go            # 长期记忆（tRPC memory/inmemory）
│   │   ├── profile.go              # 用户画像（GORM）
│   │   └── pipeline.go             # 异步记忆提取管道（LLM → profile/facts/vector）
│   ├── knowledge/                  # 知识检索
│   │   ├── embedding.go            # OpenAI 嵌入（tRPC embedder/openai）
│   │   ├── retriever.go            # 混合检索（tRPC knowledge）
│   │   └── splitter.go             # 文本分割
│   ├── hooks/                      # Agent Hook（BEFORE_MODEL）
│   │   ├── profile_hook.go         # 画像注入 Hook
│   │   ├── limit_hook.go           # 模型调用限制 Hook
│   │   └── summarization_hook.go   # 对话摘要 Hook
│   ├── interceptors/               # Agent 拦截器
│   │   └── login_interceptor.go    # 登录追踪 + 历史记忆注入
│   └── runner/                     # 对话执行器
│       ├── runner.go               # ConversationRunner
│       ├── sse.go                  # SSE 写入
│       ├── hitl.go                 # HITL 处理器（内存 channel）
│       └── hitl_cache.go           # Redis HITL 确认缓存
│
├── internal/                       # ====== 业务核心（DDD 分层） ======
│   ├── domain/                     # 领域逻辑（纯业务规则，零框架依赖）
│   │   ├── privilege/              # ✅ 已实现: domain.go（密码验证、用户校验）
│   │   └── {agent,chat,common,datasource,knowledge,  ← 🟡 这些领域逻辑暂在 usecase
│   │       modelconfig,platform,prompt,rag,semanticmodel}  # 待 Task 8 提取到这里
│   ├── usecase/                    # 用例编排（事务边界，跨 domain 协调）
│   │   ├── privilege_usecase.go    # 权限用例（登录/CURD/ACL/密码）
│   │   ├── platform_usecase.go     # 平台用例（账户/组织/租户）
│   │   ├── data_usecase.go         # ⭐ 数据用例（Agent/数据源/知识/配置 CRUD）
│   │   ├── rag_usecase.go          # RAG 用例
│   │   ├── kg_usecase.go           # KG 用例
│   │   └── agent_usecase.go        # Agent 用例
│   ├── service/                    # 应用服务（对接 handler，格式转换，透传 usecase）
│   │   ├── privilege_service.go
│   │   ├── platform_service.go
│   │   ├── data_service.go         # ← 新增 handler 方法先在这里加签名
│   │   ├── rag_service.go
│   │   ├── kg_service.go
│   │   └── agent_service.go
│   ├── repository/                 # 仓储接口（依赖反转，domain 依赖接口而非实现）
│   │   ├── privilege_repo.go       # 11 个接口
│   │   ├── platform_repo.go        # 7 个接口
│   │   ├── data_repo.go            # 6 个接口（Agent/Knowledge/Datasource/Semantic/Model/Prompt）
│   │   ├── rag_repo.go
│   │   ├── kg_repo.go
│   │   └── agent_memory_repo.go
│   ├── dao/                        # 数据访问实现（Repository 的具体实现）
│   │   ├── db/                     # GORM 实现
│   │   │   ├── privilege_repo.go   # 对应 privilege 接口
│   │   │   ├── platform_repo.go    # 对应 platform 接口
│   │   │   ├── data_repo.go        # ⭐ 对应 data 接口 ← Stub 补全主要改这里
│   │   │   ├── rag_repo.go
│   │   │   ├── kg_repo.go
│   │   │   └── agent_memory_repo.go
│   │   ├── cache/                  # Redis + BigCache
│   │   │   └── privilege_cache.go
│   │   ├── queue/                  # RabbitMQ
│   │   │   ├── producer.go         # ✅ 生产者（已实现）
│   │   │   └── consumer.go         # ✅ 消费者（已实现）
│   │   └── external/               # 外部服务 SDK
│   │       ├── milvus.go           # Milvus 向量操作
│   │       ├── docker_exec.go      # Docker Python 执行器
│   │       ├── oss.go              # 对象存储
│   │       ├── dingtalk.go         # 🟡 钉钉 SDK（待 Task 11）  ← SDK 写这里
│   │       ├── feishu.go           # 🟡 飞书 SDK（待 Task 11）  ← SDK 写这里
│   │       ├── wecom.go            # 🟡 企微 SDK（待 Task 11）  ← SDK 写这里
│   │       └── flink.go            # Flink 集成
│   ├── event/                      # 🟡 领域事件（待 Task 9）
│   │   ├── types.go
│   │   └── handler/                # 按模块拆分事件处理器
│   ├── job/                        # 🟡 定时任务（待 Task 10）
│   │   ├── scheduler.go
│   │   └── jobs/
│   ├── model/                      # 数据模型（Entity / DTO / VO）
│   │   ├── base.go                 # BaseModel（privilege 风格: ID/CreateTime/UpdateTime/...）
│   │   ├── base_models.go          # PlatformBaseModel（platform 风格: creator/updator/...）
│   │   ├── privilege_entity.go     # 12 个 privilege 实体
│   │   ├── privilege_dto.go        # DTO（请求体）
│   │   ├── privilege_vo.go         # VO（响应体）
│   │   ├── platform_entity.go      # 6 个 platform 实体
│   │   ├── platform_dto.go
│   │   ├── data_entity.go          # 14 个 data 实体
│   │   ├── data_dto.go             # 待补全（Agent/Knowledge 等 DTO）
│   │   ├── agent_entity.go         # 4 个 agent 实体
│   │   ├── agent_dto.go            # ChatModelRequest / SSE 事件
│   │   ├── common_entity.go        # PlatformInfo
│   │   ├── rag_entity.go
│   │   └── kg_entity.go
│   └── config/                     # Viper 配置结构体映射
│       ├── app.go, db.go, redis.go, milvus.go, rabbitmq.go
│       ├── agent.go, rpc.go, monitor.go
│       └── casbin_model.conf       # Casbin RBAC 模型
│
├── infra/                          # ====== 基础设施（可跨项目复用） ======
│   ├── logger/                     # Zap 封装（✅ 有测试）
│   ├── config/                     # Viper 加载器（✅ 有测试）
│   ├── response/                   # 统一响应 Response / PageResponse
│   ├── errcode/                    # 统一错误码 ErrCode
│   ├── jwt/                        # JWT 生成/解析
│   ├── pagination/                 # 分页参数解析
│   ├── sse/                        # SSE 流式输出工具
│   ├── validate/                   # 参数校验封装
│   ├── cache/                      # Redis + BigCache 初始化
│   ├── queue/                      # RabbitMQ 连接管理
│   ├── lock/                       # Redis 分布式锁（Lua 原子解锁）
│   ├── monitoring/                 # OpenTelemetry 初始化
│   ├── id/                         # Sonyflake ID 生成器
│   └── utils/                      # 🟢 通用工具函数（仅 .gitkeep）← 通用函数加这里
│
├── configs/                        # ====== YAML 配置文件 ======
│   ├── api/app.yaml                # HTTP 端口/CORS/JWT 密钥
│   ├── agent/app.yaml              # AI 模型配置（provider/model/api_key）
│   ├── rpc/app.yaml                # gRPC 端口
│   ├── job/app.yaml                # 定时任务 cronspec
│   ├── db.yaml                     # PostgreSQL 连接
│   ├── redis.yaml                  # Redis 连接
│   ├── milvus.yaml                 # Milvus 向量库
│   ├── rabbitmq.yaml               # RabbitMQ 消息队列
│   └── monitor.yaml                # OpenTelemetry + 日志级别
│
├── migrations/                     # 🟡 数据库迁移（待 Task 6）
├── scripts/                        # 构建脚本
├── storage/                        # 文件存储（上传头像等）
├── Dockerfile                      # ✅ 多阶段构建
├── docker-compose.yaml             # ✅ 8 个服务编排
├── Makefile                        # ✅ 14 个构建目标
├── go.mod
└── go.sum
```

### 重构代码归属速查表

| 你要做的事 | 代码写在哪里 | 示例 |
|:---|:---|:---|
| 新增 HTTP 接口 | `api/handler/<模块>/` + `api/router.go` | `api/handler/privilege/user.go` |
| 新增业务逻辑 | `internal/usecase/<模块>_usecase.go` | `internal/usecase/privilege_usecase.go` |
| 新增数据库操作（接口） | `internal/repository/<模块>_repo.go` | `internal/repository/privilege_repo.go` |
| 新增数据库操作（实现） | `internal/dao/db/<模块>_repo.go` | `internal/dao/db/privilege_repo.go` |
| 新增缓存操作 | `internal/dao/cache/<模块>_cache.go` | `internal/dao/cache/privilege_cache.go` |
| 新增消息队列操作 | `internal/dao/queue/` | `internal/dao/queue/consumer.go` |
| 新增外部 SDK | `internal/dao/external/<服务>.go` | `internal/dao/external/dingtalk.go` |
| 新增实体/DTO/VO | `internal/model/<模块>_entity.go` 或 `_dto.go` 或 `_vo.go` | `internal/model/privilege_entity.go` |
| 纯业务规则（无框架依赖） | `internal/domain/<模块>/` | `internal/domain/privilege/domain.go` |
| 新增 Agent 工具 | `agent/tools/<类型>/` | `agent/tools/datasource/sql_tool.go` |
| 新增工作流节点 | `agent/workflows/<工作流>/nodes/` | `agent/workflows/nl2sql/nodes/intent.go` |
| 新增 Agent Hook | `agent/hooks/` | `agent/hooks/profile_hook.go` |
| 新增 Agent 拦截器 | `agent/interceptors/` | `agent/interceptors/login_interceptor.go` |
| 新增通用工具函数 | `infra/utils/` | `infra/utils/strings.go` |
| 新增事件处理器 | `internal/event/handler/<模块>/` | `internal/event/handler/privilege/user_created.go` |
| 新增定时任务 | `internal/job/jobs/` | `internal/job/jobs/daily_report.go` |
| 新增 gRPC proto | `rpc/proto/` | `rpc/proto/agent.proto` |
| 新增 gRPC 实现 | `rpc/server/` | `rpc/server/agent_server.go` |
| 新增服务入口 | `cmd/<服务>/main.go` | `cmd/rpc/main.go` |
| 新增配置结构体 | `internal/config/<模块>.go` | `internal/config/db.go` |
| 新增中间件 | `api/middleware/` | `api/middleware/auth.go` |
| 新增基础设施组件 | `infra/<组件>/` | `infra/cache/cache.go` |

### GORM 实体约定

```go
// 每个实体文件头部有 Java 对照注释
// Java: phoenix-privilege-api/.../PrivilegeUser.java  →  Go: internal/model/privilege_entity.go

// Privilege 风格（phoenix-privilege）：
type PrivilegeUser struct {
    BaseModel                    // ID string PK + CreateTime/UpdateTime/CreateBy/UpdateBy/DelFlag int
    Username string `gorm:"column:username"`
}
func (PrivilegeUser) TableName() string { return "tbl_privilege_user" }

// Platform 风格（phoenix-platform / phoenix-common）：
type GroupInfo struct {
    PlatformBaseModel            // CreateTime/creator/UpdateTime/updator/DelFlag/keyword
    ID string `gorm:"column:id;primaryKey"`  // 各实体自声明 PK
    Name string `gorm:"column:name"`
}
func (GroupInfo) TableName() string { return "tbl_platform_group_info" }
```

### Handler 统一模式

```go
// 每个 Handler 文件遵循此模板：
type XxxHandler struct {
    svc *service.XxxService        // 注入 service
}

func NewXxxHandler(svc *service.XxxService) *XxxHandler {
    return &XxxHandler{svc: svc}
}

func (h *XxxHandler) GetByID(c *gin.Context) {
    // 1. 解析参数
    id := c.Param("id")
    // 2. 调用 service
    result, err := h.svc.GetByID(c.Request.Context(), id)
    // 3. 处理错误
    if err != nil {
        if appErr, ok := err.(*usecase.AppError); ok {
            response.Error(c, errcode.ErrCode{Code: appErr.Code, Msg: appErr.Msg})
            return
        }
        response.Error(c, errcode.InternalError)
        return
    }
    // 4. 返回结果
    response.Success(c, result)
}
```

---

## 完整的可执行任务列表

按优先级排序，每个任务独立可执行。使用 `/superpowers:subagent-driven-development` 逐个执行。

### Task 1: 数据 Usecase 补全 + 7 个 Stub Handler 激活

> **优先:** 🔴 最高 | **时间:** 4h | **依赖:** 无

**问题：** 以下 handler 已创建文件但所有方法返回硬编码值：

| Handler | 文件 | 端点 |
|:---|:---|:---|
| ChatSession/ChatMessage | `api/handler/chat/chat.go` | 9 (sessions CRUD, messages, pin, rename, report) |
| Datasource | `api/handler/datasource/datasource.go` | 13 (CRUD, test, tables, columns, logical relations) |
| ModelConfig | `api/handler/modelconfig/model_config.go` | 7 (CRUD, activate, test, check-ready) |
| PromptConfig | `api/handler/prompt/prompt_config.go` | 14 (CRUD, enable/disable, batch, priority, display-order) |
| SemanticModel | `api/handler/semanticmodel/semantic_model.go` | 11 (CRUD, batch, Excel, template) |
| BusinessKnowledge | `api/handler/knowledge/business_knowledge.go` | 8 (CRUD, recall, refresh-vector, retry-embedding) |
| Captcha | `api/handler/privilege/auth.go` | 1 (图生成 + Redis 存储) |

**修复步骤：**

```
Step 1: internal/usecase/data_usecase.go 添加所有缺失方法
  - ListChatSessions, CreateChatSession, DeleteChatSession
  - GetSessionMessages, AddChatMessage
  - PinSession, RenameSession, DeleteSession
  - CreateDatasource, UpdateDatasource, TestDatasourceConnection
  - GetDatasourceTables, GetTableColumns
  - CRUD LogicalRelation
  - CRUD ModelConfig, ActivateModelConfig, TestModelConfig
  - CRUD PromptConfig, EnablePrompt, DisablePrompt, BatchEnable, BatchDisable
  - CRUD SemanticModel, BatchDelete, BatchEnable, BatchDisable
  - CRUD BusinessKnowledge, ToggleRecall, RefreshVectorStore, RetryEmbedding

Step 2: internal/dao/db/data_repo.go 补充 GORM 方法
  - 为新的 usecase 方法添加对应的数据库操作
  - ChatSession: FindByAgentIDAndUserID, FindBySessionID, UpdatePinned, UpdateTitle
  - ChatMessage: FindBySessionID, Create
  - Datasource: FindAll, FindByType, TestConnection
  - 等等...

Step 3: 逐个 handler 替换 stub 为实际调用
  - 注入 *service.DataService
  - 每个方法: parse params → call svc → handle error → response.Success
  - 参考已完成的 api/handler/agent/agent.go 作为模板

Step 4: api/router.go 确保 routes 正确注册

Step 5: go build ./... && go vet ./... → PASS
```

**验证：** 使用 `/superpowers:requesting-code-review` 审查 diff，确保所有端点对齐 Java Controller。

---

### Task 2: Captcha 验证码实现

> **优先:** 🔴 高 | **时间:** 1h | **依赖:** Task 1

**文件:** `api/handler/privilege/auth.go`

**实现：**
```go
import (
    "image"
    "image/color"
    "image/png"
    "bytes"
    "encoding/base64"
    "math/rand"
)

func (h *AuthHandler) Captcha(c *gin.Context) {
    // 1. 生成 4 位随机验证码（字符集: ABCDEFGHJKLMNPQRSTUVWXYZ23456789）
    code := generateRandomCode(4)
    key := uuid.New().String()
    
    // 2. 存 Redis: key="captcha:"+key, TTL=60s
    h.captchaStore.Set(ctx, key, code, 60*time.Second)
    
    // 3. 生成 PNG 图片（200x80, 带干扰线和噪点）
    img := generateCaptchaImage(code)
    
    // 4. Base64 编码返回
    var buf bytes.Buffer
    png.Encode(&buf, img)
    b64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
    
    response.Success(c, model.CaptchaVO{CaptchaKey: key, Image: b64})
}
```

**Go 标准库 `image` 即可生成，无需第三方库。**

---

### Task 3: PlatformSync 实现

> **优先:** 🟡 中 | **时间:** 1h | **依赖:** 无

**文件:** `api/handler/common/platform_sync.go`

替换 3 个 stub 为实际调用：
- `POST /platform/sync/all` → 调用 DingTalk/Feishu/WeCom SDK 全量同步
- `POST /platform/sync/departments` → 同步部门
- `POST /platform/sync/users` → 同步用户

SDK 已 stub 在 `internal/dao/external/`，先调用 stub，后续 Task 13 补全 SDK 后自动生效。

---

### Task 4: KG 工作流节点

> **优先:** 🔴 高 | **时间:** 3h | **依赖:** 无

**文件：** 在 `agent/workflows/kg/nodes/` 下创建：

```
agent/workflows/kg/nodes/
├── entity_extract.go    # LLM 提取知识图谱实体
├── relation_extract.go  # LLM 识别实体间关系
├── graph_merge.go       # 融合新三元组到已有图
└── types.go             # 状态类型定义

agent/workflows/kg/graph.go  # KG 工作流图
```

**节点实现：**
- `EntityExtractNode` — LLM 调用，prompt: "从以下文本中提取实体(type/name/properties)"
- `RelationExtractNode` — LLM 调用，prompt: "识别实体间关系(subject/predicate/object)"
- `GraphMergeNode` — 通过 GORM 写入 `tbl_kg_*` 表

**图结构：** `entity_extract → relation_extract → graph_merge → END`

---

### Task 5: gRPC 服务

> **优先:** 🔴 高 | **时间:** 4h | **依赖:** 无

```
rpc/proto/
├── agent.proto       # AgentService: StreamCall, ListAgents
├── privilege.proto   # PrivilegeService: Login, GetUser
└── data.proto        # DataService: Chat

cmd/rpc/main.go       # gRPC 服务启动
rpc/server/           # proto 生成的 server 实现
rpc/client/           # 客户端封装
```

**步骤：**
1. 安装 protoc + protoc-gen-go + protoc-gen-go-grpc
2. 编写 3 个 `.proto` 文件，定义 service + message
3. `protoc --go_out=. --go-grpc_out=. rpc/proto/*.proto` 生成 Go 代码
4. 实现 server（包裹现有 service/usecase）
5. `cmd/rpc/main.go` 启动 gRPC + reflection（开发模式）

---

### Task 6: 数据库迁移文件

> **优先:** 🔴 高 | **时间:** 0.5h | **依赖:** 无

```bash
cd d:/GitHub/Phoenix-Agent-Java/godata
# 安装 golang-migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 创建初始迁移
migrate create -ext sql -dir migrations -seq init_schema

# 复制 sql/all_schema.sql 内容 → migrations/000001_init_schema.up.sql
# 编写对应的 DROP 语句 → migrations/000001_init_schema.down.sql
```

---

### Task 7: MCP 外部工具

> **优先:** 🟡 中 | **时间:** 2h | **依赖:** 无

**文件:** `agent/tools/external/mcp.go`

实现 MCP 协议客户端，发现远程工具并注册到 tRPC-Agent-Go ToolRegistry：

```go
type MCPToolSource struct {
    servers []MCPServerConfig
}

func (m *MCPToolSource) DiscoverTools(ctx context.Context) ([]tool.Tool, error)
func (m *MCPToolSource) RegisterAll(registry *tools.ToolRegistry) error
```

tRPC-Agent-Go 已有 `tool/mcp/` 包，直接使用。

---

### Task 8: DDD 领域提取

> **优先:** 🟡 中 | **时间:** 3h | **依赖:** 无

将 usecase 中的业务规则提取到 domain 层：

```
internal/usecase/
  → internal/domain/
    ├── privilege/user_rules.go     (密码规则、用户校验)
    ├── privilege/role_rules.go     (角色校验)
    ├── privilege/acl_rules.go      (权限位掩码合并)
    ├── platform/account_rules.go
    ├── platform/group_rules.go
    ├── agent/rules.go              (Agent 类型路由规则)
    ├── datasource/rules.go         (数据源连接校验)
    ├── chat/rules.go               (会话生命周期)
    ├── knowledge/rules.go
    ├── modelconfig/rules.go
    ├── prompt/rules.go
    ├── rag/rules.go
    └── semanticmodel/rules.go
```

每个文件包含纯函数（无框架依赖），useCase 调用 domain：

```go
// current (in usecase):
if privileged.CheckPassword(user.Password, dto.Password) { ... }

// after extraction, same code moves to:
// internal/domain/privilege/user_rules.go
func ValidateLogin(user *model.PrivilegeUser, plainPassword string) error { ... }
```

功能不受影响，`go build ./...` 后应保持 PASS。

---

### Task 9: 事件驱动架构

> **优先:** 🟡 中 | **时间:** 2h | **依赖:** RabbitMQ Consumer（已完成）

**文件：**
```
internal/event/
├── types.go                        # Event 结构体
├── bus.go                          # EventBus 发布/订阅
└── handler/
    ├── privilege/user_created.go   # 用户创建 → 初始化默认角色
    ├── privilege/login_success.go  # 登录成功 → 记录日志
    ├── agent/action_recorded.go    # Agent 调用 → 更新统计
    └── chat/session_created.go     # 会话创建 → 通知
```

**EventBus** 包装已有的 `internal/dao/queue/consumer.go`：
```go
type EventBus struct {
    producer *queue.Producer
    consumer *queue.Consumer
}

func (b *EventBus) Publish(ctx context.Context, event Event) error
func (b *EventBus) Subscribe(eventType string, handler EventHandler)
```

---

### Task 10: 定时任务

> **优先:** 🟡 中 | **时间:** 1h | **依赖:** 无

**文件：**
```
internal/job/
├── scheduler.go                    # Cron 调度器
└── jobs/
    ├── daily_report.go             # 每日报表
    ├── audit_sync.go               # 审计日志同步
    └── session_cleanup.go          # 过期会话清理
```

使用 `github.com/robfig/cron/v3`（已在 go.mod）：

```go
func StartScheduler() *cron.Cron {
    c := cron.New()
    c.AddFunc("0 8 * * *", jobs.DailyReport)
    c.AddFunc("*/30 * * * *", jobs.AuditSync)
    c.AddFunc("@daily", jobs.SessionCleanup)
    c.Start()
    return c
}
```

在 `cmd/api/main.go` 中 `go StartScheduler()` 启动。

---

### Task 11: 三方 SDK 实现

> **优先:** 🟡 中 | **时间:** 3h | **依赖:** 无

**文件：** 替换 `internal/dao/external/` 中的 stub：

```
internal/dao/external/
├── dingtalk.go     # GetUserIDByMobile(phone) → string, SendMessage(chatID, msg) → error
├── feishu.go       # GetUserAccessToken(code) → *Token, GetUserInfo(token) → *User
├── wecom.go        # GetAccessToken() → string, SendMessage(chatID, msg) → error
```

每个 SDK 使用标准 HTTP client + JSON，无需第三方 SDK 包。参考 Java `DingTalkSdkServiceImpl` 等。

---

### Task 12: Agent 工具实现

> **优先:** 🟢 低 | **时间:** 2h | **依赖:** 无

**文件：**
```
agent/tools/web/web_search.go       # DuckDuckGo 搜索（tRPC-Agent-Go 已有 tool/duckduckgo/）
agent/tools/rpc/rpc_tool.go         # 跨服务 RPC 调用工具
agent/tools/privilege/check.go      # 权限检查工具（封装 privilege usecase）
```

所有工具实现 `tool.Tool` 接口（`Declaration()`），注册到 `tools.ToolRegistry`。

---

### Task 13: 清理 .gitkeep + 空目录

> **优先:** 🟡 中 | **时间:** 0.5h | **依赖:** 所有 task 之后

```bash
cd d:/GitHub/Phoenix-Agent-Java/godata
# 删除所有已有 Go 代码的目录中的 .gitkeep
find . -name ".gitkeep" -not -path "./_archive/*" -delete

# 如果 infra/utils 仍为空，决定：填充工具函数 或 删除目录
# 推荐填充：
cat > infra/utils/strings.go << 'EOF'
package utils
func Contains(slice []string, item string) bool { ... }
func Coalesce(vals ...string) string { ... }
EOF

# 验证
go build ./...
```

---

### Task 14: 全量代码审查 + 测试补充

> **优先:** 🟡 中 | **时间:** 6h | **依赖:** Task 1-13

```
Step 1: /superpowers:requesting-code-review   # 全量审查
Step 2: 根据审查结果修复
Step 3: 补充单元测试（至少 core 包）
  - internal/usecase/privilege_usecase_test.go
  - internal/usecase/data_usecase_test.go
  - api/handler/privilege/auth_test.go
Step 4: 补充集成测试
  - Docker Compose 启动所有服务
  - 测试真实的 DB/Redis/RabbitMQ 操作
```

---

## 任务执行顺序

```
Task 1（4h）→ 激活 7 个 stub handler，这是最大的缺口
  ├── Task 2（1h）→ Captcha
  ├── Task 3（1h）→ PlatformSync
  ├── Task 6（0.5h）→ DB 迁移
  └── Task 7（2h）→ MCP 工具

Task 4（3h）→ KG 工作流
Task 5（4h）→ gRPC 服务

Task 8（3h）→ DDD 提取
Task 9（2h）→ 事件驱动
Task 10（1h）→ 定时任务
Task 11（3h）→ 三方 SDK
Task 12（2h）→ Agent 工具

Task 13（0.5h）→ 清理
Task 14（6h）→ 审查 + 测试
```

**总时间：33h**

---

## 使用 Superpowers 恢复工作

新会话中直接读取本文档，然后：

```bash
# 方式1：自动化逐个执行
/superpowers:subagent-driven-development

# 告诉我执行 Task N，例如：
"执行 Task 1：激活 7 个 stub handler"

# 方式2：先写计划再执行
/superpowers:writing-plans
"基于 handoff 文档中的 Task 1，写详细实施计划"
```

### 每个 Task 执行后验证

```bash
cd d:/GitHub/Phoenix-Agent-Java/godata
go build ./... && go vet ./...
```

### 快速命令参考

```bash
# 启动基础设施
docker-compose up -d postgres redis rabbitmq milvus

# 启动 API
go run ./cmd/api              # → :8066
# 健康检查
curl http://localhost:8066/echo
# → {"code":0,"message":"success","data":"ok"}

# 全部测试
go test ./... -short
```

---

*交接完毕。Ciallo～(∠・ω< )⌒☆。*
