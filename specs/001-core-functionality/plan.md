# issue2md 技术实现方案

## Version: 1.0
## Date: 2026-03-31

---

## 1. 技术上下文总结

### 1.1 技术栈

| 组件 | 技术选型 | 理由 |
|------|----------|------|
| 编程语言 | Go 1.24+ | 项目宪法明确要求 |
| CLI 框架 | 标准库 `flag` | 遵循"标准库优先"原则 |
| Web 框架 | 标准库 `net/http` | 遵循"简单性原则" |
| HTTP 客户端 | 标准库 `net/http` | 避免不必要的依赖 |
| GitHub API | REST API (v3) | 公开仓库无需认证，简单直接 |
| HTML 解析 | `golang.org/x/net/html` | 官方推荐的 HTML 解析库 |
| Markdown 转换 | 自定义实现 | 避免引入第三方库，遵循"标准库优先" |
| 数据存储 | 无 (内存) | 实时获取，无需持久化 |
| 日志 | 标准库 `log` | 简单直接，输出到 stderr |
| 测试 | 标准库 `testing` | 表格驱动测试，集成测试优先 |

### 1.2 技术选型说明

**为什么不使用 `google/go-github`？**
- 项目宪法要求"标准库优先"
- 对于公开仓库，REST API 足够简单
- 避免引入非必需的依赖
- 减少攻击面，简化依赖管理

**为什么不使用 GraphQL？**
- REST API 对于公开仓库已完全满足需求
- GraphQL 需要额外的查询构建逻辑
- 增加 API 复杂度，违反"简单性原则"

**为什么不使用第三方 Markdown 库？**
- HTML 到 Markdown 的转换逻辑相对直接
- 自定义实现可精确控制输出格式
- 符合 spec 中要求的特定输出结构

---

## 2. "合宪性"审查

### 2.1 对照 `constitution.md` 逐条审查

| 宪法条款 | 要求 | 本方案合规性 | 说明 |
|----------|------|-------------|------|
| **第一条：简单性原则** | | | |
| 1.1 (YAGNI) | 只实现 spec.md 中明确要求的功能 | ✅ 合规 | 方案严格遵循 spec.md，未添加额外功能 |
| 1.2 (标准库优先) | 优先使用 Go 标准库 | ✅ 合规 | 所有核心组件均使用标准库 |
| 1.3 (反过度工程) | 简单优于复杂 | ✅ 合规 | 无不必要抽象，数据结构简洁 |
| **第二条：测试先行铁律** | | | |
| 2.1 (TDD循环) | Red-Green-Refactor | ✅ 合规 | 方案明确要求 TDD 流程 |
| 2.2 (表格驱动) | 单元测试优先采用表格驱动 | ✅ 合规 | 测试策略明确要求表格驱动测试 |
| 2.3 (拒绝Mocks) | 优先集成测试，使用真实依赖 | ✅ 合规 | 使用 httptest.Server 而非 Mock 框架 |
| **第三条：明确性原则** | | | |
| 3.1 (错误处理) | 所有错误显式处理，使用 fmt.Errorf 包装 | ✅ 合规 | 错误类型设计支持包装和上下文 |
| 3.2 (无全局变量) | 依赖通过函数参数或结构体显式注入 | ✅ 合规 | Client 和 Config 均通过构造函数注入 |

### 2.2 合规性总结

**本方案 100% 符合项目宪法要求。**

---

## 3. 项目结构细化

### 3.1 目录结构

```
issue2md/
├── cmd/
│   ├── issue2md/
│   │   └── main.go              # CLI 入口
│   └── issue2mdweb/
│       └── main.go              # Web 入口（预留）
├── internal/
│   ├── github/
│   │   ├── client.go             # GitHub API 客户端
│   │   ├── client_test.go
│   │   ├── types.go             # 数据结构定义
│   │   └── retry.go             # 重试逻辑
│   ├── parser/
│   │   ├── url.go               # URL 解析
│   │   └── url_test.go
│   ├── converter/
│   │   ├── markdown.go          # Markdown 转换器
│   │   ├── markdown_test.go
│   │   ├── html.go              # HTML 解析和转换
│   │   └── html_test.go
│   ├── cli/
│   │   ├── config.go            # CLI 配置解析
│   │   ├── config_test.go
│   │   └── usage.go            # 帮助信息
│   └── config/
│       ├── config.go             # 应用配置加载
│       └── config_test.go
├── web/
│   ├── templates/               # Web 模板（预留）
│   └── static/                 # 静态资源（预留）
├── specs/
│   ├── spec.md
│   └── 001-core-functionality/
│       ├── api-sketch.md
│       └── plan.md
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── constitution.md
```

