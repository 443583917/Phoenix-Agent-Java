# Phoenix Go 重写 Phase 1 — 基础设施层实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 建立 godata/ 项目骨架，Gin 启动响应 `/echo`，所有基础设施包就绪

**Architecture:** 按 DDD 分层：`cmd/api` 入口 → `api/` Gin 路由 → `internal/` 业务核心 → `infra/` 基础设施。Phase 1 只搭建骨架和基础设施，不写业务逻辑

**Tech Stack:** Go 1.22+, Gin, GORM, Viper, Zap, go-redis, bigcache, Casbin, RabbitMQ AMQP, Milvus SDK, OpenTelemetry, Sonyflake

## Global Constraints

- 所有 import 路径以 `github.com/phoenix-agent-go` 为 module root
- 配置统一使用 Viper，YAML 文件按服务分目录
- 统一响应格式 `{code, message, data, traceId}`
- 统一错误码 `ErrCode{Code int, Msg string}`
- 全链路 ctx 传递
- 每个 infra 包零外部依赖或仅依赖标准库+知名三方库
- Java 版本没有对应映射的模块先空目录，后续填充

## 回测策略

| 验证项 | 方法 | 预期 |
|:---|:---|:---|
| 项目编译 | `go build ./...` | 零错误 |
| 服务启动 | `go run ./cmd/api` | 监听 :8066 |
| 健康检查 | `curl http://localhost:8066/echo` | `{"code":0,"message":"success","data":"ok"}` |
| 配置加载 | 启动日志 | 输出 db/redis 连接配置 |
| 中间件链 | 请求带 Header | 日志含 traceId + 耗时 |
| 静态检查 | `go vet ./...` | 零警告 |

---

### Task 1: 初始化项目骨架

**Files:**
- Modify: `godata/go.mod`
- Create: 全部空目录（占位 .gitkeep）

**Interfaces:**
- Produces: `module github.com/phoenix-agent-go`，Go 1.22

- [ ] **Step 1: 更新 go.mod**

```go
module github.com/phoenix-agent-go

go 1.22.0

require (
    github.com/allegro/bigcache/v3 v3.1.0
    github.com/casbin/casbin/v2 v2.103.0
    github.com/gin-gonic/gin v1.10.0
    github.com/golang-jwt/jwt/v5 v5.2.1
    github.com/rabbitmq/amqp091-go v1.10.0
    github.com/redis/go-redis/v9 v9.7.0
    github.com/robfig/cron/v3 v3.0.1
    github.com/sony/sonyflake v1.2.0
    github.com/spf13/viper v1.19.0
    github.com/stretchr/testify v1.10.0
    go.opentelemetry.io/otel v1.32.0
    go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.32.0
    go.opentelemetry.io/otel/sdk v1.32.0
    go.opentelemetry.io/otel/trace v1.32.0
    go.uber.org/zap v1.27.0
    gorm.io/driver/postgres v1.5.11
    gorm.io/gorm v1.25.12
    trpc.group/trpc-go/trpc-agent-go v1.10.0
)
```

- [ ] **Step 2: 创建目录树并放置 .gitkeep 占位**

```bash
cd d:/GitHub/Phoenix-Agent-Java/godata

# 移除旧的 main.go（保留参考）
mkdir -p _archive
[ -f main.go ] && mv main.go _archive/
[ -f trpc-agent-go-使用说明.md ] && mv trpc-agent-go-使用说明.md _archive/

# 创建全部目录
mkdir -p cmd/api cmd/rpc cmd/agent cmd/job
mkdir -p api/handler/{agent,datasource,chat,knowledge,modelconfig,prompt,semanticmodel,privilege,platform,rag,kg,common}
mkdir -p api/middleware
mkdir -p rpc/proto rpc/server rpc/client
mkdir -p agent/runtime agent/agents
mkdir -p agent/tools/{function,rpc,datasource,privilege,web,external}
mkdir -p agent/workflows/{nl2sql/nodes,rag/nodes,kg/nodes}
mkdir -p agent/memory agent/knowledge agent/runner
mkdir -p internal/domain/{agent,datasource,chat,knowledge,modelconfig,prompt,semanticmodel,privilege,platform,rag,common}
mkdir -p internal/usecase internal/service internal/repository
mkdir -p internal/dao/{db,cache,queue,external}
mkdir -p internal/event/handler/{agent,chat,privilege}
mkdir -p internal/job/jobs
mkdir -p internal/model internal/config
mkdir -p infra/{logger,config,monitoring,queue,cache,lock,id,utils}
mkdir -p configs/{api,rpc,agent,job}
mkdir -p migrations scripts docs storage

# 空目录放 .gitkeep（不含 Go 源码的目录）
for d in cmd/rpc cmd/agent cmd/job \
    api/handler/{agent,datasource,chat,knowledge,modelconfig,prompt,semanticmodel,privilege,platform,rag,kg,common} \
    rpc/proto rpc/server rpc/client \
    agent/runtime agent/agents \
    agent/tools/{function,rpc,datasource,privilege,web,external} \
    agent/workflows/{nl2sql/nodes,rag/nodes,kg/nodes} \
    agent/memory agent/knowledge agent/runner \
    internal/domain/{agent,datasource,chat,knowledge,modelconfig,prompt,semanticmodel,privilege,platform,rag,common} \
    internal/usecase internal/service internal/repository \
    internal/dao/{db,cache,queue,external} \
    internal/event/handler/{agent,chat,privilege} \
    internal/job/jobs internal/model \
    configs/{api,rpc,agent,job} \
    migrations scripts docs storage; do
    touch "$d/.gitkeep"
done
```

- [ ] **Step 3: 下载依赖并验证**

```bash
cd d:/GitHub/Phoenix-Agent-Java/godata
go mod tidy
go mod download
```

Expected: 无错误输出，`go.sum` 生成。

- [ ] **Step 4: 编译空项目验证**

```bash
go build ./...
```

Expected: 无输出（无 .go 源文件时不报错）。

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat(phase1): initialize project skeleton with go.mod and directory tree"
```

---

### Task 2: 配置文件（YAML）

**Files:**
- Create: `configs/db.yaml`
- Create: `configs/redis.yaml`
- Create: `configs/milvus.yaml`
- Create: `configs/rabbitmq.yaml`
- Create: `configs/monitor.yaml`
- Create: `configs/api/app.yaml`
- Create: `configs/agent/app.yaml`
- Create: `configs/rpc/app.yaml`
- Create: `configs/job/app.yaml`

**Interfaces:**
- Produces: 所有 YAML 配置文件，可被 Viper 读取

- [ ] **Step 1: 创建共享配置文件**

`configs/db.yaml`:
```yaml
database:
  host: "127.0.0.1"
  port: 5432
  user: "phoenix"
  password: "phoenix"
  name: "phoenix"
  sslmode: "disable"
  max_open_conns: 25
  max_idle_conns: 10
  conn_max_lifetime: "5m"
