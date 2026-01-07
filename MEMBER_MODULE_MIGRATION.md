# 会员管理模块 - 数据库迁移和测试指南

## 前置步骤

在运行数据库迁移之前，需要先重新生成 Wire 依赖注入代码，因为我们添加了新的会员管理模块。

## 步骤 1: 重新生成 Wire 代码

由于添加了新的依赖（MemberService、UsageService、MemberRepository、UsageRepository、MemberHandler），需要重新生成 Wire 代码：

```bash
# 方式1: 使用 Makefile
make wire

# 方式2: 使用 Taskfile
task wire

# 方式3: 手动执行
cd internal/infrastructure && wire
```

## 步骤 2: 执行数据库迁移

生成 Wire 代码后，执行数据库迁移以创建会员表和使用记录表：

```bash
# 方式1: 使用 Makefile
make migrate-up

# 方式2: 使用 Taskfile
task migrate-up

# 方式3: 手动执行
go run ./cmd migrate up
```

迁移将创建以下表：
- `members` - 会员表（包含套餐信息）
- `member_usages` - 会员使用记录表

## 步骤 3: 验证迁移状态

检查迁移是否成功：

```bash
# 方式1: 使用 Makefile
make migrate-status

# 方式2: 使用 Taskfile
task migrate-status

# 方式3: 手动执行
go run ./cmd migrate status
```

## 步骤 4: 运行单元测试

运行领域层的单元测试：

```bash
# 方式1: 使用 Makefile
make test

# 方式2: 使用 Taskfile
task test

# 方式3: 手动执行
go test ./internal/domain/... -v
```

## 步骤 5: 启动服务并测试 API

1. 启动服务：
```bash
# 方式1: 使用 Makefile
make run

# 方式2: 使用 Taskfile
task run

# 方式3: 手动执行
go run ./cmd server
```

2. 测试会员管理 API（需要先登录获取 token）：

### 创建会员（管理员/店长）
```bash
POST /api/v1/members
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "张三",
  "phone": "13800138000",
  "package_name": "美甲年卡",
  "service_type": "nail",
  "price": 1999.00,
  "validity_duration": 365,
  "store_id": 1,
  "purchase_amount": 1999.00,
  "description": "美甲年卡套餐"
}
```

### 获取会员列表
```bash
GET /api/v1/members?page=1&page_size=10
Authorization: Bearer {token}
```

### 创建使用记录（管理员/店长/美甲师）
```bash
POST /api/v1/members/{member_id}/usages
Authorization: Bearer {token}
Content-Type: application/json

{
  "service_item": "美甲-单色",
  "store_id": 1,
  "technician_id": 2,
  "usage_date": "2024-01-15",
  "remark": "客户满意"
}
```

### 获取使用记录列表
```bash
GET /api/v1/usages?member_id=1
Authorization: Bearer {token}
```

### 删除使用记录（仅管理员）
```bash
DELETE /api/v1/usages/{usage_id}
Authorization: Bearer {token}
```

## 数据库表结构验证

迁移成功后，可以连接数据库验证表结构：

### members 表结构
- id (uint, primary key)
- name (string, not null)
- phone (string, index)
- package_name (string)
- service_type (string)
- price (decimal(10,2))
- used_times (int, default 0)
- validity_duration (int)
- valid_from (date)
- valid_to (date)
- store_id (uint, index)
- purchase_amount (decimal(10,2))
- purchase_time (timestamp)
- status (string, default 'active')
- description (text)
- created_by (uint, index)
- created_at (timestamp)
- updated_at (timestamp)
- deleted_at (timestamp, soft delete)

### member_usages 表结构
- id (uint, primary key)
- member_id (uint, index)
- package_name (string)
- service_item (string)
- store_id (uint, index)
- store_name (string)
- technician_id (uint, index, nullable)
- technician_name (string)
- usage_date (date)
- remark (text)
- created_by (uint, index)
- created_at (timestamp)
- updated_at (timestamp)

## 常见问题

### 1. Wire 代码生成失败
- 确保已安装 wire: `go install github.com/google/wire/cmd/wire@latest`
- 检查 wireSet.go 文件中的依赖绑定是否正确

### 2. 迁移失败
- 检查数据库连接配置
- 确保数据库用户有创建表的权限
- 查看日志了解具体错误信息

### 3. 编译错误
- 运行 `go mod tidy` 确保依赖正确
- 检查所有导入的包是否存在
- 确保 wire_gen.go 已更新

## 注意事项

1. **事务处理**: 创建使用记录时会自动递增会员的 `used_times`，这两个操作应该在同一个事务中完成。当前实现中，如果创建使用记录成功但递增次数失败，需要手动处理回滚。

2. **权限验证**: 
   - 店长只能操作自己门店的会员
   - 美甲师只能记录自己门店的使用
   - 管理员可以操作所有数据

3. **有效期处理**: 
   - 支持固定时长（validity_duration）和手动指定日期（valid_from, valid_to）两种方式
   - 两种方式可以互相转换
   - 系统优先使用 valid_from 和 valid_to 作为真实值

4. **无限次使用**: 会员在有效期内可以无限次使用，`used_times` 仅用于统计显示，不影响使用权限。

