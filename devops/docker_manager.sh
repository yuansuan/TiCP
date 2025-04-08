#!/bin/bash

# 定义所有的 Docker Compose 文件
COMPOSE_FILES=(
    "docker-compose-base.yml"
    "docker-compose-cloud-base.yml"
    "docker-compose-ipaas.yml"
    "docker-compose-psp.yml"
    "docker-compose-ipaas-hpc.yml"
)

# 检查是否提供了参数
if [ $# -eq 0 ]; then
    echo "使用方法: $0 {up|down|status|restart}"
    exit 1
fi

# 获取用户输入的操作（up/down/status/restart）
ACTION=$1

# 处理不同的操作
case "$ACTION" in
    up)
        echo "启动所有 Docker Compose 服务..."
        for FILE in "${COMPOSE_FILES[@]}"; do
            echo "正在启动: $FILE"
            docker compose -f "$FILE" up -d
        done
        echo "所有服务已启动 ✅"
        ;;
    
    down)
        echo "停止所有 Docker Compose 服务..."
        for FILE in "${COMPOSE_FILES[@]}"; do
            echo "正在停止: $FILE"
            docker compose -f "$FILE" down
        done
        echo "所有服务已停止 ❌"
        ;;
    
    status)
        echo "查询所有 Docker Compose 服务状态..."
        docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
        ;;
    
    restart)
        echo "重启所有 Docker Compose 服务..."
        for FILE in "${COMPOSE_FILES[@]}"; do
            echo "正在重启: $FILE"
            docker compose -f "$FILE" down
            docker compose -f "$FILE" up -d
        done
        echo "所有服务已重启 🔄"
        ;;
    
    *)
        echo "无效命令！使用: $0 {up|down|status|restart}"
        exit 1
        ;;
esac