```

`configs/redis.yaml`:
```yaml
redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0
  pool_size: 10
  min_idle_conns: 5
  dial_timeout: "5s"
  read_timeout: "3s"
  write_timeout: "3s"
```

`configs/milvus.yaml`:
```yaml
milvus:
  addr: "127.0.0.1:19530"
  collection: "phoenix_vectors"
  dim: 512
  index_type: "HNSW"
  metric_type: "COSINE"
  timeout: "30s"
```

`configs/rabbitmq.yaml`:
```yaml
rabbitmq:
  addr: "amqp://guest:guest@127.0.0.1:5672/"
  exchange: "phoenix.events"
  prefetch_count: 10
  reconnect_delay: "3s"
  max_reconnect_attempts: 5
```

`configs/monitor.yaml`:
```yaml
monitor:
  otel_endpoint: "127.0.0.1:4317"
  service_name: "phoenix-go"
  service_version: "2.0.0"
  trace_sample_rate: 1.0
  log_level: "info"
  log_format: "json"
```

- [ ] **Step 2: 创建服务专用配置文件**

`configs/api/app.yaml`:
```yaml
server:
  port: 8066
  mode: "debug"           # debug | release | test
  read_timeout: "30s"
  write_timeout: "60s"
  max_header_bytes: 1048576
cors:
  allow_origins:
    - "*"
  allow_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allow_headers:
    - "Origin"
    - "Content-Type"
    - "Authorization"
    - "X-Trace-Id"
  expose_headers:
    - "Content-Length"
  allow_credentials: true
  max_age: "12h"
```

`configs/agent/app.yaml`:
```yaml
agent:
  model:
    provider: "deepseek"
    model: "deepseek-chat"
    api_key: "${AI_API_KEY}"
    base_url: "https://api.deepseek.com"
  stream: true
  max_tokens: 4096
  temperature: 0.7
  session_timeout: "30m"
```

`configs/rpc/app.yaml`:
```yaml
rpc:
  port: 9090
  reflection: true
```

`configs/job/app.yaml`:
```yaml
job:
  cron_specs:
    daily_report: "0 8 * * *"
    audit_sync: "*/30 * * * *"
```

- [ ] **Step 3: Commit**

```bash
git add configs/
git commit -m "feat(phase1): add all YAML configuration files"
```

---

### Task 3: 配置加载层（Viper）

**Files:**
- Create: `internal/config/app.go`
- Create: `internal/config/db.go`
- Create: `internal/config/redis.go`
- Create: `internal/config/milvus.go`
- Create: `internal/config/rabbitmq.go`
- Create: `internal/config/agent.go`
- Create: `internal/config/rpc.go`
- Create: `internal/config/monitor.go`

**Interfaces:**
- Produces: `type DBConfig struct`, `type RedisConfig struct`, `type MilvusConfig struct`, `type RabbitMQConfig struct`, `type MonitorConfig struct`, `type ServerConfig struct`, `type AgentConfig struct`, `type RPCConfig struct`, `type JobConfig struct`, `type CorsConfig struct`

- [ ] **Step 1: 定义配置结构体**

`internal/config/app.go`:
```go
package config

import "time"

type ServerConfig struct {
    Port           int           `mapstructure:"port"`
    Mode           string        `mapstructure:"mode"`
    ReadTimeout    time.Duration `mapstructure:"read_timeout"`
    WriteTimeout   time.Duration `mapstructure:"write_timeout"`
    MaxHeaderBytes int           `mapstructure:"max_header_bytes"`
}

type CorsConfig struct {
    AllowOrigins     []string      `mapstructure:"allow_origins"`
    AllowMethods     []string      `mapstructure:"allow_methods"`
    AllowHeaders     []string      `mapstructure:"allow_headers"`
    ExposeHeaders    []string      `mapstructure:"expose_headers"`
    AllowCredentials bool          `mapstructure:"allow_credentials"`
    MaxAge           time.Duration `mapstructure:"max_age"`
}
```

`internal/config/db.go`:
```go
package config

import "time"

type DBConfig struct {
    Host            string        `mapstructure:"host"`
    Port            int           `mapstructure:"port"`
    User            string        `mapstructure:"user"`
    Password        string        `mapstructure:"password"`
    Name            string        `mapstructure:"name"`
    SSLMode         string        `mapstructure:"sslmode"`
    MaxOpenConns    int           `mapstructure:"max_open_conns"`
    MaxIdleConns    int           `mapstructure:"max_idle_conns"`
    ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}
```

`internal/config/redis.go`:
```go
package config

import "time"

