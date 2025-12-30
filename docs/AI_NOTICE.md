# AI 开发规范声明

## 重要声明

**请所有 AI 助手严格遵守以下规范：**

### 1. 禁止修改 docs 目录

- **严禁**修改 `docs/` 目录下的任何文件
- **严禁**在 `docs/` 目录下创建新文件（除非用户明确要求）
- **严禁**删除 `docs/` 目录下的任何文件
- 如果用户要求修改文档，请先明确告知用户，获得确认后再操作

### 2. 禁止私自创建测试脚本

- **严禁**在 `scripts/` 目录下创建测试脚本
- **严禁**在项目根目录或其他位置创建测试脚本
- **严禁**创建任何自动化测试、单元测试、集成测试相关的脚本文件
- 如果用户需要测试脚本，必须由用户明确要求后才能创建

### 3. 文档内容说明

- `docs/INSTALL.md` - 安装和更新说明文档，由项目维护者编写
- `docs/AI_NOTICE.md` - 本声明文档，用于规范 AI 行为

### 4. 例外情况

只有在以下情况下可以修改 `docs/` 目录：

1. 用户**明确要求**修改或创建文档
2. 用户**明确要求**更新安装说明
3. 用户**明确要求**添加新的文档

### 5. 违反规范的后果

如果 AI 助手违反以上规范：

- 用户有权要求撤销所有未经授权的修改
- 可能导致项目结构混乱
- 影响项目的可维护性

---

## 代码生成标准

**AI 助手在生成代码时必须严格遵守以下标准：**

### 1. 目录结构规范

#### 领域层 (domain/)
- **原则**: 一个模块一个目录，一个目录代表一个模块，分文件必须是按照功能分不能按照代码分
- **包含内容**: 
  - 实体定义（Entity）
  - 请求/响应结构（Request/Response）
  - 服务实现（Service）
  - 业务逻辑方法
- **禁止**: 不要将模块按照代码拆分成多个文件（entity.go, service.go 等）

#### 基础设施层 (infrastructure/)
- **仓储接口**: `internal/infrastructure/persistence/repository/{module}.go`
  - 接口定义和实现放在同一个文件
  - 接口名: `{Module}Repository`
  - 实现类型: `{module}Repository` (小写，私有)
  - 构造函数: `New{Module}Repository()`
- **数据库模型**: 与仓储实现在同一文件
  - 模型名: `{Module}Model`
  - 包含 `ToEntity()` 和 `FromEntity()` 方法

#### 接口层 (interfaces/)
- **处理器**: `internal/interfaces/http/handler/{module}.go`
  - 结构体: `{Module}Handler`
  - 构造函数: `New{Module}Handler()`
- **中间件**: `internal/interfaces/http/middleware/{module}.go`
  - 结构体: `{Module}Middleware`
  - 构造函数: `New{Module}Middleware()`
- **路由**: `internal/interfaces/http/router/{module}.go`
  - 结构体: `{Module}Router`
  - 构造函数: `New{Module}Router()`
  - 实现 `RegisterRoutes(engine *gin.Engine)` 方法
- **路由设置**: 在 `internal/interfaces/http/setup.go` 中注册新路由

### 2. 命名规范

- **包名**: 使用小写，单数形式（如 `auth`, `user`, `order`）
- **结构体**: 使用大驼峰（如 `User`, `AuthService`, `UserHandler`）
- **接口**: 使用大驼峰 + `Repository` 后缀（如 `UserRepository`）
- **函数**: 
  - 构造函数: `New{Type}()`
  - 方法: 使用动词开头（如 `Get`, `Create`, `Update`, `Delete`）

### 3. 依赖注入规范

- **最小化引入**: 只引入必要的依赖
- **依赖方向**: 
  - domain 层不依赖 infrastructure 和 interfaces
  - domain 层依赖 repository 接口（通过 import）
  - infrastructure 层实现 repository 接口
  - interfaces 层依赖 domain 层服务

### 4. 代码组织原则

#### 领域层代码结构
```go
package {module}

import (
    // 最小化引入，只引入必要的包
)

// 实体定义
type Entity struct {
    // ...
}

// 请求/响应结构
type Request struct {
    // ...
}

// 服务结构
type Service struct {
    repo RepositoryInterface
    // 其他依赖
}

// NewService 构造函数
func NewService(...) *Service {
    // ...
}

// 业务方法
func (s *Service) Method() {
    // ...
}
```

#### 仓储代码结构
```go
package repository

// 接口定义
type {Module}Repository interface {
    // 方法定义
}

// 模型定义
type {Module}Model struct {
    // GORM 标签
}

// 实现结构（私有）
type {module}Repository struct {
    db    *mysql.DB
    redis *redis.Client
}

// 构造函数
func New{Module}Repository(...) {Module}Repository {
    return &{module}Repository{...}
}

// 实现方法
func (r *{module}Repository) Method() {
    // ...
}
```

#### 接口层代码结构
```go
// handler/{module}.go
type {Module}Handler struct {
    service *domain.{Module}.Service
}

func New{Module}Handler(...) *{Module}Handler {
    // ...
}

func (h *{Module}Handler) Handle(c *gin.Context) {
    // ...
}

// router/{module}.go
type {Module}Router struct {
    handler   *handler.{Module}Handler
    middleware *middleware.{Module}Middleware
}

func (r *{Module}Router) RegisterRoutes(engine *gin.Engine) {
    // 注册路由
}
```

### 5. 路由注册规范

在 `internal/interfaces/http/setup.go` 中注册新模块：

```go
func SetupRoutes(engine *gin.Engine, cfg *config.Config, db *mysql.DB, rdb *redis.Client) {
    // 创建仓储
    {module}Repo := repository.New{Module}Repository(db, rdb)
    
    // 创建服务
    {module}Service := {module}.NewService({module}Repo, ...)
    
    // 创建处理器和中间件
    {module}Handler := handler.New{Module}Handler({module}Service)
    {module}Middleware := middleware.New{Module}Middleware({module}Service)
    
    // 注册路由
    {module}Router := router.New{Module}Router({module}Handler, {module}Middleware)
    {module}Router.RegisterRoutes(engine)
}
```

### 6. 错误处理规范

- 使用 `pkg/errors` 包中的错误类型
- 领域层返回业务错误
- 接口层将错误转换为 HTTP 响应
- 使用 `pkg/utils` 中的响应工具函数

### 7. 禁止事项

- **禁止**在 domain 层创建多个文件
- **禁止**在 domain 层直接依赖 infrastructure 具体实现
- **禁止**在接口层直接操作数据库
- **禁止**创建不必要的抽象层
- **禁止**使用复杂的依赖注入框架（已使用 Wire）

### 8. 代码生成检查清单

生成代码前请确认：
- [ ] 目录结构符合规范
- [ ] 命名符合规范
- [ ] 依赖关系正确（domain → repository 接口，infrastructure 实现接口）
- [ ] 在 setup.go 中注册了路由
- [ ] 错误处理使用了项目标准
- [ ] 代码简洁，没有过度设计

---

**请所有 AI 助手在开始工作前仔细阅读并遵守本规范！**

