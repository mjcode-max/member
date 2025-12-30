.PHONY: wire gen run build migrate-up migrate-down migrate-status package clean

# 生成 Wire 代码
wire:
	@which wire > /dev/null || (echo "Wire 未安装，正在安装..." && go install github.com/google/wire/cmd/wire@latest)
	@export PATH=$$PATH:$$(go env GOPATH)/bin && cd internal/infrastructure && wire

# 运行应用
run:
	go run cmd/main.go server

# 构建应用
build: wire
	go build -o bin/server cmd/main.go

# 数据库迁移 - 执行迁移
migrate-up:
	go run cmd/main.go migrate up

# 数据库迁移 - 回滚迁移
migrate-down:
	go run cmd/main.go migrate down

# 数据库迁移 - 查看状态
migrate-status:
	go run cmd/main.go migrate status

# 下载依赖
deps:
	go mod download
	go mod tidy

# 初始化 Wire（首次使用）
wire-init:
	go install github.com/google/wire/cmd/wire@latest

# 打包 - 构建前端和后端，生成 build 文件夹并压缩
package: clean
	@echo "开始打包..."
	@mkdir -p build/bin build/configs build/dist
	@echo "构建后端二进制..."
	@go build -ldflags="-s -w" -o build/bin/server cmd/main.go
	@echo "复制配置文件..."
	@cp -r configs/* build/configs/
	@echo "复制文档和脚本..."
	@test -d docs && cp -r docs build/ || true
	@test -d scripts && cp -r scripts build/ || true
	@echo "构建前端项目..."
	@echo "  构建 admin-web..."
	@cd web/admin-web && npm run build && cp -r dist ../../build/dist/admin-web || (echo "  admin-web 构建失败，跳过" && true)
	@echo "  构建 customer-h5..."
	@cd web/customer-h5 && npm run build && cp -r dist ../../build/dist/customer-h5 || (echo "  customer-h5 构建失败，跳过" && true)
	@echo "  构建 staff-h5..."
	@cd web/staff-h5 && npm run build && cp -r dist ../../build/dist/staff-h5 || (echo "  staff-h5 构建失败，跳过" && true)
	@echo "  构建 store-h5..."
	@cd web/store-h5 && npm run build && cp -r dist ../../build/dist/store-h5 || (echo "  store-h5 构建失败，跳过" && true)
	@echo "压缩打包文件..."
	@tar -czf build.tar.gz -C build .
	@echo "打包完成！"
	@echo "  输出目录: build/"
	@echo "  压缩包: build.tar.gz"

# 清理构建产物
clean:
	@echo "清理构建产物..."
	@rm -rf build bin build.tar.gz
	@echo "清理完成"