### 3.2 包职责与依赖关系

#### `internal/github`
**职责：**
- 与 GitHub REST API 交互
- 处理 HTTP 请求和响应
- 实现重试和速率限制处理
- 解析 JSON 响应

**依赖：**
- 标准库：`net/http`, `encoding/json`, `time`, `context`, `log`

**导出接口：**
```go
type Client interface {
    FetchIssue(ctx context.Context, owner, repo string, number int) (*Issue, error)
    FetchComments(ctx context.Context, owner, repo string, number int) ([]*Comment, error)
    FetchAll(ctx context.Context, owner, repo string, number int) (*Issue, []*Comment, error)
}
```

**被依赖方：** `cmd/issue2md`

---

#### `internal/parser`
**职责：**
- 解析 GitHub Issue URL
- 验证 URL 格式
- 提取 owner、repo、number

**依赖：**
- 标准库：`regexp`, `net/url`, `fmt`

**导出接口：**
```go
func ParseIssueURL(rawURL string) (*IssueURL, error)
```

**被依赖方：** `cmd/issue2md`

---

#### `internal/converter`
**职责：**
- 将 Issue 和 Comment 数据转换为 Markdown
- HTML 到 Markdown 的转换
- 格式化元数据、评论等

**依赖：**
- 标准库：`strings`, `fmt`, `time`
- 扩展库：`golang.org/x/net/html`
- 内部：`internal/github`（仅使用其类型）

**导出接口：**
```go
type Converter interface {
    Convert(issue *github.Issue, comments []*github.Comment) string
}
```

**被依赖方：** `cmd/issue2md`

---

#### `internal/cli`
**职责：**
- 解析命令行参数
- 验证参数有效性
- 提供帮助和使用信息

**依赖：**
- 标准库：`flag`, `fmt`, `os`

**导出接口：**
```go
func ParseArgs(args []string) (*Config, error)
func Usage() string
```

**被依赖方：** `cmd/issue2md`

---

#### `internal/config`
**职责：**
- 从环境变量加载配置
- 提供默认配置
- 验证配置有效性

**依赖：**
- 标准库：`os`, `time`, `strconv`

**导出接口：**
```go
func Load() *Config
func LoadFromEnv() *Config
```

**被依赖方：** `internal/github`

---

### 3.3 依赖关系图

```
cmd/issue2md/main.go
    ├── internal/cli          (参数解析)
    ├── internal/parser       (URL 解析)
    ├── internal/github       (API 客户端)
    │       └── internal/config (配置加载)
    └── internal/converter   (Markdown 转换)
            └── internal/github (类型引用)
```

**原则验证：**
- ✅ 无循环依赖
- ✅ 依赖方向清晰（依赖注入）
- ✅ 每个包职责单一
- ✅ 包边界明确

---

## 4. 核心数据结构

### 4.1 `internal/github/types.go`

```go
// Package github 定义与 GitHub API 交互的数据结构
package github

import "time"

// Issue 表示一个 GitHub Issue
type Issue struct {
    Number      int
    Title       string
    Body        string
    State       string // "open" 或 "closed"
    HTMLURL     string
    Author      User
    CreatedAt   time.Time
    UpdatedAt   time.Time
    ClosedAt    *time.Time
    Labels      []Label
    Milestone   *Milestone
    Reactions   []Reaction
    CommentsURL string
}

// Comment 表示 Issue 上的评论
type Comment struct {
    ID        int
    Body      string
    HTMLURL   string
    Author    User
    CreatedAt time.Time
    UpdatedAt time.Time
    Reactions []Reaction
    CommitSHA string // 可选，用于 PR 评论
}

// User 表示 GitHub 用户
type User struct {
    Login     string
    ID        int64
    AvatarURL string
    HTMLURL   string
}

// Label 表示 Issue 标签
type Label struct {
    Name        string
    Color       string
    Description string
}

// Milestone 表示 Issue 里程碑
type Milestone struct {
    Title       string
    Number      int
    State       string
    Description string
    HTMLURL     string
}

// Reaction 表示对 Issue 或评论的反应
type Reaction struct {
    Content string // "+1", "-1", "laugh", "hooray", "confused", "heart", "rocket", "eyes"
    Count   int
    Users   []string // 用户登录名列表（可选，用于详细输出）
}

// IssueURL 表示解析后的 GitHub Issue URL
type IssueURL struct {
    Owner   string
    Repo    string
    Number  int
    Original string
}
```

