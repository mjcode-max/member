#!/bin/bash

# 会员预约系统 - 安装脚本
# 使用方法: ./install.sh [安装目录]

set -e

# 默认安装目录
INSTALL_DIR="${1:-/opt/member-pre}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="$(dirname "$SCRIPT_DIR")/build"

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

# 检查安装包是否存在
check_build_dir() {
    if [ ! -d "$BUILD_DIR" ]; then
        error "未找到 build 目录: $BUILD_DIR"
        error "请先执行 make package 生成安装包"
        exit 1
    fi
}

# 创建安装目录
create_install_dir() {
    info "创建安装目录: $INSTALL_DIR"
    mkdir -p "$INSTALL_DIR"
}

# 复制文件
copy_files() {
    info "复制文件到安装目录..."
    cp -r "$BUILD_DIR"/* "$INSTALL_DIR/"
    
    # 设置可执行权限
    chmod +x "$INSTALL_DIR/bin/server"
    
    info "文件复制完成"
}

# 创建配置文件（如果不存在）
create_config() {
    if [ ! -f "$INSTALL_DIR/configs/config.yaml" ]; then
        warn "配置文件不存在，将创建默认配置"
        # 这里可以创建默认配置或提示用户
    fi
}

# 执行数据库迁移
run_migration() {
    info "执行数据库迁移..."
    read -p "是否现在执行数据库迁移? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        "$INSTALL_DIR/bin/server" migrate up --config "$INSTALL_DIR/configs/config.yaml" || {
            warn "数据库迁移失败，请检查配置文件后手动执行: $INSTALL_DIR/bin/server migrate up"
        }
    else
        warn "跳过数据库迁移，请稍后手动执行: $INSTALL_DIR/bin/server migrate up"
    fi
}

# 创建 systemd 服务文件
create_systemd_service() {
    info "创建 systemd 服务文件..."
    
    SERVICE_FILE="/etc/systemd/system/member-pre.service"
    
    cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=Member Pre Booking System
After=network.target mysql.service redis.service

[Service]
Type=simple
User=root
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/bin/server server --config $INSTALL_DIR/configs/config.yaml
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

    info "systemd 服务文件已创建: $SERVICE_FILE"
    info "使用以下命令管理服务:"
    info "  启动: systemctl start member-pre"
    info "  停止: systemctl stop member-pre"
    info "  重启: systemctl restart member-pre"
    info "  状态: systemctl status member-pre"
    info "  开机自启: systemctl enable member-pre"
}

# 主函数
main() {
    info "=========================================="
    info "会员预约系统 - 安装脚本"
    info "=========================================="
    info "安装目录: $INSTALL_DIR"
    echo
    
    check_root
    check_build_dir
    create_install_dir
    copy_files
    create_config
    
    echo
    info "安装完成！"
    echo
    info "下一步操作:"
    info "1. 编辑配置文件: $INSTALL_DIR/configs/config.yaml"
    info "2. 执行数据库迁移: $INSTALL_DIR/bin/server migrate up"
    info "3. 创建 systemd 服务（可选）"
    
    echo
    read -p "是否创建 systemd 服务? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        create_systemd_service
        systemctl daemon-reload
        info "systemd 服务已创建并重载"
    fi
    
    echo
    run_migration
    
    echo
    info "=========================================="
    info "安装完成！"
    info "安装目录: $INSTALL_DIR"
    info "=========================================="
}

# 执行主函数
main

