# issue2md - 任务列表

## Version: 1.0
## Date: 2026-03-31

---

## Phase 1: Foundation (数据结构定义)

### 测试任务

1. **[P]** 创建 `internal/github/types_test.go` - 为数据结构定义测试用例
   - 表格驱动测试
   - 测试 Issue、Comment、User、Label、Milestone、Reaction 结构
   - 测试 JSON 序列化/反序列化
   - 测试结构验证逻辑

2. **[P]** 创建 `internal/parser/url_test.go` - 为 URL 解析定义测试用例
   - 表格驱动测试
   - 测试有效 URL 解析
   - 测试各种无效 URL 格式
   - 测试 URL 验证逻辑

3. **[P]** 创建 `internal/cli/config_test.go` - 为 CLI 配置定义测试用例
   - 表格驱动测试
   - 测试参数解析
   - 测试配置验证
   - 测试帮助信息生成

4. **[P]** 创建 `internal/config/config_test.go` - 为应用配置定义测试用例
   - 表格驱动测试
   - 测试从环境变量加载配置
   - 测试默认值设置
   - 测试配置验证

### 实现任务

5. **[P]** 创建 `internal/github/types.go` - 定义核心数据结构
   - 实现 Issue、Comment、User、Label、Milestone、Reaction 结构体
   - 添加 JSON 标签
   - 实现结构验证方法

6. **[P]** 创建 `internal/github/const.go` - 定义常量
   - 定义默认超时时间
   - 定义重试配置
   - 定义 HTTP Headers
   - 定义 API 端点常量

7. **[P]** 创建 `internal/parser/url.go` - 实现 URL 解析
   - 实现正则表达式匹配
   - 实现 ParseIssueURL 函数
   - 实现 Validate 方法
   - 实现 String 方法

8. **[P]** 创建 `internal/cli/config.go` - 实现 CLI 配置
   - 实现配置结构体
   - 实现参数解析
   - 实现验证逻辑
   - 实现帮助信息

9. **[P]** 创建 `internal/cli/usage.go` - 实现帮助和版本信息
   - 实现使用说明
   - 实现版本字符串
   - 格式化输出

10. **[P]** 创建 `internal/config/config.go` - 实现应用配置
    - 实现配置结构体
    - 实现配置加载逻辑
    - 实现环境变量加载
    - 实现验证逻辑

---

## Phase 2: GitHub Fetcher (API交互逻辑，TDD)

### 测试任务

11. **[P]** 创建 `internal/github/client_test.go` - 为 GitHub 客户端定义测试用例
    - 表格驱动测试
    - 使用 httptest.Server 模拟 API
    - 测试 HTTP 客户端配置
    - 测试请求构建和响应解析
    - 测试错误处理（404、403、429、5xx）
    - 测试重试逻辑
    - 测试速率限制处理

12. **[P]** 创建 `internal/github/retry_test.go` - 为重试逻辑定义测试用例
    - 表格驱动测试
    - 测试指数退避算法
    - 测试重试条件判断
    - 测试重试计数

### 实现任务

13. **[P]** 创建 `internal/github/retry.go` - 实现重试逻辑
    - 实现指数退避算法
    - 实现重试条件判断
    - 实现重试包装器
    - 添加随机抖动

14. **[P]** 创建 `internal/github/client.go` - 实现 GitHub API 客户端
    - 实现客户端结构体
    - 实现构造函数
    - 实现基本 HTTP 请求方法
    - 实现错误处理和包装
    - 实现速率限制处理
    - 实现分页处理

15. **[P]** 创建 `internal/github/fetch.go` - 实现 API 数据获取
    - 实现 FetchIssue 方法
    - 实现 FetchComments 方法
    - 实现 FetchAll 方法
    - 实现响应解析
    - 实现分页逻辑

---

## Phase 3: Markdown Converter (转换逻辑，TDD)

### 测试任务

16. **[P]** 创建 `internal/converter/html_test.go` - 为 HTML 转换定义测试用例
    - 表格驱动测试
    - Golden 文件测试
    - 测试基本 HTML 元素转换
    - 测试代码块保留
    - 测试链接转换
    - 测试列表转换
    - 测试表格转换
    - 测试 GitHub 特殊元素

17. **[P]** 创建 `internal/converter/markdown_test.go` - 为 Markdown 转换定义测试用例
    - 表格驱动测试
    - Golden 文件测试
    - 测试 Issue 转换
    - 测试评论转换
    - 测试元数据格式化
    - 测试反应信息格式化
    - 测试完整文档生成

### 实现任务