### 4.2 `internal/cli/config.go`

```go
// Package cli 处理命令行参数
package cli

// Config 表示 CLI 配置
type Config struct {
    URL     string
    Verbose bool
    Version bool
    Help    bool
}
```

### 4.3 `internal/config/config.go`

```go
// Package config 管理应用配置
package config

import "time"

// Config 表示应用配置
type Config struct {
    // GitHub API 设置
    BaseURL  string
    Timeout  time.Duration

    // 重试设置
    MaxRetries    int
    InitialBackoff time.Duration
    MaxBackoff    time.Duration

    // 日志设置
    Verbose bool
}
```

### 4.4 `internal/converter/markdown.go`

```go
// Package converter 处理 Markdown 转换
package converter

import "github.com/bigwhite/issue2md/internal/github"

// Converter 处理 GitHub 数据到 Markdown 的转换
type Converter struct {
    includeReactions bool
    includeMetadata  bool
}
```

---

## 5. 接口设计

### 5.1 `internal/github` 包接口

```go
// Package github - GitHub API 客户端接口

// Client 表示 GitHub API 客户端
type Client struct {
    httpClient *http.Client
    baseURL    string
    verbose    bool
    maxRetries int
}

// NewClient 创建新的 GitHub API 客户端
func NewClient(httpClient *http.Client, verbose bool, maxRetries int) *Client

// FetchIssue 获取单个 Issue
// 如果 Issue 不存在或仓库是私有的，返回错误
func (c *Client) FetchIssue(ctx context.Context, owner, repo string, number int) (*Issue, error)

// FetchComments 获取 Issue 的所有评论
// 自动处理分页以获取所有评论（包括折叠的）
func (c *Client) FetchComments(ctx context.Context, owner, repo string, number int) ([]*Comment, error)

// FetchAll 一次性获取 Issue 和所有评论
// 这是推荐的典型使用方式
func (c *Client) FetchAll(ctx context.Context, owner, repo string, number int) (*Issue, []*Comment, error)

// Close 关闭客户端并释放资源
func (c *Client) Close() error
```

### 5.2 `internal/converter` 包接口

```go
// Package converter - Markdown 转换器接口

// Converter 处理 GitHub 数据到 Markdown 的转换
type Converter struct{}

// NewConverter 创建新的 Markdown 转换器
func NewConverter() *Converter

// Convert 将 Issue 和评论转换为 Markdown
// 返回完整的 Markdown 文档
func (c *Converter) Convert(issue *github.Issue, comments []*github.Comment) string

// ConvertIssue 仅转换 Issue（不含评论）为 Markdown
func (c *Converter) ConvertIssue(issue *github.Issue) string

// ConvertComment 将单个评论转换为 Markdown
func (c *Converter) ConvertComment(comment *github.Comment) string
```

### 5.3 `internal/parser` 包接口

```go
// Package parser - URL 解析器接口

// IssueURL 表示解析后的 GitHub Issue URL
type IssueURL struct {
    Owner   string
    Repo    string
    Number  int
    Original string
}

// ParseIssueURL 解析 GitHub Issue URL
// 如果 URL 无效，返回错误
func ParseIssueURL(rawURL string) (*IssueURL, error)

// String 返回完整 URL 字符串
func (u *IssueURL) String() string

// Validate 验证 URL 组件
func (u *IssueURL) Validate() error
```

### 5.4 `internal/cli` 包接口

```go
// Package cli - CLI 解析器接口

// Config 表示 CLI 配置
type Config struct {
    URL     string
    Verbose bool
    Version bool
    Help    bool
}

// ParseArgs 解析命令行参数
func ParseArgs(args []string) (*Config, error)

// Validate 验证 CLI 配置
func (c *Config) Validate() error

// Usage 返回使用说明
func Usage() string

// Version 返回版本字符串
func Version() string
```

### 5.5 `internal/config` 包接口

```go
// Package config - 配置管理器接口

// Config 表示应用配置
type Config struct {
    BaseURL        string
    Timeout        time.Duration
    MaxRetries     int
    InitialBackoff time.Duration
    MaxBackoff     time.Duration
    Verbose        bool
}

// Load 加载配置（环境变量 + 默认值）
func Load() *Config

// LoadFromEnv 从环境变量加载配置
// 前缀: ISSUE2MD_
// 示例:
//   ISSUE2MD_TIMEOUT=30s
//   ISSUE2MD_VERBOSE=true
func LoadFromEnv() *Config

// Validate 验证配置
func (c *Config) Validate() error
```