type RedisConfig struct {
    Addr         string        `mapstructure:"addr"`
    Password     string        `mapstructure:"password"`
    DB           int           `mapstructure:"db"`
    PoolSize     int           `mapstructure:"pool_size"`
    MinIdleConns int           `mapstructure:"min_idle_conns"`
    DialTimeout  time.Duration `mapstructure:"dial_timeout"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
}
```

`internal/config/milvus.go`:
```go
package config

import "time"

type MilvusConfig struct {
    Addr       string        `mapstructure:"addr"`
    Collection string        `mapstructure:"collection"`
    Dim        int           `mapstructure:"dim"`
    IndexType  string        `mapstructure:"index_type"`
    MetricType string        `mapstructure:"metric_type"`
    Timeout    time.Duration `mapstructure:"timeout"`
}
```

`internal/config/rabbitmq.go`:
```go
package config

import "time"

type RabbitMQConfig struct {
    Addr                 string        `mapstructure:"addr"`
    Exchange             string        `mapstructure:"exchange"`
    PrefetchCount        int           `mapstructure:"prefetch_count"`
    ReconnectDelay       time.Duration `mapstructure:"reconnect_delay"`
    MaxReconnectAttempts int           `mapstructure:"max_reconnect_attempts"`
}
```

`internal/config/agent.go`:
```go
package config

import "time"

type AgentModelConfig struct {
    Provider string  `mapstructure:"provider"`
    Model    string  `mapstructure:"model"`
    APIKey   string  `mapstructure:"api_key"`
    BaseURL  string  `mapstructure:"base_url"`
}

type AgentConfig struct {
    Model          AgentModelConfig `mapstructure:"model"`
    Stream         bool             `mapstructure:"stream"`
    MaxTokens      int              `mapstructure:"max_tokens"`
    Temperature    float64          `mapstructure:"temperature"`
    SessionTimeout time.Duration    `mapstructure:"session_timeout"`
}
```

`internal/config/rpc.go`:
```go
package config

type RPCConfig struct {
    Port       int  `mapstructure:"port"`
    Reflection bool `mapstructure:"reflection"`
}
```

`internal/config/monitor.go`:
```go
package config

type MonitorConfig struct {
    OTELEndpoint    string  `mapstructure:"otel_endpoint"`
    ServiceName     string  `mapstructure:"service_name"`
    ServiceVersion  string  `mapstructure:"service_version"`
    TraceSampleRate float64 `mapstructure:"trace_sample_rate"`
    LogLevel        string  `mapstructure:"log_level"`
    LogFormat       string  `mapstructure:"log_format"`
}
```

- [ ] **Step 2: 编译验证**

```bash
cd d:/GitHub/Phoenix-Agent-Java/godata
go build ./internal/config/...
```

Expected: 无错误。

- [ ] **Step 3: Commit**

```bash
git add internal/config/
git commit -m "feat(phase1): add Viper config struct definitions"
```

---

### Task 4: 日志包（infra/logger）

**Files:**
- Create: `infra/logger/logger.go`
- Test: `infra/logger/logger_test.go`

**Interfaces:**
- Produces: `func Init(cfg *config.MonitorConfig) error`, `func L() *zap.Logger`, `func Sync()`

- [ ] **Step 1: 编写测试**

`infra/logger/logger_test.go`:
```go
package logger

import (
    "testing"

    "github.com/phoenix-agent-go/internal/config"
    "github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
    cfg := &config.MonitorConfig{
        LogLevel:  "debug",
        LogFormat: "console",
    }
    err := Init(cfg)
    assert.NoError(t, err)
    assert.NotNil(t, L())

    L().Info("test log message",
        zap.String("key", "value"),
    )
    Sync()
}
```

- [ ] **Step 2: 运行测试验证失败**

```bash
go test ./infra/logger/ -v
```

Expected: FAIL — `undefined: Init`

- [ ] **Step 3: 实现 logger**

`infra/logger/logger.go`:
```go
package logger

import (
    "sync"

    "github.com/phoenix-agent-go/internal/config"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var (
    logger *zap.Logger
    once   sync.Once
)

func Init(cfg *config.MonitorConfig) error {
    var err error
    once.Do(func() {
        level := zapcore.InfoLevel
        if cfg != nil {
            _ = level.UnmarshalText([]byte(cfg.LogLevel))
        }

        var zapCfg zap.Config
        if cfg != nil && cfg.LogFormat == "json" {
            zapCfg = zap.NewProductionConfig()
        } else {
            zapCfg = zap.NewDevelopmentConfig()
        }
        zapCfg.Level = zap.NewAtomicLevelAt(level)
        zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

        logger, err = zapCfg.Build(zap.AddCallerSkip(1))
    })
    return err
}

func L() *zap.Logger {
    if logger == nil {
        l, _ := zap.NewDevelopment()
        return l
    }
    return logger
}

func Sync() {
    if logger != nil {
        _ = logger.Sync()
    }
}
```

- [ ] **Step 4: 修复测试中的 import**

在 `logger_test.go` 顶部添加:
```go
import (
    "testing"

    "github.com/phoenix-agent-go/internal/config"
    "github.com/stretchr/testify/assert"
    "go.uber.org/zap"
)
```

- [ ] **Step 5: 运行测试**

```bash
go test ./infra/logger/ -v
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add infra/logger/
git commit -m "feat(phase1): add Zap logger wrapper"
```

---

### Task 5: 响应格式 + 错误码（infra/response + infra/errcode）

**Files:**
- Create: `infra/response/response.go`
- Create: `infra/errcode/errcode.go`

**Interfaces:**
- Produces: `type Response struct`, `type PageResponse struct`, `func Success(c *gin.Context, data interface{})`, `func SuccessPage(c *gin.Context, data interface{}, total int64, page, size int)`, `func Error(c *gin.Context, ec ErrCode)`, `func ErrorWithMsg(c *gin.Context, ec ErrCode, msg string)`
- Produces: `type ErrCode struct{Code int; Msg string}` + 预定义错误码

- [ ] **Step 1: 实现错误码**

`infra/errcode/errcode.go`:
```go
package errcode

type ErrCode struct {
    Code int
    Msg  string
}

var (
    Success       = ErrCode{0, "success"}
    Unauthorized  = ErrCode{401, "未认证"}
    Forbidden     = ErrCode{403, "无权限"}
    NotFound      = ErrCode{404, "资源不存在"}
    InternalError = ErrCode{500, "服务器内部错误"}

    InvalidParams   = ErrCode{1001, "参数校验失败"}
    TooManyRequests = ErrCode{1002, "请求过于频繁"}

    AgentOffline     = ErrCode{2001, "智能体已下线"}
    AgentNotFound    = ErrCode{2002, "智能体不存在"}
    SessionExpired   = ErrCode{2003, "会话已过期"}

    DatasourceError   = ErrCode{3001, "数据源连接失败"}
    DatasourceNotFound = ErrCode{3002, "数据源不存在"}
    SQLError           = ErrCode{3003, "SQL执行失败"}

    ModelError     = ErrCode{4001, "模型调用失败"}
    EmbeddingError = ErrCode{4002, "向量化失败"}

    MilvusError  = ErrCode{5001, "向量检索失败"}
    QueueError   = ErrCode{5002, "消息队列异常"}
    CacheError   = ErrCode{5003, "缓存服务异常"}
)
```

- [ ] **Step 2: 实现响应工具**

`infra/response/response.go`:
```go
package response

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/phoenix-agent-go/infra/errcode"
)

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

func getTraceID(c *gin.Context) string {
    if tid := c.GetString("trace_id"); tid != "" {
        return tid
    }
    return ""
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    errcode.Success.Code,
        Message: errcode.Success.Msg,
        Data:    data,
        TraceID: getTraceID(c),
    })
}

func SuccessPage(c *gin.Context, data interface{}, total int64, page, size int) {
    c.JSON(http.StatusOK, PageResponse{
        Code:    errcode.Success.Code,
        Message: errcode.Success.Msg,
        Data:    data,
        Total:   total,
        Page:    page,
        Size:    size,
        TraceID: getTraceID(c),
    })
}

func Error(c *gin.Context, ec errcode.ErrCode) {
    c.JSON(http.StatusOK, Response{
        Code:    ec.Code,
        Message: ec.Msg,
        TraceID: getTraceID(c),
    })
}

func ErrorWithMsg(c *gin.Context, ec errcode.ErrCode, msg string) {
    c.JSON(http.StatusOK, Response{
        Code:    ec.Code,
        Message: msg,
        TraceID: getTraceID(c),
    })
}

func ErrorWithStatus(c *gin.Context, httpStatus int, ec errcode.ErrCode) {
    c.JSON(httpStatus, Response{
        Code:    ec.Code,
        Message: ec.Msg,
        TraceID: getTraceID(c),
    })
}
```

- [ ] **Step 3: 编译验证**

```bash
go build ./infra/response/ ./infra/errcode/
```

Expected: 无错误。

- [ ] **Step 4: Commit**

```bash
git add infra/response/ infra/errcode/
git commit -m "feat(phase1): add unified response and error code packages"
```

---

### Task 6: 工具包（infra/validate, infra/pagination, infra/jwt, infra/sse, infra/id）

**Files:**
- Create: `infra/validate/validate.go`
- Create: `infra/pagination/pagination.go`
- Create: `infra/jwt/jwt.go`
- Create: `infra/sse/sse.go`
- Create: `infra/id/id.go`

**Interfaces:**
- Produces: `func ValidateStruct(v interface{}) error`
- Produces: `type PageQuery struct`, `func ParsePageQuery(c *gin.Context) PageQuery`
- Produces: `type JWTManager struct`, `func NewJWTManager(secret string, expire time.Duration) *JWTManager`, `GenerateToken`, `ParseToken`
- Produces: `type Event struct`, `func Stream(c *gin.Context, events <-chan Event)`
- Produces: `func GenerateID() uint64`

- [ ] **Step 1: 实现参数校验**

`infra/validate/validate.go`:
```go
package validate

import (
    "github.com/go-playground/validator/v10"
)

var v = validator.New()

func ValidateStruct(s interface{}) error {
    return v.Struct(s)
}
```

Note: 需在 go.mod 添加 `github.com/go-playground/validator/v10 v10.23.0`，之后执行 `go mod tidy`。

- [ ] **Step 2: 实现分页工具**

`infra/pagination/pagination.go`:
```go
package pagination

import (
    "strconv"

    "github.com/gin-gonic/gin"
)

type PageQuery struct {
    Page int `json:"page"`
    Size int `json:"size"`
}

const (
    DefaultPage = 1
    DefaultSize = 10
    MaxSize     = 100
)

func ParsePageQuery(c *gin.Context) PageQuery {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

    if page < 1 {
        page = DefaultPage
    }
    if size < 1 {
        size = DefaultSize
    }
    if size > MaxSize {
        size = MaxSize
    }

    return PageQuery{Page: page, Size: size}
}

func (p PageQuery) Offset() int {
    return (p.Page - 1) * p.Size
}
```

- [ ] **Step 3: 实现 JWT 工具**

`infra/jwt/jwt.go`:
```go
package jwt

import (
    "time"

    jwtlib "github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
    secret []byte
    expire time.Duration
}

type Claims struct {
    UserID   uint64 `json:"userId"`
    Username string `json:"username"`
    jwtlib.RegisteredClaims
}

func NewJWTManager(secret string, expire time.Duration) *JWTManager {
    return &JWTManager{
        secret: []byte(secret),
        expire: expire,
    }
}

func (m *JWTManager) GenerateToken(userID uint64, username string) (string, error) {
    now := time.Now()
    claims := Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwtlib.RegisteredClaims{
            IssuedAt:  jwtlib.NewNumericDate(now),
            ExpiresAt: jwtlib.NewNumericDate(now.Add(m.expire)),
        },
    }
    token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
    return token.SignedString(m.secret)
}

func (m *JWTManager) ParseToken(tokenStr string) (*Claims, error) {
    token, err := jwtlib.ParseWithClaims(tokenStr, &Claims{},
        func(t *jwtlib.Token) (interface{}, error) {
            return m.secret, nil
        })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    return nil, jwtlib.ErrSignatureInvalid
}
```

- [ ] **Step 4: 实现 SSE 工具**

`infra/sse/sse.go`:
```go
package sse

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
)

type Event struct {
    Type string      `json:"type"`
    Data interface{} `json:"data"`
}

func Stream(c *gin.Context, events <-chan Event) {
    c.Writer.Header().Set("Content-Type", "text/event-stream")
    c.Writer.Header().Set("Cache-Control", "no-cache")
    c.Writer.Header().Set("Connection", "keep-alive")
    c.Writer.Header().Set("X-Accel-Buffering", "no")
    c.Writer.WriteHeader(http.StatusOK)

    flusher, ok := c.Writer.(http.Flusher)
    if !ok {
        return
    }

    for {
        select {
        case <-c.Request.Context().Done():
            return
        case event, ok := <-events:
            if !ok {
                return
            }
            data, err := json.Marshal(event)
            if err != nil {
                continue
            }
            fmt.Fprintf(c.Writer, "data: %s\n\n", data)
            flusher.Flush()
        }
    }
}
```

- [ ] **Step 5: 实现 ID 生成器**

`infra/id/id.go`:
```go
package id

import (
    "sync"

    "github.com/sony/sonyflake"
)

var (
    sf   *sonyflake.Sonyflake
    once sync.Once
)

func initSonyflake() {
    once.Do(func() {
        sf = sonyflake.NewSonyflake(sonyflake.Settings{})
    })
}

func GenerateID() (uint64, error) {
    initSonyflake()
    return sf.NextID()
}

func MustGenerateID() uint64 {
    id, err := GenerateID()
    if err != nil {
        panic("failed to generate ID: " + err.Error())
    }
    return id
}
```

- [ ] **Step 6: 编译验证**

```bash
go mod tidy
go build ./infra/validate/ ./infra/pagination/ ./infra/jwt/ ./infra/sse/ ./infra/id/
```

Expected: 无错误。

- [ ] **Step 7: Commit**

```bash
git add infra/validate/ infra/pagination/ infra/jwt/ infra/sse/ infra/id/ go.mod go.sum
git commit -m "feat(phase1): add validate, pagination, jwt, sse, and id utility packages"
```

---

### Task 7: 基础设施连接管理（infra/cache, infra/queue, infra/lock, infra/monitoring）

**Files:**
- Create: `infra/cache/cache.go`
- Create: `infra/queue/queue.go`
- Create: `infra/lock/lock.go`
- Create: `infra/monitoring/monitoring.go`

**Interfaces:**
- Produces: `func InitRedis(cfg *config.RedisConfig) (*redis.Client, error)`, `func InitBigCache() (*bigcache.BigCache, error)`
- Produces: `func InitRabbitMQ(cfg *config.RabbitMQConfig) (*amqp.Connection, *amqp.Channel, error)`
- Produces: `type Locker interface { Lock/Unlock }`, `func NewRedisLocker(client *redis.Client) Locker`
- Produces: `func InitTracer(cfg *config.MonitorConfig) (*trace.TracerProvider, error)`

- [ ] **Step 1: 实现缓存初始化**

`infra/cache/cache.go`:
```go
package cache

import (
    "context"
    "time"

    "github.com/allegro/bigcache/v3"
    "github.com/phoenix-agent-go/internal/config"
    "github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.RedisConfig) (*redis.Client, error) {
    rdb := redis.NewClient(&redis.Options{
        Addr:         cfg.Addr,
        Password:     cfg.Password,
        DB:           cfg.DB,
        PoolSize:     cfg.PoolSize,
        MinIdleConns: cfg.MinIdleConns,
        DialTimeout:  cfg.DialTimeout,
        ReadTimeout:  cfg.ReadTimeout,
        WriteTimeout: cfg.WriteTimeout,
    })

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := rdb.Ping(ctx).Err(); err != nil {
        return nil, err
    }
    return rdb, nil
}

func InitBigCache() (*bigcache.BigCache, error) {
    return bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
}
```

- [ ] **Step 2: 实现消息队列初始化**

`infra/queue/queue.go`:
```go
package queue

import (
    "github.com/phoenix-agent-go/internal/config"
    amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ(cfg *config.RabbitMQConfig) (*amqp.Connection, *amqp.Channel, error) {
    conn, err := amqp.Dial(cfg.Addr)
    if err != nil {
        return nil, nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, nil, err
    }

    if err := ch.Qos(cfg.PrefetchCount, 0, false); err != nil {
        ch.Close()
        conn.Close()
        return nil, nil, err
    }

    err = ch.ExchangeDeclare(
        cfg.Exchange,
        "topic",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        ch.Close()
        conn.Close()
        return nil, nil, err
    }

    return conn, ch, nil
}
```

- [ ] **Step 3: 实现分布式锁**

`infra/lock/lock.go`:
```go
package lock

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
)

type Locker interface {
    Lock(ctx context.Context, key string, ttl time.Duration) (bool, error)
    Unlock(ctx context.Context, key string) error
}

type redisLocker struct {
    client *redis.Client
}

func NewRedisLocker(client *redis.Client) Locker {
    return &redisLocker{client: client}
}

func (l *redisLocker) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
    return l.client.SetNX(ctx, "lock:"+key, "1", ttl).Result()
}

func (l *redisLocker) Unlock(ctx context.Context, key string) error {
    return l.client.Del(ctx, "lock:"+key).Err()
}
```

- [ ] **Step 4: 实现可观测性初始化**

`infra/monitoring/monitoring.go`:
```go
package monitoring

import (
    "context"

    "github.com/phoenix-agent-go/internal/config"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func InitTracer(ctx context.Context, cfg *config.MonitorConfig) (*sdktrace.TracerProvider, error) {
    exporter, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint(cfg.OTELEndpoint),
        otlptracegrpc.WithInsecure(),
    )
    if err != nil {
        return nil, err
    }

    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceName(cfg.ServiceName),
            semconv.ServiceVersion(cfg.ServiceVersion),
        ),
    )
    if err != nil {
        return nil, err
    }

    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
        sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.TraceSampleRate)),
    )

    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ))

    return tp, nil
}
```

- [ ] **Step 5: 编译验证**

```bash
go mod tidy
go build ./infra/cache/ ./infra/queue/ ./infra/lock/ ./infra/monitoring/
```

Expected: 无错误。

- [ ] **Step 6: Commit**

```bash
git add infra/cache/ infra/queue/ infra/lock/ infra/monitoring/ go.mod go.sum
git commit -m "feat(phase1): add cache, queue, lock, and monitoring infrastructure packages"
```

---

### Task 8: Viper 配置加载器（infra/config）

**Files:**
- Create: `infra/config/config.go`
- Create: `infra/config/config_test.go`

**Interfaces:**
- Produces: `type AppConfig struct` (聚合所有子配置), `func Load(serviceName string) (*AppConfig, error)`

- [ ] **Step 1: 编写测试**

`infra/config/config_test.go`:
```go
package config

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
    cfg, err := Load("api")
    assert.NoError(t, err)
    assert.NotNil(t, cfg)
    assert.Equal(t, 8066, cfg.Server.Port)
    assert.Equal(t, "127.0.0.1", cfg.DB.Host)
}
```

- [ ] **Step 2: 实现配置加载器**

`infra/config/config.go`:
```go
package config

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    appcfg "github.com/phoenix-agent-go/internal/config"
    "github.com/spf13/viper"
)

type AppConfig struct {
    Server   appcfg.ServerConfig
    DB       appcfg.DBConfig
    Redis    appcfg.RedisConfig
    Milvus   appcfg.MilvusConfig
    RabbitMQ appcfg.RabbitMQConfig
    Monitor  appcfg.MonitorConfig
    Agent    appcfg.AgentConfig
    RPC      appcfg.RPCConfig
    Cors     appcfg.CorsConfig
}

func Load(serviceName string) (*AppConfig, error) {
    v := viper.New()

    // 配置搜索路径
    v.SetConfigName("db")
    v.SetConfigType("yaml")
    v.AddConfigPath("./configs")
    v.AddConfigPath("../configs")
    v.AddConfigPath("../../configs")

    // 环境变量覆盖
    v.AutomaticEnv()
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    v.SetEnvPrefix("PHOENIX")

    if err := loadConfigFile(v, "db"); err != nil {
        return nil, err
    }
    if err := loadConfigFile(v, "redis"); err != nil {
        return nil, err
    }
    if err := loadConfigFile(v, "milvus"); err != nil {
        return nil, err
    }
    if err := loadConfigFile(v, "rabbitmq"); err != nil {
        return nil, err
    }
    if err := loadConfigFile(v, "monitor"); err != nil {
        return nil, err
    }

    // 服务专用配置
    serviceConfigPath := filepath.Join("configs", serviceName, "app.yaml")
    if _, err := os.Stat(serviceConfigPath); err == nil {
        v.SetConfigFile(serviceConfigPath)
        if err := v.MergeInConfig(); err != nil {
            return nil, fmt.Errorf("failed to load service config %s: %w", serviceConfigPath, err)
        }
    }

    // 展开 ${VAR} 环境变量
    for _, key := range v.AllKeys() {
        val := v.GetString(key)
        if strings.HasPrefix(val, "${") && strings.HasSuffix(val, "}") {
            envVar := val[2 : len(val)-1]
            v.Set(key, os.Getenv(envVar))
        }
    }

    cfg := &AppConfig{}
    if err := v.UnmarshalKey("database", &cfg.DB); err != nil {
        return nil, err
    }
    if err := v.UnmarshalKey("redis", &cfg.Redis); err != nil {
        return nil, err
    }
    if err := v.UnmarshalKey("milvus", &cfg.Milvus); err != nil {
        return nil, err
    }
    if err := v.UnmarshalKey("rabbitmq", &cfg.RabbitMQ); err != nil {
        return nil, err
    }
    if err := v.UnmarshalKey("monitor", &cfg.Monitor); err != nil {
        return nil, err
    }
    if err := v.UnmarshalKey("server", &cfg.Server); err != nil {
        return nil, err
    }
    if err := v.UnmarshalKey("cors", &cfg.Cors); err != nil {
        return nil, err
    }
    if err := v.UnmarshalKey("agent", &cfg.Agent); err != nil {
        return nil, err
    }
    if err := v.UnmarshalKey("rpc", &cfg.RPC); err != nil {
        return nil, err
    }

    return cfg, nil
}

func loadConfigFile(v *viper.Viper, name string) error {
    v.SetConfigName(name)
    if err := v.MergeInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return fmt.Errorf("failed to load %s.yaml: %w", name, err)
        }
        // 配置文件不存在可接受
    }
    return nil
}
```

- [ ] **Step 3: 运行测试**

```bash
cd d:/GitHub/Phoenix-Agent-Java/godata
go test ./infra/config/ -v
```

Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add infra/config/
git commit -m "feat(phase1): add Viper configuration loader"
```

---

### Task 9: API 中间件

**Files:**
- Create: `api/middleware/recovery.go`
- Create: `api/middleware/logger.go`
- Create: `api/middleware/tracing.go`
- Create: `api/middleware/cors.go`
- Create: `api/middleware/auth.go` (骨架)
- Create: `api/middleware/rbac.go` (骨架)
- Create: `api/middleware/ratelimit.go` (骨架)

**Interfaces:**
- Produces: `func Recovery() gin.HandlerFunc`, `func Logger() gin.HandlerFunc`, `func Tracing() gin.HandlerFunc`, `func CORS(cfg *config.CorsConfig) gin.HandlerFunc`, `func Auth() gin.HandlerFunc`, `func RBAC() gin.HandlerFunc`, `func RateLimit() gin.HandlerFunc`

- [ ] **Step 1: 实现 Recovery 中间件**

`api/middleware/recovery.go`:
```go
package middleware

import (
    "net"
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/phoenix-agent-go/infra/errcode"
    "github.com/phoenix-agent-go/infra/response"
    "go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                var brokenPipe bool
                if ne, ok := err.(*net.OpError); ok {
                    if se, ok := ne.Err.(*os.SyscallError); ok {
                        if strings.Contains(strings.ToLower(se.Error()), "broken pipe") {
                            brokenPipe = true
                        }
                    }
                }

                if brokenPipe {
                    c.Abort()
                    return
                }

                zap.L().Error("panic recovered",
                    zap.Any("error", err),
                    zap.String("path", c.Request.URL.Path),
                )
                response.ErrorWithStatus(c, http.StatusInternalServerError, errcode.InternalError)
                c.Abort()
            }
        }()
        c.Next()
    }
}
```

- [ ] **Step 2: 实现 Logger 中间件**

`api/middleware/logger.go`:
```go
package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery

        c.Next()

        latency := time.Since(start)
        status := c.Writer.Status()
        clientIP := c.ClientIP()
        method := c.Request.Method

        zap.L().Info("request",
            zap.Int("status", status),
            zap.String("method", method),
            zap.String("path", path),
            zap.String("query", query),
            zap.String("ip", clientIP),
            zap.Duration("latency", latency),
            zap.String("trace_id", c.GetString("trace_id")),
        )
    }
}
```

- [ ] **Step 3: 实现 Tracing 中间件**

`api/middleware/tracing.go`:
```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func Tracing() gin.HandlerFunc {
    return func(c *gin.Context) {
        traceID := c.GetHeader("X-Trace-Id")
        if traceID == "" {
            traceID = uuid.New().String()
        }
        c.Set("trace_id", traceID)
        c.Header("X-Trace-Id", traceID)

        tracer := otel.Tracer("phoenix-api")
        ctx, span := tracer.Start(c.Request.Context(),
            c.Request.Method+" "+c.Request.URL.Path,
            trace.WithSpanKind(trace.SpanKindServer),
        )
        defer span.End()

        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}
```

- [ ] **Step 4: 实现 CORS 中间件**

`api/middleware/cors.go`:
```go
package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/phoenix-agent-go/internal/config"
)