18. **[P]** 创建 `internal/converter/html.go` - 实现 HTML 解析和转换
    - 实现流式 HTML 解析
    - 实现基本元素转换函数
    - 实现代码块处理
    - 实现链接处理
    - 实现列表处理
    - 实现表格处理
    - 实现 GitHub 扩展

19. **[P]** 创建 `internal/converter/markdown.go` - 实现 Markdown 转换器
    - 实现 Converter 结构体
    - 实现 Convert 方法
    - 实现 ConvertIssue 方法
    - 实现 ConvertComment 方法
    - 实现格式化辅助函数
    - 实现元数据生成
    - 实现反应信息生成

20. **[P]** 创建 `internal/converter/utils.go` - 实现转换工具函数
    - 实现 Markdown 转义函数
    - 实现时间格式化函数
    - 实现字符串处理函数

---

## Phase 4: CLI Assembly (命令行入口集成)

### 测试任务

21. **[P]** 创建 `cmd/issue2md/main_test.go` - 为主程序定义测试用例
    - 集成测试
    - 测试命令行参数解析
    - 测试完整工作流程
    - 测试错误场景
    - 测试 verbose 模式

### 实现任务

22. **[P]** 创建 `cmd/issue2md/main.go` - 实现 CLI 入口
    - 实现 main 函数
    - 实现主逻辑流程
    - 集成所有包
    - 实现错误处理
    - 实现退出码处理
    - 实现 verbose 日志

23. **[P]** 创建 `Makefile` - 实现构建脚本
    - 实现构建目标
    - 实现测试目标
    - 实现覆盖率目标
    - 实现交叉编译
    - 实现清理目标

24. **[P]** 创建 `go.mod` - 初始化 Go 模块
    - 设置模块路径
    - 设置 Go 版本要求
    - 添加必要的依赖

---

## Phase 5: Polishing (完善和发布)

### 测试任务

25. **[P]** 创建 `internal/testdata/setup.go` - 测试辅助代码
    - 创建测试数据生成器
    - 创建测试 Issue 模板
    - 创建测试用例数据

26. **[P]** 创建 `internal/testdata/expected/simple.md` - 创建简单的预期输出
27. **[P]** 创建 `internal/testdata/expected/complex.md` - 创建复杂的预期输出
28. **[P]** 创建 `internal/testdata/expected/with_comments.md` - 创建带评论的预期输出

### 实现任务

29. **[P]** 更新 `README.md` - 创建项目文档
    - 项目描述
    - 安装说明
    - 使用示例
    - 命令行选项
    - 限制说明
    - 贡献指南

30. **[P]** 创建 `CHANGELOG.md` - 创建变更日志
    - 初始版本记录
    - 未来规划

31. **[P]** 创建 `.gitignore` - 完善 Git 忽略规则
    - 添加测试数据
    - 添加构建产物
    - 添加 IDE 文件

---

## 任务依赖关系

### Phase 1 前置任务
- 任务 24 (go.mod) 必须在创建任何 Go 文件之前完成

### Phase 2 依赖
- 依赖于 Phase 1 中所有任务（特别是类型定义）

### Phase 3 依赖
- 依赖于 Phase 1 中的类型定义
- 需要正确理解 Phase 2 中的 GitHub 数据结构

### Phase 4 依赖
- 依赖于前 3 个阶段的所有实现任务
- 需要完整的集成测试

### Phase 5 依赖
- 依赖于前 4 个阶段的任务完成
- 是最后的完善和文档工作

---

## 执行策略

### 并行执行标记说明
- **[P]** - 可以并行执行，没有依赖关系

### 测试先行策略
1. 每个实现任务必须先完成对应的测试任务
2. 测试应该先写失败用例（Red）
3. 然后实现代码使测试通过（Green）
4. 最后重构（Refactor）

### TDD 验证点
- 所有单元测试必须通过
- 覆盖率必须 > 85%
- 使用 `go test -v` 验证详细输出
- 使用 `go test -cover` 验证覆盖率

---

## 任务完成检查清单

### 每个任务完成后检查
- [ ] 测试文件已创建并测试通过
- [ ] 实现代码满足测试用例
- [ ] 代码格式正确 (`gofmt`)
- [ ] 无警告 (`go vet`)
- [ ] 错误处理完善（所有错误都被处理）

### 每个阶段完成后检查
- [ ] 所有任务完成
- [ ] 集成测试通过
- [ ] 文档更新
- [ ] 代码审查通过

### 项目完成后检查
- [ ] 所有阶段任务完成
- [ ] 整体测试覆盖率 > 85%
- [ ] 文档完整
- [ ] 发布准备就绪

---

**文档状态:** 草稿
**最后更新:** 2026-03-31
**下次更新:** 首个任务开始前