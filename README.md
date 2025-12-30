# 会员预约系统

基于 Go 语言开发的会员预约系统，采用领域驱动设计（DDD）架构。

## 技术栈

- **Gin**: HTTP Web 框架
- **GORM**: ORM 框架
- **Viper**: 配置管理
- **Wire**: 依赖注入
- **Cobra**: 命令行工具
- **MySQL**: 关系型数据库
- **Redis**: 缓存数据库
- **Zap**: 日志库

## 项目结构

```
member-pre/
├── cmd/server/          # 应用入口
├── internal/            # 内部代码
│   ├── domain/          # 领域模块（自行组织）
│   ├── infrastructure/  # 基础设施层
│   │   ├── config/      # 配置管理
│   │   ├── logger/      # 日志
│   │   ├── persistence/ # 持久化
│   │   │   ├── mysql/   # MySQL
│   │   │   └── redis/   # Redis
│   │   ├── http/        # HTTP 服务器
│   │   └── wire.go      # Wire 依赖注入
│   └── interfaces/      # 接口层
│       └── http/        # HTTP 接口
├── configs/             # 配置文件
├── pkg/                 # 公共包
│   ├── errors/          # 错误处理
│   ├── utils/           # 工具函数
│   └── validator/       # 验证器
└── web/                 # 前端代码
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
go mod tidy
```

### 2. 安装 Wire

```bash
go install github.com/google/wire/cmd/wire@latest
```

### 3. 生成 Wire 代码

```bash
make wire
# 或者
cd internal/infrastructure && wire
```

### 4. 配置数据库

编辑 `configs/config.yaml`，配置 MySQL 和 Redis 连接信息。

### 5. 运行应用

```bash
make run
# 或者
go run cmd/server/main.go server
```

### 6. 数据库迁移

```bash
# 执行迁移
make migrate-up
# 或
go run cmd/server/main.go migrate up

# 回滚迁移
make migrate-down
# 或
go run cmd/server/main.go migrate down

# 查看迁移状态
make migrate-status
# 或
go run cmd/server/main.go migrate status
```

### 7. 构建应用

```bash
make build
```

## 配置说明

配置文件位于 `configs/config.yaml`，包含以下配置：

- **server**: 服务器配置（端口、模式等）
- **database**: MySQL 数据库配置
- **redis**: Redis 配置
- **log**: 日志配置

## 开发指南

### 添加新的路由

1. 在 `internal/interfaces/http/` 下创建路由文件
2. 实现 `Router` 接口
3. 在 `router.go` 中注册路由

### 添加领域模块

在 `internal/domain/` 目录下自行组织领域模块代码。

### 使用数据库

通过 Wire 注入的 `*mysql.DB` 实例使用 GORM 操作数据库。

### 使用 Redis

通过 Wire 注入的 `*redis.Client` 实例操作 Redis。

### 使用日志

通过 Wire 注入的 `*logger.ZapLogger` 实例记录日志。

## 命令行工具 (Cobra)

项目使用 Cobra 管理命令行，支持以下命令：

- `server`: 启动HTTP服务器
  ```bash
  go run cmd/server/main.go server
  ```

- `migrate`: 数据库迁移
  - `migrate up`: 执行所有待执行的迁移
  - `migrate down`: 回滚最后一次迁移
  - `migrate status`: 查看迁移状态

所有命令都支持 `--config` 或 `-c` 参数指定配置文件路径。

## Makefile 命令

- `make wire`: 生成 Wire 代码
- `make run`: 运行应用（server命令）
- `make build`: 构建应用
- `make migrate-up`: 执行数据库迁移
- `make migrate-down`: 回滚数据库迁移
- `make migrate-status`: 查看迁移状态
- `make deps`: 下载并整理依赖
- `make wire-init`: 安装 Wire 工具

## 工具包说明

### pkg/errors
统一的错误处理包，提供错误码和错误类型定义。

### pkg/utils
工具函数包，包含：
- `string.go`: 字符串处理工具
- `time.go`: 时间处理工具
- `response.go`: HTTP响应工具

### pkg/validator
数据验证包，提供常用的验证方法。