func CORS(cfg *config.CorsConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        allowOrigin := ""

        for _, o := range cfg.AllowOrigins {
            if o == "*" || o == origin {
                allowOrigin = origin
                if o == "*" {
                    allowOrigin = "*"
                }
                break
            }
        }

        if allowOrigin == "" {
            c.Next()
            return
        }

        c.Header("Access-Control-Allow-Origin", allowOrigin)
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization,X-Trace-Id")
        c.Header("Access-Control-Expose-Headers", "Content-Length,X-Trace-Id")
        c.Header("Access-Control-Max-Age", "43200")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}
```

- [ ] **Step 5: 实现 Auth / RBAC / RateLimit 骨架**

`api/middleware/auth.go`:
```go
package middleware

import "github.com/gin-gonic/gin"

// Auth JWT + OAuth2 认证中间件 — Phase 1 为骨架，Phase 2 实现
func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // TODO Phase 2: JWT 验证 + Casbin 加载角色
        c.Next()
    }
}
```

`api/middleware/rbac.go`:
```go
package middleware

import "github.com/gin-gonic/gin"

// RBAC Casbin 权限中间件 — Phase 1 为骨架，Phase 2 实现
func RBAC() gin.HandlerFunc {
    return func(c *gin.Context) {
        // TODO Phase 2: Casbin enforce
        c.Next()
    }
}
```

`api/middleware/ratelimit.go`:
```go
package middleware

