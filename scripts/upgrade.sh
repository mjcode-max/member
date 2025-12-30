#!/bin/bash

# 会员预约系统 - 升级脚本
# 使用方法: ./upgrade.sh [安装目录]

set -e

# 默认安装目录
INSTALL_DIR="${1:-/opt/member-pre}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="$(dirname "$SCRIPT_DIR")/build"
BACKUP_DIR="$INSTALL_DIR/backup_$(date +%Y%m%d_%H%M%S)"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印信息
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查是否为 root 用户
check_root() {
    if [ "$EUID" -ne 0 ]; then
        error "请使用 root 用户运行此脚本"
        exit 1
    fi
}

# 检查安装目录是否存在
check_install_dir() {
    if [ ! -d "$INSTALL_DIR" ]; then
        error "安装目录不存在: $INSTALL_DIR"
        error "请先执行安装脚本: ./install.sh"
        exit 1
    fi
}

# 检查安装包是否存在
check_build_dir() {
    if [ ! -d "$BUILD_DIR" ]; then
        error "未找到 build 目录: $BUILD_DIR"
        error "请先执行 make package 生成安装包"
        exit 1
    fi
}

# 停止服务
stop_service() {
    info "停止服务..."
    if systemctl is-active --quiet member-pre 2>/dev/null; then
        systemctl stop member-pre
        info "服务已停止"
    else
        warn "服务未运行或未使用 systemd 管理"
        # 尝试查找并停止进程
        PID=$(pgrep -f "bin/server server" || true)
        if [ -n "$PID" ]; then
            warn "发现运行中的进程 PID: $PID"
            read -p "是否停止该进程? (y/n): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                kill "$PID"
                sleep 2
                info "进程已停止"
            else
                error "用户取消操作"
                exit 1
            fi
        fi
    fi
}

# 备份当前版本
backup_current() {
    info "备份当前版本到: $BACKUP_DIR"
    mkdir -p "$BACKUP_DIR"
    
    # 备份二进制文件
    if [ -d "$INSTALL_DIR/bin" ]; then
        cp -r "$INSTALL_DIR/bin" "$BACKUP_DIR/"
    fi
    
    # 备份配置文件
    if [ -d "$INSTALL_DIR/configs" ]; then
        cp -r "$INSTALL_DIR/configs" "$BACKUP_DIR/"
    fi
    
    # 备份前端文件
    if [ -d "$INSTALL_DIR/dist" ]; then
        cp -r "$INSTALL_DIR/dist" "$BACKUP_DIR/"
    fi
    
    # 备份数据库（如果可能）
    if command -v mysqldump &> /dev/null; then
        info "备份数据库..."
        # 从配置文件读取数据库信息（这里简化处理）
        DB_NAME="member_pre"
        if [ -f "$INSTALL_DIR/configs/config.yaml" ]; then
            # 尝试从配置文件提取数据库名
            DB_NAME=$(grep -A 10 "^database:" "$INSTALL_DIR/configs/config.yaml" | grep "dbname:" | awk '{print $2}' | tr -d '"' || echo "member_pre")
        fi
        mysqldump -u root -p"$MYSQL_ROOT_PASSWORD" "$DB_NAME" > "$BACKUP_DIR/database_$(date +%Y%m%d_%H%M%S).sql" 2>/dev/null || {
            warn "数据库备份失败，请手动备份"
        }
    fi
    
    info "备份完成: $BACKUP_DIR"
}

# 替换文件
replace_files() {
    info "替换文件..."
    
    # 删除旧文件（保留配置和备份）
    rm -rf "$INSTALL_DIR/bin"
    rm -rf "$INSTALL_DIR/dist"
    rm -rf "$INSTALL_DIR/docs"
    rm -rf "$INSTALL_DIR/scripts"
    
    # 复制新文件
    cp -r "$BUILD_DIR/bin" "$INSTALL_DIR/"
    cp -r "$BUILD_DIR/dist" "$INSTALL_DIR/"
    cp -r "$BUILD_DIR/docs" "$INSTALL_DIR/" 2>/dev/null || true
    cp -r "$BUILD_DIR/scripts" "$INSTALL_DIR/" 2>/dev/null || true
    
    # 设置可执行权限
    chmod +x "$INSTALL_DIR/bin/server"
    
    info "文件替换完成"
}

# 执行数据库迁移
run_migration() {
    info "执行数据库迁移..."
    read -p "是否现在执行数据库迁移? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        "$INSTALL_DIR/bin/server" migrate up --config "$INSTALL_DIR/configs/config.yaml" || {
            error "数据库迁移失败！"
            warn "可以使用备份恢复: $BACKUP_DIR"
            exit 1
        }
        info "数据库迁移完成"
    else
        warn "跳过数据库迁移，请稍后手动执行: $INSTALL_DIR/bin/server migrate up"
    fi
}

# 启动服务
start_service() {
    info "启动服务..."
    if systemctl list-unit-files | grep -q "member-pre.service"; then
        systemctl start member-pre
        sleep 2
        if systemctl is-active --quiet member-pre; then
            info "服务启动成功"
            systemctl status member-pre --no-pager
        else
            error "服务启动失败，请检查日志: journalctl -u member-pre -n 50"
            exit 1
        fi
    else
        warn "未找到 systemd 服务，请手动启动: $INSTALL_DIR/bin/server server"
    fi
}

# 主函数
main() {
    info "=========================================="
    info "会员预约系统 - 升级脚本"
    info "=========================================="
    info "安装目录: $INSTALL_DIR"
    echo
    
    check_root
    check_install_dir
    check_build_dir
    
    # 确认升级
    warn "即将升级系统，当前版本将被备份"
    read -p "确认继续? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        info "用户取消操作"
        exit 0
    fi
    
    stop_service
    backup_current
    replace_files
    run_migration
    start_service
    
    echo
    info "=========================================="
    info "升级完成！"
    info "安装目录: $INSTALL_DIR"
    info "备份目录: $BACKUP_DIR"
    info "=========================================="
    echo
    info "如果遇到问题，可以从备份恢复:"
    info "  cp -r $BACKUP_DIR/* $INSTALL_DIR/"
}

# 执行主函数
main