---

## 6. 实现流程

### 6.1 典型执行流程

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLI 入口                            │
│                    cmd/issue2md/main.go                     │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
              ┌─────────────────────────┐
              │  1. 解析命令行参数    │
              │  internal/cli.ParseArgs │
              └─────────────────────────┘
                            │
                            ▼
              ┌─────────────────────────┐
              │  2. 解析 GitHub URL   │
              │  internal/parser.Parse  │
              └─────────────────────────┘
                            │
                            ▼
              ┌─────────────────────────┐
              │  3. 加载配置         │
              │  internal/config.Load   │
              └─────────────────────────┘
                            │
                            ▼
              ┌─────────────────────────┐
              │  4. 创建 GitHub 客户端 │
              │  internal/github.New   │
              └─────────────────────────┘
                            │
                            ▼
              ┌─────────────────────────┐
              │  5. 获取 Issue 和评论  │
              │  client.FetchAll()     │
              └─────────────────────────┘
                            │
                            ▼
              ┌─────────────────────────┐
              │  6. 转换为 Markdown  │
              │  converter.Convert()    │
              └─────────────────────────┘
                            │
                            ▼
              ┌─────────────────────────┐
              │  7. 输出到 stdout     │
              │  fmt.Println()         │
              └─────────────────────────┘
```

### 6.2 错误处理流程

```
┌─────────────┐
│ 任何错误    │
└─────────────┘
       │
       ▼
┌─────────────────────────┐
│ 检查错误类型           │
│ - ParseError (URL 无效) │
│ - APIError (API 错误)   │
│ - NetworkError (网络错误)  │
│ - RateLimit (速率限制)    │
└─────────────────────────┘
       │
       ▼
┌─────────────────────────┐
│ 使用 fmt.Errorf 包装    │
│ 保留错误上下文          │
└─────────────────────────┘
       │
       ▼
┌─────────────────────────┐
│ 输出清晰的错误信息      │
│ 到 stderr             │
└─────────────────────────┘
       │
       ▼
┌─────────────────────────┐
│ 使用适当的退出码退出    │
│ 0: 成功               │
│ 1: 一般错误            │
│ 2: 参数错误            │
│ 3: 速率限制            │
└─────────────────────────┘
```

---

## 7. 测试策略

### 7.1 测试优先级（TDD 顺序）

1. **第一阶段：URL 解析测试** (`internal/parser`)
   - 表格驱动测试
   - 覆盖所有有效和无效情况
   - 集成测试：无（纯逻辑）

2. **第二阶段：配置测试** (`internal/config`, `internal/cli`)
   - 表格驱动测试
   - 测试参数解析和配置加载
   - 集成测试：无（纯逻辑）

3. **第三阶段：HTML 转换测试** (`internal/converter`)
   - 表格驱动测试
   - Golden 文件测试（预期输出比对）
   - 集成测试：无（纯逻辑）

4. **第四阶段：GitHub 客户端测试** (`internal/github`)
   - 表格驱动测试
   - 使用 `httptest.Server` 模拟 API
   - 集成测试：可选（真实 API，条件执行）

5. **第五阶段：端到端测试** (`cmd/issue2md`)
   - 集成测试
   - 测试完整流程
   - 需要真实 GitHub Issue

### 7.2 测试覆盖率目标

| 包 | 目标覆盖率 | 说明 |
|----|----------|------|
| `internal/parser` | 100% | 纯逻辑，应完全覆盖 |
| `internal/config` | 100% | 纯逻辑，应完全覆盖 |
| `internal/cli` | 100% | 纯逻辑，应完全覆盖 |
| `internal/converter` | 100% | 纯逻辑，应完全覆盖 |
| `internal/github` | 90%+ | 包含网络逻辑，部分场景难以测试 |
| `cmd/issue2md` | 80%+ | 集成测试为主 |

**总体目标：> 85%**

### 7.3 测试示例结构

```go
// internal/parser/url_test.go