import "github.com/gin-gonic/gin"

// RateLimit 限流中间件 — Phase 1 为骨架，后续 Phase 实现
func RateLimit() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
    }
}
```

- [ ] **Step 6: 编译验证**

```bash
go mod tidy
go build ./api/middleware/...
```

Expected: 无错误。

- [ ] **Step 7: Commit**

```bash
git add api/middleware/ go.mod go.sum
git commit -m "feat(phase1): add API middleware (recovery, logger, tracing, cors, auth/rbac/ratelimit stubs)"
```

---

### Task 10: 路由注册 + 入口

**Files:**
- Create: `api/router.go`
- Create: `cmd/api/main.go`

**Interfaces:**
- Produces: `func SetupRouter(cfg *config.AppConfig) *gin.Engine`
- Produces: `main()` — 加载配置 → 初始化 infra → 启动 Gin

- [ ] **Step 1: 实现路由注册**

`api/router.go`:
```go
package api

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/phoenix-agent-go/api/middleware"
    "github.com/phoenix-agent-go/infra/config"
    "github.com/phoenix-agent-go/infra/response"
)

func SetupRouter(cfg *config.AppConfig) *gin.Engine {
    if cfg.Server.Mode == "release" {
        gin.SetMode(gin.ReleaseMode)
    }

    r := gin.New()

    // 全局中间件
    r.Use(middleware.Recovery())
    r.Use(middleware.Logger())
    r.Use(middleware.Tracing())
    r.Use(middleware.CORS(&cfg.Cors))

    // 健康检查
    r.GET("/echo", func(c *gin.Context) {
        response.Success(c, "ok")
    })

    // 静态文件（头像等）
    r.Static("/api/upload", "./storage/upload")

    // 认证路由（无需 JWT）
    {
        auth := r.Group("/api/privilege/auth")
        auth.Use(middleware.RateLimit())
        // Phase 2: auth.POST("/login", handler.Login)
        // Phase 2: auth.POST("/logout", handler.Logout)
        // Phase 2: auth.GET("/captcha", handler.Captcha)
        _ = auth
    }

    // API 路由（需 JWT）
    api := r.Group("/api")
    api.Use(middleware.Auth())
    api.Use(middleware.RBAC())
    {
        // Phase 4: agentGroup := api.Group("/agent")
        // Phase 5: datasourceGroup := api.Group("/datasource")
        // Phase 5: chatGroup := api.Group("")
        _ = api
    }

    // 平台管理路由
    platform := r.Group("/platform")
    platform.Use(middleware.Auth())
    {
        // Phase 3: platform handler registration
        _ = platform
    }

    // 404
    r.NoRoute(func(c *gin.Context) {
        c.JSON(http.StatusNotFound, response.Response{
            Code:    404,
            Message: "not found",
        })
    })

    return r
}
```

- [ ] **Step 2: 实现 main.go**

`cmd/api/main.go`:
```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/phoenix-agent-go/api"
    "github.com/phoenix-agent-go/infra/config"
    "github.com/phoenix-agent-go/infra/logger"
    "github.com/phoenix-agent-go/infra/monitoring"
    "go.uber.org/zap"
)

