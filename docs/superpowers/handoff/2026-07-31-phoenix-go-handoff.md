# Phoenix Go 重写 — 交接与可执行重构计划

> 日期：2026-07-31 | 分支：main | `go build ./...` PASS | `go vet ./...` PASS

## 当前状态

161 个 Go 文件，21,551 行代码。所有 12 个 Java 模块已映射，核心架构齐全。发现 **16 个功能缺口**（7 个 stub handler + 9 个空目录/占位），需要补全。

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
