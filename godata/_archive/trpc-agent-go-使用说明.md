# tRPC-Agent-Go 框架使用说明

> 模块路径: `trpc.group/trpc-go/trpc-agent-go`  
> 语言版本: Go 1.21+  
> 开源协议: Apache 2.0  
> 仓库地址: https://github.com/trpc-group/trpc-agent-go

---

## 目录

1. [框架概述](#1-框架概述)
2. [安装与初始化](#2-安装与初始化)
3. [核心包详解](#3-核心包详解)
   - [3.1 agent — 核心 Agent 接口](#31-agent--核心-agent-接口)
   - [3.2 agent/llmagent — LLM Agent](#32-agentllmagent--llm-agent)
   - [3.3 agent/graphagent — Graph Agent](#33-agentgraphagent--graph-agent)
   - [3.4 agent/chainagent — Chain Agent](#34-agentchainagent--chain-agent)
   - [3.5 agent/parallelagent — Parallel Agent](#35-agentparallelagent--parallel-agent)
   - [3.6 agent/cycleagent — Cycle Agent](#36-agentcycleagent--cycle-agent)
   - [3.7 runner — Agent 运行器](#37-runner--agent-运行器)
   - [3.8 model — LLM 模型接口](#38-model--llm-模型接口)
   - [3.9 tool — 工具系统](#39-tool--工具系统)
   - [3.10 graph — 状态图引擎](#310-graph--状态图引擎)
   - [3.11 memory — 记忆系统](#311-memory--记忆系统)
   - [3.12 session — 会话管理](#312-session--会话管理)
   - [3.13 event — 事件系统](#313-event--事件系统)
   - [3.14 knowledge — 知识库/RAG](#314-knowledge--知识库rag)
   - [3.15 skill — 技能系统](#315-skill--技能系统)
   - [3.16 planner — 规划器](#316-planner--规划器)
   - [3.17 artifact — 产物管理](#317-artifact--产物管理)
   - [3.18 evolution — 自进化](#318-evolution--自进化)
   - [3.19 evaluation — 评估框架](#319-evaluation--评估框架)
   - [3.20 telemetry — 可观测性](#320-telemetry--可观测性)
   - [3.21 server — 服务器框架](#321-server--服务器框架)
   - [3.22 codeexecutor — 代码执行器](#322-codeexecutor--代码执行器)
   - [3.23 prompt — 提示词模板](#323-prompt--提示词模板)
   - [3.24 plugin — 插件系统](#324-plugin--插件系统)
   - [3.25 team — 多 Agent 团队](#325-team--多-agent-团队)
4. [完整示例](#4-完整示例)
   - [4.1 基础 LLM Agent 对话](#41-基础-llm-agent-对话)
   - [4.2 带工具的 Agent](#42-带工具的-agent)
   - [4.3 Graph 工作流 Agent](#43-graph-工作流-agent)
   - [4.4 多 Agent 链式协作](#44-多-agent-链式协作)
   - [4.5 带记忆的 Agent](#45-带记忆的-agent)
   - [4.6 带知识库的 Agent](#46-带知识库的-agent)
   - [4.7 动态 Agent 工厂](#47-动态-agent-工厂)
   - [4.8 运行取消与停止](#48-运行取消与停止)
   - [4.9 AG-UI 服务端](#49-ag-ui-服务端)
   - [4.10 A2A 交互](#410-a2a-交互)

---

## 1. 框架概述

tRPC-Agent-Go 是腾讯 tRPC 团队开源的 **Go 语言 AI Agent 框架**，用于构建生产级的智能 Agent 系统。它提供了 LLM Agent、图工作流、工具调用、会话与记忆状态、知识检索、Agent 自进化、评估和 OpenTelemetry 可观测性于一体。

### 架构设计

```
┌─────────────────────────────────────────────────────────┐
│                    Runner                                │
│  (会话管理、事件循环、运行追踪)                            │
├─────────────────────────────────────────────────────────┤
│                    Agent Interface                       │
│  ┌──────┬───────┬──────┬──────┬──────┬──────┬──────┐   │
│  │LLM   │Graph  │Chain │Cycle │Parallel│A2A  │Codex │   │
│  │Agent │Agent  │Agent │Agent │Agent  │Agent│Agent │   │
│  └──────┴───────┴──────┴──────┴──────┴──────┴──────┘   │
├─────────────────────────────────────────────────────────┤
│  Model      Tool       Memory    Knowledge    Skill     │
│  (LLM       (Function  (事实 &    (RAG/        (SKILL.md │
│   适配器)   工具)      情景记忆)   搜索)       仓库)     │
├─────────────────────────────────────────────────────────┤
│  session    event     graph     telemetry    codeexec   │
│  (历史、    (流、      (状态图    (Otel、      (沙箱、    │
│   摘要、    持久化)    引擎)      Langfuse)    本地)     │
│   分支)                                                │
└─────────────────────────────────────────────────────────┘
```

### 核心特性

- **多 Agent 编排**: ChainAgent（链式）、ParallelAgent（并行）、CycleAgent（循环）
- **GraphAgent**: 类型安全的状态图工作流，支持条件路由（功能等价 LangGraph for Go）
- **丰富工具生态**: Function 工具、MCP 工具、网页搜索、代码执行
- **持久化状态**: 会话、记忆、产物、知识检索
- **Agent 技能**: 可复用的 `SKILL.md` 工作流
- **Agent 自进化**: Hermes 式会话审查，自动提取技能
- **提示词缓存**: 自动成本优化，缓存内容节省 90%
- **评估与基准**: Eval 集 + 指标衡量质量
- **协议集成**: AG-UI、A2A、MCP
- **生产级可观测性**: OpenTelemetry + Langfuse

---

## 2. 安装与初始化

```bash
# 初始化模块
go mod init your-module-name

# 安装依赖
go get trpc.group/trpc-go/trpc-agent-go
```

---

## 3. 核心包详解

### 3.1 `agent` — 核心 Agent 接口

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/agent`

所有 Agent 必须实现的核心接口。

#### Agent 接口

```go
type Agent interface {
    Run(ctx context.Context, invocation *Invocation) (<-chan *event.Event, error)
    Info() Info
    Tools() []tool.Tool
    SubAgents() []agent.Agent
    FindSubAgent(name string) agent.Agent
}
```

#### 关键类型

| 类型 | 说明 |
|------|------|
| `Info` | Agent 元数据（Name, Description, InputSchema, OutputSchema） |
| `Invocation` | 一次 Agent 运行的执行上下文 |
| `RunOptions` | 单次运行配置（model, tools, streaming, timeouts 等） |
| `RunOption` | `RunOptions` 的函数式选项 |
| `StopError` | 停止 Agent 执行的信号 |
| `TransferInfo` | 待处理的 Agent 转移元数据 |
| `StreamMode` | 事件流模式枚举 |

#### 关键函数

```go
func NewStopError(message string) *StopError
func AsStopError(err error) (*StopError, bool)
```

#### RunOption 函数（部分）

```go
agent.WithModel(model model.Model)
agent.WithModelName(name string)
agent.WithStream(enabled bool)
agent.WithInstruction(instruction string)
agent.WithGlobalInstruction(instruction string)
agent.WithToolFilter(filter tool.FilterFunc)
agent.WithAdditionalTools(tools []tool.Tool)
agent.WithExternalTools(tools []tool.Tool)
agent.WithMessages(messages []*model.Message)
agent.WithResume(enabled bool)
agent.WithRequestID(id string)
agent.WithAgent(name string)
agent.WithRuntimeState(key string, value any)
agent.WithKnowledgeFilter(filter *knowledge.SearchFilter)
agent.WithStructuredOutputJSONSchema(schema map[string]any)
agent.WithCodeExecutor(exec codeexecutor.CodeExecutor)
```

---

### 3.2 `agent/llmagent` — LLM Agent

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/agent/llmagent`

基于大语言模型的主要 Agent 实现，支持工具调用、技能、代码执行等。

#### 构造函数

```go
func New(name string, opts ...Option) *LLMAgent
```

**实现**: `agent.Agent`

#### 配置选项

```go
llmagent.WithModel(model.Model)
llmagent.WithModels(map[string]model.Model)
llmagent.WithModelSelector(agent.ModelSelector)
llmagent.WithDescription(description string)
llmagent.WithInstruction(instruction string)
llmagent.WithGlobalInstruction(instruction string)
llmagent.WithGenerationConfig(model.GenerationConfig)
llmagent.WithTools(tools []tool.Tool)
llmagent.WithToolSets(toolSets ...tool.ToolSet)
llmagent.WithSubAgents(agents []agent.Agent)
llmagent.WithCodeExecutor(exec codeexecutor.CodeExecutor)
llmagent.WithPlanner(planner.Planner)
llmagent.WithOutputKey(key string)
llmagent.WithOutputSchema(schema map[string]any)
llmagent.WithStructuredOutput(schema *model.StructuredOutput)
llmagent.WithKnowledge(knowledge knowledge.Knowledge)
llmagent.WithSkills(repo skill.Repository)
llmagent.WithAgentCallbacks(callbacks *agent.Callbacks)
llmagent.WithModelCallbacks(callbacks *model.Callbacks)
llmagent.WithToolCallbacks(callbacks *tool.Callbacks)
llmagent.WithMaxHistoryRuns(n int)
llmagent.WithAddCurrentTime(enabled bool)
llmagent.WithEnableContextCompaction(enabled bool)
llmagent.WithEnableParallelTools(enabled bool)
llmagent.WithReasoningContentMode(mode string)
```

---

### 3.3 `agent/graphagent` — Graph Agent

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/agent/graphagent`

将 `graph.StateGraph` 包装为 `agent.Agent`，使图工作流可作为可运行的 Agent。

#### 构造函数

```go
func New(name string, g *graph.Graph, opts ...Option) (*GraphAgent, error)
```

#### 配置选项

```go
graphagent.WithDescription(desc string)
graphagent.WithInitialState(state graph.State)
```

---

### 3.4 `agent/chainagent` — Chain Agent

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/agent/chainagent`

按顺序执行子 Agent 的链式 Agent。

#### 构造函数

```go
func New(name string, opts ...Option) *ChainAgent
```

#### 配置选项

```go
chainagent.WithSubAgents(agents []agent.Agent)
chainagent.WithDescription(desc string)
```

---

### 3.5 `agent/parallelagent` — Parallel Agent

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/agent/parallelagent`

并发执行子 Agent 并合并结果的并行 Agent。

#### 构造函数

```go
func New(name string, opts ...Option) *ParallelAgent
```

#### 配置选项

```go
parallelagent.WithSubAgents(agents []agent.Agent)
parallelagent.WithDescription(desc string)
```

---

### 3.6 `agent/cycleagent` — Cycle Agent

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/agent/cycleagent`

循环执行规划器+执行器直到停止信号的循环 Agent。

#### 构造函数

```go
func New(name string, opts ...Option) *CycleAgent
```

#### 配置选项

```go
cycleagent.WithPlanner(agent.Agent)
cycleagent.WithExecutor(agent.Agent)
cycleagent.WithMaxIterations(n int)
cycleagent.WithDescription(desc string)
```

---

### 3.7 `runner` — Agent 运行器

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/runner`

执行编排器，管理会话、事件处理、Agent 生命周期和运行追踪。

#### 核心接口

```go
type Runner interface {
    Run(ctx context.Context, userID, sessionID string,
         message *model.Message, opts ...agent.RunOption) (<-chan *event.Event, error)
}

type ManagedRunner interface {
    Runner
    Cancel(requestID string) error
    RunStatus(requestID string) (*RunStatus, error)
}

type SteerableRunner interface {
    ManagedRunner
    EnqueueUserMessage(requestID string, msg *model.Message) error
}
```

#### 构造函数

```go
func NewRunner(appName string, ag agent.Agent, opts ...Option) Runner
func NewRunnerWithAgentFactory(appName, defaultAgentName string,
    factory AgentFactory, opts ...Option) Runner
```

#### 配置选项

```go
runner.WithSessionService(service session.Service)
runner.WithMemoryService(service memory.Service)
runner.WithSessionIngestor(ingestor session.Ingestor)
runner.WithArtifactService(service artifact.Service)
runner.WithEvolutionService(service evolution.Service)
runner.WithAgent(name string, ag agent.Agent)
runner.WithAgentFactory(name string, factory AgentFactory)
runner.WithPlugins(plugins ...plugin.Plugin)
runner.WithAwaitUserReplyRouting(enabled bool)
```

#### 错误常量

```go
runner.ErrRunNotFound
runner.ErrQueuedUserMessageUnsupported
runner.ErrInvalidQueuedUserMessage
```

---

### 3.8 `model` — LLM 模型接口

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/model`

LLM 交互的抽象层。

#### 核心接口

```go
type Model interface {
    GenerateContent(ctx context.Context, request *Request) (<-chan *Response, error)
    Info() Info
}

type IterModel interface {
    Model
    GenerateContentIter(ctx context.Context, request *Request) (Seq[*Response], error)
}
```

#### 关键类型

| 类型 | 说明 |
|------|------|
| `Request` | 完整 LLM 请求（Messages, Tools, GenerationConfig 等） |
| `Response` | LLM 响应（Choices, Usage, Error, Done） |
| `Message` | 消息（Role, Content, ToolCalls, ToolID） |
| `Choice` | 响应选择（Index, Message, Delta, FinishReason） |
| `Usage` | Token 用量统计 |
| `GenerationConfig` | 生成配置（MaxTokens, Temperature, TopP, Stream） |
| `StructuredOutput` | JSON Schema 或格式约束 |
| `TimingInfo` | 首字节和完成时间 |

#### 工厂函数

```go
model.NewUserMessage(content string) *Message
model.NewAssistantMessage(content string) *Message
model.NewSystemMessage(content string) *Message
model.NewToolMessage(content string, toolCallID string) *Message
```

#### 模型实现子包

| 子包 | 说明 |
|------|------|
| `model/openai` | OpenAI 兼容 API |
| `model/anthropic` | Anthropic Claude |
| `model/bedrock` | AWS Bedrock |
| `model/gemini` | Google Gemini |
| `model/hunyuan` | 腾讯混元 |
| `model/ollama` | Ollama 本地模型 |
| `model/huggingface` | HuggingFace 模型 |
| `model/failover` | 故障转移模型路由 |
| `model/hedge` | 对冲/竞速模型策略 |

##### OpenAI 模型示例

```go
openai.New(modelName string, opts ...Option) *Client

// 选项
openai.WithVariant(variant OpenAIVariant) // VariantDeepSeek, VariantOpenAI, VariantMoonshot 等
openai.WithBaseURL(url string)
openai.WithAPIKey(key string)
openai.WithHTTPClient(client *http.Client)
```

---

### 3.9 `tool` — 工具系统

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/tool`

定义 Agent 使用的工具抽象。

#### 核心接口

```go
type Tool interface {
    Declaration() *Declaration
}

type CallableTool interface {
    Tool
    Call(ctx context.Context, jsonArgs []byte) (any, error)
}

type StreamableTool interface {
    Tool
    StreamableCall(ctx context.Context, jsonArgs []byte) (*StreamReader, error)
}
```

#### 关键类型

| 类型 | 说明 |
|------|------|
| `Declaration` | 工具元数据（Name, Description, InputSchema, OutputSchema） |
| `Schema` | JSON Schema 定义 |
| `StreamReader` | 流式响应读取器 |
| `FilterFunc` | `func(ctx context.Context, t Tool) bool` |
| `PermissionPolicy` | 工具级权限检查接口 |
| `ToolSet` | 工具集合接口 |

#### 工具实现子包

| 子包 | 说明 |
|------|------|
| `tool/function` | Go 函数转工具适配器 |
| `tool/file` | 文件操作 |
| `tool/codeexec` | 代码执行 |
| `tool/duckduckgo` | DuckDuckGo 搜索 |
| `tool/google` | Google 搜索 API |
| `tool/webfetch` | 网页内容获取 |
| `tool/wikipedia` | Wikipedia 搜索 |
| `tool/arxivsearch` | ArXiv 学术搜索 |
| `tool/mcp` | MCP 协议工具 |
| `tool/mcpbroker` | MCP 代理 |
| `tool/skill` | 技能管理工具 |
| `tool/agent` | Agent 委托工具 |
| `tool/vision` | 图片分析工具 |
| `tool/openapi` | OpenAPI 工具生成 |
| `tool/email` | 邮件工具 |
| `tool/awaitreply` | 等待用户回复工具 |
| `tool/transfer` | Agent 间转移工具 |
| `tool/hostexec` | 主机命令执行 |
| `tool/workspaceexec` | 工作空间命令执行 |

#### `tool/function` — 函数工具

```go
function.NewFunctionTool(fn any, opts ...Option) tool.CallableTool

// 选项
function.WithName(name string)
function.WithDescription(desc string)
function.WithInputSchema(schema map[string]any)
```

---

### 3.10 `graph` — 状态图引擎

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/graph`

受 LangGraph 启发的状态图执行引擎，支持 DAG 和循环工作流，以及条件路由。

#### StateGraph 方法

```go
// 创建状态图
func NewStateGraph(schema *StateSchema) *StateGraph

// 添加函数节点
func (sg *StateGraph) AddNode(nodeID string, fn NodeFunc) *StateGraph

// 添加 LLM 节点（自动处理与 LLM 的对话）
func (sg *StateGraph) AddLLMNode(nodeID string, m model.Model,
    instruction string, tools map[string]tool.Tool) *StateGraph

// 添加工具节点
func (sg *StateGraph) AddToolsNode(nodeID string,
    tools map[string]tool.Tool) *StateGraph

// 添加 Agent 节点（将任意 Agent 作为图节点）
func (sg *StateGraph) AddAgentNode(nodeID string, ag agent.Agent) *StateGraph

// 添加子图节点
func (sg *StateGraph) AddSubgraphNode(nodeID string, subgraph *Graph) *StateGraph

// 添加普通边
func (sg *StateGraph) AddEdge(from, to string) *StateGraph

// 添加条件边（返回一个目标节点名）
func (sg *StateGraph) AddConditionalEdges(from string,
    cond ConditionalFunc, pathMap map[string]string) *StateGraph

// 添加多条件边（返回多个目标节点名，并行执行）
func (sg *StateGraph) AddMultiConditionalEdges(from string,
    cond MultiConditionalFunc, pathMap map[string]string) *StateGraph

// 工具条件边快捷方法（LLM→工具→后处理的标准模式）
func (sg *StateGraph) AddToolsConditionalEdges(llmNode, toolsNode,
    afterNode string) *StateGraph

// 设置入口点
func (sg *StateGraph) SetEntryPoint(nodeID string) *StateGraph

// 设置结束点
func (sg *StateGraph) SetFinishPoint(nodeID string) *StateGraph

// 设置条件入口点
func (sg *StateGraph) SetConditionalEntryPoint(cond ConditionalFunc,
    pathMap map[string]string) *StateGraph

// 编译图（返回不可变更的 Graph 实例）
func (sg *StateGraph) Compile() (*Graph, error)
```

#### 节点函数签名

```go
type NodeFunc func(ctx context.Context, state State) (any, error)

// 条件函数
type ConditionalFunc func(ctx context.Context, state State) (string, error)

// 多条件函数（返回多个目标，并行执行）
type MultiConditionalFunc func(ctx context.Context, state State) ([]string, error)
```

#### 关键常量

```go
graph.Start = "__start__"   // 起始节点
graph.End = "__end__"       // 结束节点

// 状态键
graph.StateKeyUserInput    // 用户输入
graph.StateKeyMessages     // 消息列表
graph.StateKeyLastResponse // 最后响应
graph.StateKeyNodeResponses // 节点响应

// 元数据键
graph.MetadataKeyNode  // 当前节点
graph.MetadataKeyTool  // 当前工具
graph.MetadataKeyModel // 当前模型
```

#### 命令式路由

```go
type Command struct {
    Update State    // 状态更新
    GoTo   string   // 下一个节点
    Resume State    // 恢复状态
}
```

---

### 3.11 `memory` — 记忆系统

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/memory`

长期记忆接口，用于持久化和检索用户信息。

#### 核心接口

```go
type Reader interface {
    ReadMemories(ctx context.Context, userKey UserKey, limit int) ([]*Entry, error)
    SearchMemories(ctx context.Context, userKey UserKey,
        query string, opts ...SearchOption) ([]*Entry, error)
}

type Service interface {
    Reader
    AddMemory(ctx context.Context, userKey UserKey, memory string,
        topics []string, opts ...AddOption) error
    UpdateMemory(ctx context.Context, memoryKey Key, memory string,
        topics []string, opts ...UpdateOption) error
    DeleteMemory(ctx context.Context, memoryKey Key) error
    ClearMemories(ctx context.Context, userKey UserKey) error
    Tools() []tool.Tool
    EnqueueAutoMemoryJob(ctx context.Context, sess *session.Session) error
    Close() error
}
```

#### 关键类型

| 类型 | 说明 |
|------|------|
| `Entry` | 存储的记忆（ID, 内容, 分数） |
| `Kind` | `"fact"` 或 `"episode"` |
| `Key` | 复合键 `{AppName, UserID, MemoryID}` |
| `UserKey` | `{AppName, UserID}` |
| `SearchOptions` | 高级搜索（Kind, TimeRange, HybridSearch, Deduplicate） |

#### 记忆实现子包

| 子包 | 说明 |
|------|------|
| `memory/inmemory` | 内存（开发/测试用） |
| `memory/sqlite` | SQLite 存储 |
| `memory/sqlitevec` | SQLite + 向量搜索 |
| `memory/postgres` | PostgreSQL |
| `memory/pgvector` | PostgreSQL + pgvector |
| `memory/mysql` | MySQL |
| `memory/mysqlvec` | MySQL + 向量搜索 |
| `memory/redis` | Redis |
| `memory/chromadb` | ChromaDB |
| `memory/mem0` | Mem0 平台 |

---

### 3.12 `session` — 会话管理

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/session`

管理对话会话，包括事件存储、状态、摘要、分支和搜索。

#### 核心接口

```go
type Service interface {
    GetSession(ctx context.Context, key Key) (*Session, error)
    UpsertSession(ctx context.Context, session *Session) error
    DeleteSession(ctx context.Context, key Key) error
    ListSessions(ctx context.Context, appName, userID string,
        opts ...Options) ([]*Session, error)
}

type SearchableService interface {
    Service
    SearchEvents(ctx context.Context, request *EventSearchRequest) (*EventSearchResult, error)
}

type WindowService interface {
    Service
    GetEventWindow(ctx context.Context, key Key, start, end int) ([]event.Event, error)
}

type Ingestor interface {
    Ingest(ctx context.Context, session *Session) error
}
```

#### 关键类型

| 类型 | 说明 |
|------|------|
| `Session` | 核心会话（Events, State, Tracks, Summaries） |
| `Key` | `{AppName, UserID, SessionID}` |
| `StateMap` | `map[string][]byte` |
| `Summary` | 对话摘要 + 边界元数据 |
| `Track` / `TrackEvents` | 多 Agent 对话的分支追踪 |
| `EventSearchRequest` | 会话事件的语义搜索请求 |
| `EventSearchResult` | 搜索结果（命中 + 分数 + 元数据） |

#### Session 方法

```go
func (s *Session) GetEvents() []event.Event
func (s *Session) GetState(key string) ([]byte, bool)
func (s *Session) SetState(key string, value []byte)
func (s *Session) SnapshotState() StateMap
func (s *Session) Clone() *Session
```

#### 会话实现子包

| 子包 | 说明 |
|------|------|
| `session/inmemory` | 内存 |
| `session/summary` | 摘要生成 |
| `session/noop` | 空操作 |
| `session/sqlite` | SQLite |
| `session/redis` | Redis |
| `session/mysql` | MySQL |
| `session/tdsql` | TDSQL |
| `session/postgres` | PostgreSQL |
| `session/pgvector` | PGVector |
| `session/clickhouse` | ClickHouse |
| `session/mongodb` | MongoDB |

---

### 3.13 `event` — 事件系统

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/event`

定义 Agent 到调用方通信、会话持久化和流式传输的事件模型。

#### 关键类型

```go
type Event struct {
    Response          *model.Response
    RequestID         string
    InvocationID      string
    Author            string
    Branch            string
    FilterKey         string
    StateDelta        map[string]any
    Extensions        map[string]any
    Error             error
    Done              bool
    SkipSummarization bool
}
```

#### 关键函数

```go
func New(invocationID, author string, opts ...Option) *Event
func NewErrorEvent(invocationID, author, errorType, errorMessage string, opts ...Option) *Event
func NewResponseEvent(invocationID, author string, response *model.Response, opts ...Option) *Event
func EmitEvent(ctx context.Context, ch chan<- *Event, e *Event) error
func EmitEventWithTimeout(ctx context.Context, ch chan<- *Event, e *Event, timeout time.Duration) error
```

#### Event 方法

```go
func (e *Event) IsRunnerCompletion() bool
func (e *Event) IsError() bool
func (e *Event) IsTerminalError() bool
func (e *Event) ContainsTag(tag string) bool
func (e *Event) Filter(filterKey string) bool
func (e *Event) Clone() *Event
```

---

### 3.14 `knowledge` — 知识库/RAG

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/knowledge`

检索增强生成（RAG）能力接口。

#### 核心接口

```go
type Knowledge interface {
    Search(ctx context.Context, req *SearchRequest) (*SearchResult, error)
}
```

#### 关键类型

| 类型 | 说明 |
|------|------|
| `SearchRequest` | 查询（Query, History, UserID, MaxResults, MinScore, SearchFilter, SearchMode） |
| `SearchResult` | 最佳 Document + 顶部 N 个 Results |
| `Result` | Document + Score |
| `SearchFilter` | DocumentIDs, Metadata, FilterCondition |

---

### 3.15 `skill` — 技能系统

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/skill`

与模型无关的技能仓库系统。技能是包含 `SKILL.md` 文件的文件夹。

#### 核心接口

```go
type Repository interface {
    Summaries() []Summary
    Get(name string) (*Skill, error)
    Path(name string) (string, error)
}

type RootedRepository interface {
    Repository
    Roots() []string
}

type RefreshableRepository interface {
    Repository
    Refresh() error
}
```

#### 关键类型

| 类型 | 说明 |
|------|------|
| `Summary` | 技能名称 + 描述 |
| `Skill` | 完整技能（Summary + Body + Docs） |
| `Doc` | 辅助文档（Path + Content） |
| `FSRepository` | 文件系统实现 |

#### FSRepository 方法

```go
func NewFSRepository(roots ...string) (*FSRepository, error)
func (r *FSRepository) Summaries() []Summary
func (r *FSRepository) Get(name string) (*Skill, error)
func (r *FSRepository) Path(name string) (string, error)
func (r *FSRepository) Refresh() error
```

---

### 3.16 `planner` — 规划器

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/planner`

Agent 规划和推理能力的接口。

```go
type Planner interface {
    Plan(ctx context.Context, request PlanRequest) (*Plan, error)
}
```

---

### 3.17 `artifact` — 产物管理

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/artifact`

管理 Agent 产生的版本化文件/输出。

```go
type Service interface {
    Create(ctx context.Context, artifact *Artifact) (*Artifact, error)
    Get(ctx context.Context, id string) (*Artifact, error)
    List(ctx context.Context, sessionID string) ([]*Artifact, error)
    Update(ctx context.Context, artifact *Artifact) (*Artifact, error)
    Delete(ctx context.Context, id string) error
}
```

---

### 3.18 `evolution` — 自进化

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/evolution`

审查已完成的会话并将可复用的流程提取为可管理的 Agent 技能。

---

### 3.19 `evaluation` — 评估框架

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/evaluation`

使用可重复的评估集和可插拔指标评估 Agent。

```go
evaluator, err := evaluation.New("app", runner, evaluation.WithNumRuns(3))
result, err := evaluator.Evaluate(ctx, "eval-set-name")
_ = result.OverallStatus
```

---

### 3.20 `telemetry` — 可观测性

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/telemetry`

基于 OpenTelemetry 的可观测性基础设施。

#### 子包

| 子包 | 说明 |
|------|------|
| `telemetry/appid` | 应用身份注册 |
| `telemetry/errs` | 错误分类 |
| `telemetry/langfuse` | Langfuse 追踪集成 |
| `telemetry/metric` | 指标（OpenTelemetry + Prometheus） |
| `telemetry/semconv` | 语义约定 |
| `telemetry/trace` | 追踪工具 |

##### Langfuse 集成

```go
clean, err := langfuse.Start(ctx)
defer clean(ctx)

// 在 RunOption 中添加 Langfuse 属性
agent.WithSpanAttributes(
    attribute.String("langfuse.user.id", userID),
    attribute.String("langfuse.session.id", sessionID),
)
```

---

### 3.21 `server` — 服务器框架

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/server`

将 Agent 作为网络服务运行的基础设施。

#### 子包

| 子包 | 说明 |
|------|------|
| `server/a2a` | A2A（Agent-to-Agent）协议服务器 |
| `server/openai` | OpenAI 兼容 API 服务器 |
| `server/agui` | AG-UI 服务器 |

##### AG-UI 服务器

提供 Server-Sent Events (SSE) 端点，支持 CopilotKit 和 TDesign Chat 等前端。

---

### 3.22 `codeexecutor` — 代码执行器

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/codeexecutor`

代码执行抽象层，支持本地、沙箱和工作空间执行。

---

### 3.23 `prompt` — 提示词模板

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/prompt`

提示词模板/文本管理。

---

### 3.24 `plugin` — 插件系统

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/plugin`

Runner 级别的插件系统。

---

### 3.25 `team` — 多 Agent 团队

**导入路径**: `trpc.group/trpc-go/trpc-agent-go/team`

多 Agent 团队编排。

---

## 4. 完整示例

### 4.1 基础 LLM Agent 对话

```go
package main

import (
    "context"
    "fmt"
    "log"

    "trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
    "trpc.group/trpc-go/trpc-agent-go/model"
    "trpc.group/trpc-go/trpc-agent-go/model/openai"
    "trpc.group/trpc-go/trpc-agent-go/runner"
)

func main() {
    ctx := context.Background()

    // 1. 创建模型
    m := openai.New("gpt-4o-mini")

    // 2. 创建生成配置
    genConfig := model.GenerationConfig{
        MaxTokens:   intPtr(1000),
        Temperature: floatPtr(0.7),
        Stream:      true,
    }

    // 3. 创建 LLM Agent
    agent := llmagent.New("assistant",
        llmagent.WithModel(m),
        llmagent.WithDescription("一个友好的 AI 助手"),
        llmagent.WithInstruction("你是一个友好的 AI 助手，请用中文回答。"),
        llmagent.WithGenerationConfig(genConfig),
    )

    // 4. 创建 Runner
    r := runner.NewRunner("chat-demo", agent)

    // 5. 运行对话
    events, err := r.Run(ctx, "user-001", "session-001",
        model.NewUserMessage("你好！请介绍一下你自己。"))
    if err != nil {
        log.Fatal(err)
    }

    // 6. 处理事件流
    for event := range events {
        if event.Error != nil {
            log.Printf("错误: %v", event.Error)
            break
        }
        if len(event.Response.Choices) > 0 {
            content := event.Response.Choices[0].Delta.Content
            fmt.Print(content)
        }
        if event.Done {
            break
        }
    }
    fmt.Println()
}

func intPtr(v int) *int { return &v }

func floatPtr(v float64) *float64 { return &v }
```

### 4.2 带工具的 Agent

```go
package main

import (
    "context"
    "fmt"
    "log"

    "trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
    "trpc.group/trpc-go/trpc-agent-go/model"
    "trpc.group/trpc-go/trpc-agent-go/model/openai"
    "trpc.group/trpc-go/trpc-agent-go/runner"
    "trpc.group/trpc-go/trpc-agent-go/tool"
    "trpc.group/trpc-go/trpc-agent-go/tool/function"
    "trpc.group/trpc-go/trpc-agent-go/tool/duckduckgo"
)

type calculatorReq struct {
    A  float64 `json:"A" jsonschema:"description=第一个数字,required"`
    B  float64 `json:"B" jsonschema:"description=第二个数字,required"`
    Op string  `json:"Op" jsonschema:"description=运算类型,enum=add,enum=sub,enum=mul,enum=div,required"`
}

type calculatorRsp struct {
    Result float64 `json:"result"`
}

func calculator(ctx context.Context, req calculatorReq) (calculatorRsp, error) {
    var result float64
    switch req.Op {
    case "add", "+":
        result = req.A + req.B
    case "sub", "-":
        result = req.A - req.B
    case "mul", "*":
        result = req.A * req.B
    case "div", "/":
        if req.B == 0 {
            return calculatorRsp{}, fmt.Errorf("除数不能为0")
        }
        result = req.A / req.B
    default:
        return calculatorRsp{}, fmt.Errorf("不支持的运算: %s", req.Op)
    }
    return calculatorRsp{Result: result}, nil
}

func main() {
    ctx := context.Background()

    m := openai.New("gpt-4o-mini",
        openai.WithVariant(openai.VariantDeepSeek),
    )

    // 创建计算器工具
    calculatorTool := function.NewFunctionTool(
        calculator,
        function.WithName("calculator"),
        function.WithDescription("执行加减乘除四则运算。"+
            "参数: a, b 是数值, op 取值 add/sub/mul/div; 返回计算结果。"),
    )

    // 创建网页搜索工具（DuckDuckGo）
    searchTool := duckduckgo.NewWebFetchTool()

    agent := llmagent.New("assistant",
        llmagent.WithModel(m),
        llmagent.WithTools([]tool.Tool{calculatorTool, searchTool}),
        llmagent.WithInstruction("你是一个有用的助手，可以使用工具来回答问题。请用中文回答。"),
    )

    r := runner.NewRunner("tool-demo", agent)
    genConfig := model.GenerationConfig{Stream: true}
    m2 := openai.New("gpt-4o-mini")

    agent2 := llmagent.New("assistant",
        llmagent.WithModel(m2),
        llmagent.WithGenerationConfig(genConfig),
        llmagent.WithInstruction("请将用户的输入翻译成英文。"),
    )

    // 链式: 先用工具 Agent，再用翻译 Agent
    pipeline := chainagent.New("pipeline",
        chainagent.WithSubAgents([]agent.Agent{agent, agent2}),
    )

    r2 := runner.NewRunner("chain-demo", pipeline)
    events, err := r2.Run(ctx, "user-001", "session-002",
        model.NewUserMessage("计算 25 * 4 等于多少，然后翻译结果"))
    if err != nil {
        log.Fatal(err)
    }

    for event := range events {
        if event.Error != nil {
            log.Printf("错误: %v", event.Error)
            break
        }
        if len(event.Response.Choices) > 0 {
            fmt.Print(event.Response.Choices[0].Delta.Content)
        }
        if event.Done {
            break
        }
    }
    fmt.Println()
}
```

### 4.3 Graph 工作流 Agent

```go
package main

import (
    "context"
    "fmt"
    "log"
    "strings"

    "trpc.group/trpc-go/trpc-agent-go/agent/graphagent"
    "trpc.group/trpc-go/trpc-agent-go/graph"
    "trpc.group/trpc-go/trpc-agent-go/model"
    "trpc.group/trpc-go/trpc-agent-go/model/openai"
    "trpc.group/trpc-go/trpc-agent-go/runner"
    "trpc.group/trpc-go/trpc-agent-go/tool"
    "trpc.group/trpc-go/trpc-agent-go/tool/function"
)

type analyzeReq struct {
    Text string `json:"text" jsonschema:"description=需要分析的文本,required"`
}

type analyzeRsp struct {
    WordCount int `json:"word_count"`
    CharCount int `json:"char_count"`
}

func textAnalyzer(ctx context.Context, req analyzeReq) (analyzeRsp, error) {
    return analyzeRsp{
        WordCount: len(strings.Fields(req.Text)),
        CharCount: len(req.Text),
    }, nil
}

func main() {
    ctx := context.Background()

    // 1. 创建状态 Schema
    schema := graph.MessagesStateSchema()

    // 2. 创建状态图构建器
    sg := graph.NewStateGraph(schema)

    // 3. 添加节点
    // 3a. 函数节点 — 文本分析
    sg.AddNode("analyze", func(ctx context.Context, state graph.State) (any, error) {
        input, _ := state[graph.StateKeyUserInput].(string)
        words := len(strings.Fields(input))
        chars := len(input)
        return graph.State{
            "word_count": words,
            "char_count": chars,
            "analysis":   fmt.Sprintf("文本分析结果: %d 个词, %d 个字符", words, chars),
        }, nil
    })

    // 3b. LLM 节点
    m := openai.New("gpt-4o-mini")
    analyzerTool := function.NewFunctionTool(textAnalyzer,
        function.WithName("text_analyzer"),
        function.WithDescription("分析文本的统计信息（词数、字符数）"),
    )

    sg.AddLLMNode("enhance", m,
        "你是一个文本增强助手。基于分析结果，给出改进建议。",
        map[string]tool.Tool{"text_analyzer": analyzerTool},
    )

    // 3c. 条件路由节点
    sg.AddNode("route", func(ctx context.Context, state graph.State) (any, error) {
        words, _ := state["word_count"].(int)
        if words < 10 {
            return map[string]string{"next": "short_reply"}, nil
        }
        return map[string]string{"next": "long_reply"}, nil
    })

    sg.AddNode("short_reply", func(ctx context.Context, state graph.State) (any, error) {
        return graph.State{"reply": "文本较短，无需过多处理。"}, nil
    })

    sg.AddNode("long_reply", func(ctx context.Context, state graph.State) (any, error) {
        return graph.State{"reply": "文本较长，已进行详细分析。"}, nil
    })

    // 4. 设置边
    sg.SetEntryPoint("analyze")
    sg.AddEdge("analyze", "enhance")
    sg.AddEdge("enhance", "route")
    sg.AddConditionalEdges("route",
        func(ctx context.Context, state graph.State) (string, error) {
            next, _ := state["next"].(string)
            return next, nil
        },
        map[string]string{
            "short_reply": "short_reply",
            "long_reply":  "long_reply",
        },
    )
    sg.SetFinishPoint("short_reply").SetFinishPoint("long_reply")

    // 5. 编译图
    g, err := sg.Compile()
    if err != nil {
        log.Fatal(err)
    }

    // 6. 创建 Graph Agent
    graphAgent, err := graphagent.New("text-processor", g)
    if err != nil {
        log.Fatal(err)
    }

    // 7. 通过 Runner 运行
    r := runner.NewRunner("graph-demo", graphAgent)
    events, err := r.Run(ctx, "user-001", "session-001",
        model.NewUserMessage("你好，今天天气真不错，适合出去散步。"))
    if err != nil {
        log.Fatal(err)
    }

    for event := range events {
        if event.Error != nil {
            log.Printf("错误: %v", event.Error)
            break
        }
        if len(event.Response.Choices) > 0 {
            fmt.Print(event.Response.Choices[0].Delta.Content)
        }
        if event.Done {
            break
        }
    }
    fmt.Println()
}
```

### 4.4 多 Agent 链式协作

```go
package main

import (
    "context"
    "fmt"
    "log"

    "trpc.group/trpc-go/trpc-agent-go/agent"
    "trpc.group/trpc-go/trpc-agent-go/agent/chainagent"
    "trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
    "trpc.group/trpc-go/trpc-agent-go/agent/parallelagent"
    "trpc.group/trpc-go/trpc-agent-go/model"
    "trpc.group/trpc-go/trpc-agent-go/model/openai"
    "trpc.group/trpc-go/trpc-agent-go/runner"
)

func main() {
    ctx := context.Background()

    m := openai.New("gpt-4o-mini")

    // 1. 创建三个专用 Agent
    planner := llmagent.New("planner",
        llmagent.WithModel(m),
        llmagent.WithInstruction("你是一个规划专家。请将用户的请求分解为详细的执行步骤。"),
    )

    researcher := llmagent.New("researcher",
        llmagent.WithModel(m),
        llmagent.WithInstruction("你是一个研究专家。请基于规划步骤进行深入研究，提供详细信息。"),
    )

    writer := llmagent.New("writer",
        llmagent.WithModel(m),
        llmagent.WithInstruction("你是一个写作专家。请基于研究结果撰写最终报告。"),
    )

    // 2. 创建并行 Agent（研究阶段并行执行）
    parallel := parallelagent.New("research-phase",
        parallelagent.WithSubAgents([]agent.Agent{researcher, writer}),
    )

    // 3. 创建链式 Agent（规划→并行执行→写作）
    pipeline := chainagent.New("full-pipeline",
        chainagent.WithSubAgents([]agent.Agent{planner, parallel}),
    )

    // 4. 运行
    r := runner.NewRunner("multi-agent-demo", pipeline)
    events, err := r.Run(ctx, "user-001", "session-001",
        model.NewUserMessage("请为我制定一个学习 Go 语言的三个月计划"))
    if err != nil {
        log.Fatal(err)
    }

    for event := range events {
        if event.Error != nil {
            log.Printf("错误: %v", event.Error)
            break
        }
        if len(event.Response.Choices) > 0 {
            fmt.Print(event.Response.Choices[0].Delta.Content)
        }
        if event.Done {
            break
        }
    }
    fmt.Println()
}
```

### 4.5 带记忆的 Agent

```go
package main

import (
    "context"
    "fmt"
    "log"

    "trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
    "trpc.group/trpc-go/trpc-agent-go/memory/inmemory"
    "trpc.group/trpc-go/trpc-agent-go/model"
    "trpc.group/trpc-go/trpc-agent-go/model/openai"
    "trpc.group/trpc-go/trpc-agent-go/runner"
    "trpc.group/trpc-go/trpc-agent-go/session/inmemory"
    sessioninmemory "trpc.group/trpc-go/trpc-agent-go/session/inmemory"
)

func main() {
    ctx := context.Background()

    m := openai.New("gpt-4o-mini")

    // 1. 创建记忆服务
    memoryService := inmemory.NewService()

    // 2. 创建 Agent（注入记忆工具）
    agent := llmagent.New("memory-assistant",
        llmagent.WithModel(m),
        llmagent.WithTools(memoryService.Tools()),
        llmagent.WithInstruction("你是一个有记忆的助手。记住用户的偏好和信息，在对话中利用这些信息。"+
            "使用 memory_add 工具记住重要信息，使用 memory_search 工具检索记忆。"),
    )

    // 3. 创建 Runner（注入记忆服务和会话服务）
    r := runner.NewRunner("memory-demo", agent,
        runner.WithMemoryService(memoryService),
        runner.WithSessionService(sessioninmemory.NewSessionService()),
    )

    // 第一轮对话：让 Agent 记住用户信息
    events, err := r.Run(ctx, "user-001", "session-001",
        model.NewUserMessage("我叫小明，我喜欢编程和阅读科幻小说。"))
    if err != nil {
        log.Fatal(err)
    }
    for event := range events {
        if len(event.Response.Choices) > 0 {
            fmt.Print(event.Response.Choices[0].Delta.Content)
        }
        if event.Done {
            break
        }
    }
    fmt.Println()

    // 第二轮对话：Agent 应该记得用户信息
    events, err = r.Run(ctx, "user-001", "session-002",
        model.NewUserMessage("你还记得我的名字和爱好吗？"))
    if err != nil {
        log.Fatal(err)
    }
    for event := range events {
        if len(event.Response.Choices) > 0 {
            fmt.Print(event.Response.Choices[0].Delta.Content)
        }
        if event.Done {
            break
        }
    }
    fmt.Println()
}
```

### 4.6 带知识库的 Agent

```go
package main

import (
    "context"
    "fmt"
    "log"

    "trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
    "trpc.group/trpc-go/trpc-agent-go/knowledge"
    "trpc.group/trpc-go/trpc-agent-go/model"
    "trpc.group/trpc-go/trpc-agent-go/model/openai"
    "trpc.group/trpc-go/trpc-agent-go/runner"
)

func main() {
    ctx := context.Background()

    m := openai.New("gpt-4o-mini")

    // 1. 创建知识库（需先创建 Embedder 和 VectorStore）
    // 此处省略知识库的具体创建细节，见 knowledge 子包文档
    var knowledgeBase knowledge.Knowledge

    // 2. 创建 Agent（注入知识库）
    agent := llmagent.New("knowledge-agent",
        llmagent.WithModel(m),
        llmagent.WithKnowledge(knowledgeBase),
        llmagent.WithInstruction("你是一个知识库助手。利用提供的知识来回答用户问题。"+
            "如果你不确定答案，请明确说明。"),
    )

    r := runner.NewRunner("knowledge-demo", agent)
    events, err := r.Run(ctx, "user-001", "session-001",
        model.NewUserMessage("请根据知识库中的数据回答：公司去年的营收情况如何？"))
    if err != nil {
        log.Fatal(err)
    }

    for event := range events {
        if event.Error != nil {
            log.Printf("错误: %v", event.Error)
            break
        }
        if len(event.Response.Choices) > 0 {
            fmt.Print(event.Response.Choices[0].Delta.Content)
        }
        if event.Done {
            break
        }
    }
    fmt.Println()
}
```

### 4.7 动态 Agent 工厂

```go
package main

import (
    "context"
    "fmt"
    "log"

    "trpc.group/trpc-go/trpc-agent-go/agent"
    "trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
    "trpc.group/trpc-go/trpc-agent-go/model"
    "trpc.group/trpc-go/trpc-agent-go/model/openai"
    "trpc.group/trpc-go/trpc-agent-go/runner"
)

func main() {
    ctx := context.Background()

    // 使用 AgentFactory 为每次请求创建不同的 Agent
    r := runner.NewRunnerWithAgentFactory(
        "dynamic-app",
        "assistant",
        func(ctx context.Context, ro agent.RunOptions) (agent.Agent, error) {
            a := llmagent.New("assistant",
                llmagent.WithModel(openai.New("gpt-4o-mini")),
                llmagent.WithInstruction(ro.Instruction),
            )
            return a, nil
        },
    )

    events, err := r.Run(ctx, "user-001", "session-001",
        model.NewUserMessage("你好，请帮我制定一个健身计划。"),
        agent.WithInstruction("你是一个专业的健身教练，请用中文回答。"),
    )
    if err != nil {
        log.Fatal(err)
    }

    for event := range events {
        if event.Error != nil {
            log.Printf("错误: %v", event.Error)
            break
        }
        if len(event.Response.Choices) > 0 {
            fmt.Print(event.Response.Choices[0].Delta.Content)
        }
        if event.Done {
            break
        }
    }
    fmt.Println()
}
```

### 4.8 运行取消与停止

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "time"

    "trpc.group/trpc-go/trpc-agent-go/agent"
    "trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
    "trpc.group/trpc-go/trpc-agent-go/model"
    "trpc.group/trpc-go/trpc-agent-go/model/openai"
    "trpc.group/trpc-go/trpc-agent-go/runner"
)

// 方式 A: Ctrl+C 取消
func cancelWithSignal() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    m := openai.New("gpt-4o-mini")
    agent := llmagent.New("assistant", llmagent.WithModel(m))
    r := runner.NewRunner("cancel-demo", agent)

    events, err := r.Run(ctx, "user-001", "session-001",
        model.NewUserMessage("请写一篇 5000 字的文章"),
    )
    if err != nil {
        log.Fatal(err)
    }
    for range events {
        // 持续消费直到结束或取消
    }
}

// 方式 B: 代码中取消
func cancelInCode() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    m := openai.New("gpt-4o-mini")
    agent := llmagent.New("assistant", llmagent.WithModel(m))
    r := runner.NewRunner("cancel-demo", agent)

    events, err := r.Run(ctx, "user-001", "session-001",
        model.NewUserMessage("请写一篇 5000 字的文章"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 2 秒后取消
    go func() {
        time.Sleep(2 * time.Second)
        cancel()
    }()

    for range events {
        // 持续消费直到通道关闭
    }
}

// 方式 C: 通过 requestID 取消（服务端/后台运行）
func cancelByRequestID() {
    ctx := context.Background()

    m := openai.New("gpt-4o-mini")
    agent := llmagent.New("assistant", llmagent.WithModel(m))

    // 需要 ManagedRunner 来支持 Cancel
    r := runner.NewRunner("cancel-demo", agent).(runner.ManagedRunner)

    requestID := "req-123"
    events, err := r.Run(ctx, "user-001", "session-001",
        model.NewUserMessage("请写一篇 5000 字的文章"),
        agent.WithRequestID(requestID),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 在另一处取消
    err = r.Cancel(requestID)
    if err != nil {
        log.Printf("取消失败: %v", err)
    }

    for range events {
    }
}

func main() {
    // 选择一种方式运行
    cancelWithSignal()
}
```

### 4.9 AG-UI 服务端

```go
package main

import (
    "context"
    "log"
    "net/http"

    "trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
    "trpc.group/trpc-go/trpc-agent-go/model/openai"
    "trpc.group/trpc-go/trpc-agent-go/runner"
)

func main() {
    ctx := context.Background()

    m := openai.New("gpt-4o-mini")
    agent := llmagent.New("assistant",
        llmagent.WithModel(m),
        llmagent.WithInstruction("你是一个友好的 AI 助手。"),
    )

    r := runner.NewRunner("agui-demo", agent)

    // 启动 AG-UI SSE 服务器
    // 这需要导入 server/agui 包并设置路由
    // 完整实现请参考 examples/agui

    log.Println("AG-UI 服务启动在 :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))

    _ = ctx
    _ = r
}
```

### 4.10 A2A 交互

```go
package main

import (
    "context"
    "log"

    "trpc.group/trpc-go/trpc-agent-go/agent/a2aagent"
    "trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
    "trpc.group/trpc-go/trpc-agent-go/model"
    "trpc.group/trpc-go/trpc-agent-go/model/openai"
    "trpc.group/trpc-go/trpc-agent-go/runner"
)

func main() {
    ctx := context.Background()

    m := openai.New("gpt-4o-mini")

    // 创建 A2A Agent（与另一个 Agent 服务通信）
    a2aAgent := a2aagent.New("remote-agent",
        a2aagent.WithRemoteURL("http://localhost:9091/a2a"),
        a2aagent.WithModel(openai.New("gpt-4o-mini")),
    )

    // 也可以使用本地 Agent 作为 fallback
    localAgent := llmagent.New("local-agent",
        llmagent.WithModel(m),
        llmagent.WithInstruction("你是本地助手。"),
    )

    r := runner.NewRunner("a2a-demo", localAgent,
        runner.WithAgent("remote", a2aAgent),
    )

    events, err := r.Run(ctx, "user-001", "session-001",
        model.NewUserMessage("请远程 Agent 帮我查询一些信息"),
        // 通过 RunOption 指定使用远程 Agent
        // agent.WithAgentByName("remote"),
    )
    if err != nil {
        log.Fatal(err)
    }

    for event := range events {
        if event.Error != nil {
            log.Printf("错误: %v", event.Error)
            break
        }
        _ = event
    }
}
```

---

## 附录: 快速参考

### 常用导入路径速查

| 包 | 导入路径 |
|----|----------|
| 核心 Agent 接口 | `trpc.group/trpc-go/trpc-agent-go/agent` |
| LLM Agent | `trpc.group/trpc-go/trpc-agent-go/agent/llmagent` |
| Graph Agent | `trpc.group/trpc-go/trpc-agent-go/agent/graphagent` |
| Chain Agent | `trpc.group/trpc-go/trpc-agent-go/agent/chainagent` |
| Parallel Agent | `trpc.group/trpc-go/trpc-agent-go/agent/parallelagent` |
| Cycle Agent | `trpc.group/trpc-go/trpc-agent-go/agent/cycleagent` |
| Runner | `trpc.group/trpc-go/trpc-agent-go/runner` |
| 模型接口 | `trpc.group/trpc-go/trpc-agent-go/model` |
| OpenAI 模型 | `trpc.group/trpc-go/trpc-agent-go/model/openai` |
| 工具接口 | `trpc.group/trpc-go/trpc-agent-go/tool` |
| 函数工具 | `trpc.group/trpc-go/trpc-agent-go/tool/function` |
| 状态图引擎 | `trpc.group/trpc-go/trpc-agent-go/graph` |
| 记忆系统 | `trpc.group/trpc-go/trpc-agent-go/memory` |
| 会话管理 | `trpc.group/trpc-go/trpc-agent-go/session` |
| 事件系统 | `trpc.group/trpc-go/trpc-agent-go/event` |
| 知识库 | `trpc.group/trpc-go/trpc-agent-go/knowledge` |
| 技能系统 | `trpc.group/trpc-go/trpc-agent-go/skill` |
| 规划器 | `trpc.group/trpc-go/trpc-agent-go/planner` |
| 产物管理 | `trpc.group/trpc-go/trpc-agent-go/artifact` |
| 代码执行 | `trpc.group/trpc-go/trpc-agent-go/codeexecutor` |
| 提示模板 | `trpc.group/trpc-go/trpc-agent-go/prompt` |
| 插件系统 | `trpc.group/trpc-go/trpc-agent-go/plugin` |
| 评估 | `trpc.group/trpc-go/trpc-agent-go/evaluation` |
| 自进化 | `trpc.group/trpc-go/trpc-agent-go/evolution` |
| 可观测性 | `trpc.group/trpc-go/trpc-agent-go/telemetry` |
| Server | `trpc.group/trpc-go/trpc-agent-go/server` |

### 典型 Agent 创建流程

```
1. 创建 Model       → openai.New("model-name", opts...)
2. 创建 Tool        → function.NewFunctionTool(fn, opts...)
3. 创建 Agent       → llmagent.New("name", opts...)
4. 创建 Runner      → runner.NewRunner("app", agent, opts...)
5. 执行 Run         → runner.Run(ctx, userID, sessionID, message, opts...)
6. 处理事件流       → for event := range events { ... }
```

---

> 更多详细文档请访问: https://trpc-group.github.io/trpc-agent-go/  
> 示例代码: https://github.com/trpc-group/trpc-agent-go/tree/main/examples