func TestParseIssueURL(t *testing.T) {
    tests := []struct {
        name    string
        url     string
        want    *IssueURL
        wantErr bool
    }{
        {
            name:    "valid URL",
            url:     "https://github.com/golang/go/issues/12345",
            want:    &IssueURL{Owner: "golang", Repo: "go", Number: 12345},
            wantErr: false,
        },
        {
            name:    "invalid URL - missing number",
            url:     "https://github.com/golang/go/issues",
            wantErr: true,
        },
        // ... 更多测试用例
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseIssueURL(tt.url)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseIssueURL() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ParseIssueURL() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

## 8. 性能考虑

### 8.1 内存管理

**策略：**
- HTML 解析使用流式处理，避免构建完整 DOM
- 使用 `strings.Builder` 高效构建字符串
- 评论按批处理（非全部加载到内存）

**目标：**
- 典型 Issue（<100 评论）：< 10MB 内存
- 大型 Issue（1000 评论）：< 50MB 内存

### 8.2 网络性能

**策略：**
- 默认 HTTP 超时：30 秒
- 重试机制：最多 3 次，指数退避
- 连接复用：使用 HTTP 连接池

**目标：**
- 典型 Issue 获取：< 30 秒
- API 调用延迟：< 500ms（p95）

### 8.3 并发策略

**第一阶段（v1.0）：**
- 顺序获取 Issue 和评论
- 简单直接，易于调试

**未来优化：**
- 并行获取 Issue 和评论（减少总延迟）
- 使用 `sync.WaitGroup` 协调

---

## 9. 安全考虑

### 9.1 输入验证

- URL 解析前验证格式
- 拒绝非 `github.com` 域名
- 验证 owner、repo、number 的有效性

### 9.2 网络安全

- 使用 HTTPS（强制）
- 默认超时防止挂起
- 重试限制防止无限循环

### 9.3 输出安全

- 保留原始格式，不执行任何代码
- 用户提供的 HTML 经过净化（去除危险标签）
- 不存储敏感信息

---

## 10. 发布清单

### 10.1 代码质量

- [ ] 所有测试通过
- [ ] 测试覆盖率 > 85%
- [ ] `go vet` 无警告
- [ ] `gofmt` 格式化
- [ ] 代码审查通过

### 10.2 文档

- [ ] README.md 完整
- [ ] 使用示例清晰
- [ ] API 文档完整
- [ ] CHANGELOG.md 更新

### 10.3 构建和发布

- [ ] 交叉编译（macOS, Linux, Windows）
- [ ] 二进制文件签名（如需要）
- [ ] Release notes 准备
- [ ] GitHub Release 创建

---

## 11. 未来扩展点

### 11.1 已预留的扩展

1. **Web 接口** (`cmd/issue2mdweb/`)
   - 使用 `net/http` 实现
   - 模板目录已预留
   - 可选功能，不阻塞 v1.0

2. **私有仓库支持**
   - 可通过环境变量传入 Token
   - `internal/config` 已预留加载逻辑
   - HTTP Header 添加 Authorization

3. **输出格式扩展**
   - JSON 格式
   - HTML 格式
   - 可配置模板

### 11.2 技术债务跟踪

| 项目 | 优先级 | 说明 |
|------|--------|------|
| 并发获取 | 低 | 顺序获取已满足 v1.0 需求 |
| 缓存 | 低 | 单次运行场景，无需缓存 |
| PR 支持 | 低 | v1.0 不支持 |

---

## 附录 A: 关键常量定义

```go
// internal/github/const.go

const (
    // API 相关
    DefaultBaseURL = "https://api.github.com"
    APIVersion     = "v3"

    // 重试配置
    DefaultMaxRetries     = 3
    DefaultInitialBackoff = 1 * time.Second
    DefaultMaxBackoff     = 4 * time.Second

    // 超时配置
    DefaultTimeout = 30 * time.Second

    // 分页配置
    DefaultPerPage = 100

    // HTTP Headers
    HeaderUserAgent      = "User-Agent"
    HeaderRateLimit     = "X-RateLimit-Limit"
    HeaderRateRemaining = "X-RateLimit-Remaining"
    HeaderRateReset    = "X-RateLimit-Reset"
    HeaderRetryAfter    = "Retry-After"
)
```

---

## 附录 B: 错误码映射

| 场景 | 内部错误类型 | 退出码 |
|--------|-------------|--------|
| 成功 | 无 | 0 |
| URL 格式无效 | ParseError | 2 |
| 参数错误 | CLIError | 2 |
| Issue 不存在 | APIError(NotFound) | 1 |
| 仓库是私有的 | APIError(Forbidden) | 1 |
| 网络失败 | APIError(Network) | 1 |
| API 速率限制 | APIError(RateLimit) | 3 |
| 服务器错误 | APIError(ServerError) | 1 |
| 未知错误 | AppError | 1 |

---

**文档状态:** 草稿
**最后更新:** 2026-03-31
**下次审查:** Phase 1 完成后
