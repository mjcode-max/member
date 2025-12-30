# 安装和更新说明

## 安装说明

### 1. 解压安装包

```bash
tar -xzf build.tar.gz
cd build
```

### 2. 配置数据库

编辑 `configs/config.yaml` 文件，配置数据库和 Redis 连接信息：

```yaml
database:
  host: localhost
  port: 3306
  user: root
  password: your_password
  dbname: member_pre

redis:
  host: localhost
  port: 6379
  password: your_redis_password
  db: 0
```

### 3. 执行数据库迁移

```bash
./bin/server migrate up
```

### 4. 启动服务

```bash
./bin/server server
```

或者使用配置文件：

```bash
./bin/server server --config configs/config.yaml
```

### 5. 验证安装

访问健康检查接口：

```bash
curl http://localhost:8080/api/v1/health
```

## 更新说明

### 1. 停止服务

```bash
# 使用 Ctrl+C 或 kill 命令停止当前运行的服务
kill <pid>
```

### 2. 备份数据

```bash
# 备份数据库
mysqldump -u root -p member_pre > backup_$(date +%Y%m%d).sql

# 备份配置文件
cp configs/config.yaml configs/config.yaml.bak
```

### 3. 解压新版本

```bash
# 解压新版本到新目录
tar -xzf build.tar.gz -C /tmp/new_build

# 复制配置文件到新版本
cp configs/config.yaml /tmp/new_build/configs/
```

### 4. 执行数据库迁移

```bash
cd /tmp/new_build
./bin/server migrate up
```

### 5. 替换旧版本

```bash
# 停止服务后，替换文件
cd /path/to/old/build
rm -rf bin dist
cp -r /tmp/new_build/bin .
cp -r /tmp/new_build/dist .
cp -r /tmp/new_build/docs .
cp -r /tmp/new_build/scripts .
```

### 6. 启动新服务

```bash
./bin/server server --config configs/config.yaml
```

### 7. 验证更新

```bash
curl http://localhost:8080/api/v1/health
```

## 常见问题

### 数据库连接失败

- 检查数据库服务是否启动
- 检查配置文件中的数据库连接信息是否正确
- 检查数据库用户权限

### Redis 连接失败

- 检查 Redis 服务是否启动
- 检查配置文件中的 Redis 连接信息是否正确
- 检查防火墙设置

### 端口被占用

修改 `configs/config.yaml` 中的端口配置：

```yaml
server:
  port: 8080  # 修改为其他可用端口
```

### 权限问题

确保二进制文件有执行权限：

```bash
chmod +x bin/server
```

## 目录说明

- `bin/` - 可执行文件
- `configs/` - 配置文件
- `dist/` - 前端静态资源
  - `admin-web/` - 管理后台
  - `customer-h5/` - 客户H5
  - `staff-h5/` - 员工H5
  - `store-h5/` - 门店H5
- `docs/` - 文档
- `scripts/` - 脚本文件

