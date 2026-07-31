# Phoenix Java → Go 完整差异对照

> 生成时间: 2026-07-31 | 基于 Java 全模块提取 + Go 全文件审查

---

## 1. 缺失的 HTTP 端点

### 1.1 Platform Front 端点 (phoenix-platform 模块)

| Java 端点 | 方法 | Go 对应 | 状态 |
|:---|:---|:---|:---|
| `/api/front/stream/chat` | POST | 无 `/front/` 路由 | ❌ 缺失 |
| `/api/front/stream/chatsql` | GET | 无 `/front/` 路由 | ❌ 缺失 |
| `/api/front/harness/chat` | POST | 无 `/front/` 路由 | ❌ 缺失 |
| `/api/front/harness/confirm` | POST | 无 `/front/` 路由 | ❌ 缺失 |
| `/api/front/{agentId}/preset-questions` | GET | `/api/agent/:agentId/preset-question/` | ✅ 有 |
| `/api/front/addPresetQuestion` | POST | 无 | ❌ 缺失 |
| `/api/front/deletePresetQuestion/{id}` | DELETE | 无 | ❌ 缺失 |

### 1.2 Account/Group 查询端点

| Java 端点 | 方法 | Go 对应 | 状态 |
|:---|:---|:---|:---|
| `/platform/account-info/status/{status}` | GET | 无 | ❌ 缺失 |
| `/platform/account-tenant-info/account/{accountId}` | GET | 无 | ❌ 缺失 |
| `/platform/account-group-info/group/{groupId}` | GET | 无 | ❌ 缺失 |
| `/platform/account-group-info/account/{accountId}` | GET | 无 | ❌ 缺失 |

### 1.3 Privilege 端点

| Java 端点 | 方法 | Go 对应 | 状态 |
|:---|:---|:---|:---|
| `/api/privilege/user/username/{username}` | GET | 无 | ❌ 缺失 |
| `/api/privilege/module/{moduleId}/pvalues` | GET | 无 | ❌ 缺失 |
| `/api/privilege/module/{moduleId}/pvalue/{position}/{enabled}` | PUT | 无 | ❌ 缺失 |

### 1.4 Platform Sync 端点

| Java 端点 | 方法 | Go 对应 | 状态 |
|:---|:---|:---|:---|
| `/platform/sync/depts/{deptId}` | POST | 无 | ❌ 缺失 |
| `/platform/sync/depts/users/{deptId}` | POST | 无 | ❌ 缺失 |
| `/platform/sync/users/{userId}` | POST | 无 | ❌ 缺失 |

---

## 2. 缺失的运行时功能

| Java 功能 | 说明 | Go 状态 |
|:---|:---|:---|
| **MultiTurnContextManager** | 多轮对话上下文管理：beginTurn/finishTurn/buildContext/appendPlannerChunk | ❌ 缺失 |
| **Langfuse 链路追踪** | graph-stream/graph-feedback span 追踪，TRACE_THREAD_ID 传递 | ❌ 未接入 |
| **Graph Checkpoint/Saver** | SaverConfig 状态持久化，支持 interrupt/resume | ❌ 缺失 |
| **SessionTitleService** | 发消息后异步 LLM 标题生成 + SSE 推送 | ❌ 未触发 |
| **知识向量化事件** | ApplicationEvent → TransactionalEventListener → 异步向量化 | ❌ 缺失 |
| **PgVectorStore 向量存储** | 512维/余弦距离/HNSW 索引，schema/knowledge 检索 | ❌ 未接入 |
| **文本分割策略** | 5种分割器(token/recursive/sentence/semantic/paragraph) | ❌ 缺失 |
| **动态模型代理** | EmbeddingModel JDK 动态代理 + AiModelRegistry 热切换 | ❌ 缺失 |
| **HTML 报告下载** | ResponseEntity<byte[]> 二进制下载 | ❌ 返回 JSON 非二进制 |
| **StreamContext 管理** | ConcurrentHashMap 管理 sink/disposable/span/textType | ❌ 缺失 |

---

## 3. NL2SQL 工作流差异

### 3.1 状态键对比

| Java 状态键 | Go 对应 | 状态 |
|:---|:---|:---|
| `MULTI_TURN_CONTEXT` | 无 | ❌ 缺失 |
| `TRACE_THREAD_ID` | 无 | ❌ 缺失 |
| `SQL_GENERATE_SCHEMA_MISSING_ADVICE` | 无 | ❌ 缺失 |
| `PYTHON_FALLBACK_MODE` | 无 | ❌ 缺失 |
| `RESULT` | 无 | ❌ 缺失 |
| 其他 25+ 状态键 | 全部有 | ✅ |

### 3.2 边/路由差异

