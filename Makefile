.PHONY: wire gen run build migrate-up migrate-down migrate-status package clean test link frontend-dev frontend-stop

# 生成 Wire 代码
wire:
	@which wire > /dev/null || (echo "Wire 未安装，正在安装..." && go install github.com/google/wire/cmd/wire@latest)
	@export PATH=$$PATH:$$(go env GOPATH)/bin && cd internal/infrastructure && wire

# 运行应用
run:
	go run ./cmd server

# 构建应用
build: wire swagger
	go build -o bin/server ./cmd

# 数据库迁移 - 执行迁移
migrate-up:
	go run ./cmd migrate up

# 数据库迁移 - 回滚迁移
migrate-down:
	go run ./cmd migrate down

# 数据库迁移 - 查看状态
migrate-status:
	go run ./cmd migrate status

# 下载依赖
deps:
	go mod download
	go mod tidy

# 初始化 Wire（首次使用）
wire-init:
	go install github.com/google/wire/cmd/wire@latest

# 生成 Swagger 文档
swagger:
	@which swag > /dev/null || (echo "swag 未安装，正在安装..." && go install github.com/swaggo/swag/cmd/swag@latest)
	@export PATH=$$PATH:$$(go env GOPATH)/bin && swag init -g cmd/main.go -o ./docs --parseDependency --parseInternal
	@echo "Swagger 文档生成完成！访问 http://localhost:8080/swagger/index.html 查看文档"

# 运行 domain 包的单元测试
test:
	@echo "运行 domain 包的单元测试..."
	@go test ./internal/domain/... -v
	@echo "测试完成！"

# 运行 Go 静态分析工具
link:
	@echo "运行 Go 静态分析..."
	@go vet ./...
	@echo "静态分析完成！"

# 打包 - 构建前端和后端，生成 build 文件夹并压缩
package: clean deps wire swagger link test
	@echo "开始打包..."
	@mkdir -p build/bin build/configs build/dist
	@echo "构建后端二进制..."
	@go build -ldflags="-s -w" -o build/bin/server ./cmd
	@echo "复制配置文件..."
	@cp -r configs/* build/configs/
	@echo "复制文档和脚本..."
	@test -d docs && cp -r docs build/ || true
	@test -d scripts && cp -r scripts build/ || true
	@echo "构建前端项目..."
	@echo "  构建 admin-web..."
	@cd web/admin-web && npm install && npm run build && cp -r dist ../../build/dist/admin-web || (echo "  admin-web 构建失败，跳过" && true)
	@echo "  构建 customer-h5..."
	@cd web/customer-h5 && npm install && npm run build && cp -r dist ../../build/dist/customer-h5 || (echo "  customer-h5 构建失败，跳过" && true)
	@echo "  构建 staff-h5..."
	@cd web/staff-h5 && npm install && npm run build && cp -r dist ../../build/dist/staff-h5 || (echo "  staff-h5 构建失败，跳过" && true)
	@echo "  构建 store-h5..."
	@cd web/store-h5 && npm install && npm run build && cp -r dist ../../build/dist/store-h5 || (echo "  store-h5 构建失败，跳过" && true)
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

pre-commit: clean deps wire swagger link test

# 启动所有前端项目（开发模式）
frontend-dev:
	@echo "启动所有前端项目..."
	@mkdir -p logs
	@echo "  启动 admin-web (后台管理)..."
	@cd web/admin-web && npm run dev > ../../logs/admin-web.log 2>&1 &
	@echo "  启动 customer-h5 (客户H5)..."
	@cd web/customer-h5 && npm run dev > ../../logs/customer-h5.log 2>&1 &
	@echo "  启动 staff-h5 (员工H5)..."
	@cd web/staff-h5 && npm run dev > ../../logs/staff-h5.log 2>&1 &
	@echo "  启动 store-h5 (门店H5)..."
	@cd web/store-h5 && npm run dev > ../../logs/store-h5.log 2>&1 &
	@sleep 2
	@echo ""
	@echo "所有前端项目已启动！"
	@echo "  后台管理 (admin-web): http://localhost:3000 (查看日志: tail -f logs/admin-web.log)"
	@echo "  客户H5 (customer-h5): http://localhost:3001 (查看日志: tail -f logs/customer-h5.log)"
	@echo "  员工H5 (staff-h5): http://localhost:3002 (查看日志: tail -f logs/staff-h5.log)"
	@echo "  门店H5 (store-h5): http://localhost:3003 (查看日志: tail -f logs/store-h5.log)"
	@echo ""
	@echo "使用 'make frontend-stop' 可以停止所有前端服务"
	@echo "使用 'tail -f logs/*.log' 可以查看所有日志"

# 停止所有前端项目
frontend-stop:
	@echo "停止所有前端项目..."
	@pkill -f "web/admin-web.*vite" || true
	@pkill -f "web/customer-h5.*vite" || true
	@pkill -f "web/staff-h5.*vite" || true
	@pkill -f "web/store-h5.*vite" || true
	@echo "所有前端项目已停止"