func main() {
    // 加载配置
    cfg, err := config.Load("api")
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    // 初始化日志
    if err := logger.Init(&cfg.Monitor); err != nil {
        log.Fatalf("failed to init logger: %v", err)
    }
    defer logger.Sync()

    zap.L().Info("starting Phoenix API server",
        zap.Int("port", cfg.Server.Port),
        zap.String("version", cfg.Monitor.ServiceVersion),
    )

    // 初始化 OpenTelemetry
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    tp, err := monitoring.InitTracer(ctx, &cfg.Monitor)
    if err != nil {
        zap.L().Warn("failed to init tracer", zap.Error(err))
    } else {
        defer func() {
            shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer shutdownCancel()
            tp.Shutdown(shutdownCtx)
        }()
    }

    // 设置路由
    router := api.SetupRouter(cfg)

    // 启动服务
    srv := &http.Server{
        Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
        Handler:        router,
        ReadTimeout:    cfg.Server.ReadTimeout,
        WriteTimeout:   cfg.Server.WriteTimeout,
        MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
    }

    // 优雅关闭
    go func() {
        zap.L().Info("server listening", zap.String("addr", srv.Addr))
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("server error: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    zap.L().Info("shutting down server...")
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer shutdownCancel()

    if err := srv.Shutdown(shutdownCtx); err != nil {
        zap.L().Fatal("server forced to shutdown", zap.Error(err))
    }

    zap.L().Info("server stopped")
}
```

- [ ] **Step 3: 编译项目**

```bash
go mod tidy
go build ./...
```

Expected: 零错误。

- [ ] **Step 4: 启动验证**

```bash
go run ./cmd/api &
sleep 2
curl http://localhost:8066/echo
```

Expected: `{"code":0,"message":"success","data":"ok"}`

- [ ] **Step 5: Commit**

```bash
git add api/router.go cmd/api/main.go go.mod go.sum
git commit -m "feat(phase1): add router and main entry point with graceful shutdown"
```

---

### Task 11: Docker + docker-compose + Makefile

**Files:**
- Create: `Dockerfile`
- Create: `docker-compose.yaml`
- Create: `Makefile`

- [ ] **Step 1: Dockerfile**

`Dockerfile`:
```dockerfile
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /phoenix-api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata python3 py3-pip
COPY --from=builder /phoenix-api /usr/local/bin/phoenix-api
COPY configs/ /etc/phoenix/configs/
WORKDIR /app
EXPOSE 8066
CMD ["phoenix-api"]
```

- [ ] **Step 2: docker-compose.yaml**

`docker-compose.yaml`:
```yaml
version: "3.8"

services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8066:8066"
    environment:
      - PHOENIX_DATABASE_HOST=postgres
      - PHOENIX_REDIS_ADDR=redis:6379
      - PHOENIX_MILVUS_ADDR=milvus:19530
      - PHOENIX_RABBITMQ_ADDR=amqp://guest:guest@rabbitmq:5672/
      - PHOENIX_MONITOR_OTEL_ENDPOINT=jaeger:4317
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    restart: unless-stopped
    networks:
      - phoenix-net

  postgres:
    image: pgvector/pgvector:pg16
    environment:
      POSTGRES_USER: phoenix
      POSTGRES_PASSWORD: phoenix
      POSTGRES_DB: phoenix
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ../sql/all_schema.sql:/docker-entrypoint-initdb.d/01_schema.sql
      - ../sql/all_data.sql:/docker-entrypoint-initdb.d/02_data.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U phoenix"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - phoenix-net

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - phoenix-net

  milvus:
    image: milvusdb/milvus:v2.4.0
    ports:
      - "19530:19530"
      - "9091:9091"
    environment:
      ETCD_ENDPOINTS: etcd:2379
      MINIO_ADDRESS: minio:9000
    depends_on:
      - etcd
      - minio
    command: ["milvus", "run", "standalone"]
    networks:
      - phoenix-net

  etcd:
    image: quay.io/coreos/etcd:v3.5.5
    environment:
      ETCD_AUTO_COMPACTION_MODE: revision
      ETCD_AUTO_COMPACTION_RETENTION: "1000"
      ETCD_QUOTA_BACKEND_BYTES: "4294967296"
    command: etcd -advertise-client-urls=http://127.0.0.1:2379 -listen-client-urls http://0.0.0.0:2379
    networks:
      - phoenix-net

  minio:
    image: minio/minio:latest
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    command: minio server /data --console-address ":9001"
    networks:
      - phoenix-net

  rabbitmq:
    image: rabbitmq:3-management-alpine
    ports:
      - "5672:5672"
      - "15672:15672"
    healthcheck:
      test: ["CMD", "rabbitmq-diagnostics", "check_port_connectivity"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - phoenix-net

  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"
      - "4317:4317"
    environment:
      COLLECTOR_OTLP_ENABLED: "true"
    networks:
      - phoenix-net

volumes:
  pgdata:

networks:
  phoenix-net:
    driver: bridge
```

- [ ] **Step 3: Makefile**

`Makefile`:
```makefile
.PHONY: build run test lint clean docker-build docker-up docker-down migrate

APP_NAME := phoenix-api
GO := go
GOFMT := gofmt

build:
	$(GO) build -o bin/$(APP_NAME) ./cmd/api

run:
	$(GO) run ./cmd/api

test:
	$(GO) test -v -race -coverprofile=coverage.out ./...

test-short:
	$(GO) test -short ./...

lint:
	golangci-lint run ./...

fmt:
	$(GOFMT) -s -w .

vet:
	$(GO) vet ./...

clean:
	rm -rf bin/ coverage.out

docker-build:
	docker build -t phoenix-api .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f api

migrate-up:
	migrate -path migrations -database "postgres://phoenix:phoenix@localhost:5432/phoenix?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://phoenix:phoenix@localhost:5432/phoenix?sslmode=disable" down

mod:
	$(GO) mod tidy

all: fmt vet test build
```

- [ ] **Step 4: 编译验证**

```bash
go build ./...
```

- [ ] **Step 5: Commit**

```bash
git add Dockerfile docker-compose.yaml Makefile
git commit -m "feat(phase1): add Dockerfile, docker-compose, and Makefile"
```

---

### Task 12: 集成验证 + 回测

**目标**：全量编译、启动、验证 `/echo` 端点

- [ ] **Step 1: 全量编译**

```bash
cd d:/GitHub/Phoenix-Agent-Java/godata
go build ./...
```

Expected: 无错误。

- [ ] **Step 2: 静态检查**

```bash
go vet ./...
```

Expected: 无警告。

- [ ] **Step 3: 运行所有测试**

```bash
go test ./... -v -short
```

Expected: 所有 PASS。

- [ ] **Step 4: 启动服务 + 验证 /echo**

```bash
# Terminal 1: 启动服务
go run ./cmd/api

# Terminal 2: 测试
curl -v http://localhost:8066/echo
```

Expected response:
```json
{"code":0,"message":"success","data":"ok"}
```

- [ ] **Step 5: 验证中间件（trace id）**

```bash
curl -v http://localhost:8066/echo -H "X-Trace-Id: test-trace-123"
```

Expected: 响应头含 `X-Trace-Id: test-trace-123`，日志输出 trace_id。

- [ ] **Step 6: 验证配置加载**

检查启动日志中是否输出：
```
INFO  starting Phoenix API server  {"port": 8066, "version": "2.0.0"}
```

- [ ] **Step 7: 验证优雅关闭**

```bash
# 启动服务
go run ./cmd/api &
SERVER_PID=$!

# 发送 SIGTERM
kill -TERM $SERVER_PID

# 检查日志输出
```

Expected: 日志输出 `shutting down server...` 和 `server stopped`。

- [ ] **Step 8: 验证 Docker 构建**

```bash
docker build -t phoenix-api .
```

Expected: 构建成功，输出 `Successfully tagged phoenix-api:latest`。

- [ ] **Step 9: Commit**

```bash
git add -A
git commit -m "feat(phase1): integration verification passed — build, test, run, /echo ok"
```

---

## 里程碑

| 里程碑 | Task | 验收标准 |
|:---|:---|:---|
| M1: 项目骨架 | 1-2 | `go build ./...` 无错误，目录结构完整 |
| M2: 配置体系 | 3, 8 | Viper 加载所有 YAML，测试通过 |
| M3: 基础设施包 | 4-7 | 所有 `infra/` 包可编译，单元测试通过 |
| M4: HTTP 骨架 | 9-10 | Gin 启动 → `/echo` 返回正确响应 |
| M5: 部署就绪 | 11 | `docker build` + `docker-compose up` 成功 |
| M6: Phase 1 完成 | 12 | 全量回测通过 |

## 回退策略

Phase 1 不涉及 Java 端的任何变更，回退无需操作 Java。若 Go 项目骨架出现问题：

```bash
# 完全重置 godata/ 目录
cd d:/GitHub/Phoenix-Agent-Java
rm -rf godata/*
git checkout HEAD -- godata/go.mod godata/go.sum
mv godata/_archive/main.go godata/
mv godata/_archive/trpc-agent-go-使用说明.md godata/
```

每个 Task 独立 commit，可通过 `git revert <commit>` 撤销单个 Task。