| Java 路由 | Go 对应 | 状态 |
|:---|:---|:---|
| `table_relation → table_relation` (重试) | 无重试自环 | ❌ 缺失 |
| `query_enhance → END` (空输出) | 无条件终止 | ❌ 缺失 |
| `schema_recall → END` (无表) | 无条件终止 | ❌ 缺失 |
| `interruptBefore(HUMAN_FEEDBACK_NODE)` | 无编译时中断配置 | ❌ 缺失 |
| `SaverConfig` 状态持久化 | 无 | ❌ 缺失 |

### 3.3 节点实现差异

| 节点 | 差异 |
|:---|:---|
| `EvidenceRecallNode` | Retriever 传入 nil，运行时不检索 |
| `SchemaRecallNode` | Retriever 传入 nil，运行时不检索 |
| `SqlExecuteNode` | 只支持 PostgreSQL，Java 支持 7 种数据库 |
| `PythonExecuteNode` | os/exec 无沙箱，Java 有 Docker 隔离 |
| `ReportGeneratorNode` | 无图表生成，Java 有 ECharts 配置 |

---

## 4. 配置差异

| Java 配置项 | Go 状态 |
|:---|:---|
| `vectorStore.tableTopkLimit=10` | ❌ 硬编码 |
| `vectorStore.similarityThreshold=0.2/0.4` | ❌ 缺失 |
| `vectorStore.batchDelTopkLimit=5000` | ❌ 缺失 |
| `maxSqlRetryCount=10` | ❌ 硬编码 3 |
| `maxturnhistory=5` | ❌ 缺失 |
| `enableSqlResultChart=true` | ❌ 缺失 |
| `textSplitter.chunkSize=1000` | ❌ 缺失 |
| `codePoolExecutor=DOCKER` | ❌ 只有 simulation |
| `limitMemory=500` / `cpuCore=1` | ❌ 缺失 |
| `pythonMaxTriesCount=5` | ❌ 硬编码 3 |

---

## 5. 缺失的实体字段

| Java 实体 | 字段 | Go 状态 |
|:---|:---|:---|
| AgentKnowledge | `isDeleted`, `isResourceCleaned` | ❌ `isDeleted` 用 DelFlag 替代, `isResourceCleaned` 缺失 |
| AgentDatasource | `datasource`(transient), `selectTables`(transient) | ❌ 无关联查询 |
| PrivilegeUser | `image`(byte[]), `aclTimestamp`, `pwdFtime`, `pwdInit` | ❌ 缺失 |
| PrivilegeModule | `component`, `image`, `isShow`, `categoryId` | 部分缺失 |
| RagFileInfo | `pdfType`, `pageTopMargin`, `pagesPerDocument`, `textSplitter` | ❌ 缺失 |

---

## 6. 已修复项

| 修复项 | 状态 |
|:---|:---|
| `/api/front/` 前端路由组 (chat/chatsql/harness/preset-questions) | ✅ 已添加 |
| `/platform/account-info/status/:status` 查询端点 | ✅ 已添加 |
| `/platform/sync/depts/:deptId` 子部门同步 | ✅ 已添加 |
| `/platform/sync/depts/users/:deptId` 部门用户同步 | ✅ 已添加 |
| `/platform/sync/users/:userId` 单用户同步 | ✅ 已添加 |
| `/api/privilege/user/username/:username` 查询端点 | ✅ 已添加 |
| 消息创建触发 SessionTitleService 异步标题生成 | ✅ 已添加 |
| Graph Handler 缓存 + 客户端断连处理 | ✅ 已修复 |
| ShortTermMemory StateDelta 消息存储 | ✅ 已修复 |
| RAG 工作流 StateGraph 实现 | ✅ 已修复 |
| web_search 工具 HTML 解析 | ✅ 已修复 |
| privilege check roleCode 检查 | ✅ 已修复 |
| NL2SQL Retriever 注入 | ✅ 已修复 |

## 7. 仍待解决 (需要外部依赖或较大工作量)

| 待解决项 | 说明 | 优先级 |
|:---|:---|:---|
| Langfuse 链路追踪接入 | 需要 OpenTelemetry 配置 | 🟡 中 |
| PgVectorStore 向量存储 | 需要向量库配置 + embedding 模型 | 🟡 中 |
| MultiTurnContextManager | 多轮对话上下文管理 | 🟡 中 |
| Python Docker 沙箱 | 需要 Docker SDK | 🟡 中 |
| 5 种文本分割策略 | 需要文本处理库 | 🟢 低 |
| Graph Checkpoint/Saver | 状态持久化 | 🟢 低 |
| HTML 报告二进制下载 | 需要模板文件 | 🟢 低 |
| 动态模型代理/热切换 | 需要 AIModelRegistry | 🟢 低 |
