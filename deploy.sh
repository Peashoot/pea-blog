#!/bin/bash

echo "=== Pea Blog Docker 部署脚本 ==="

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装Docker Compose"
    exit 1
fi

echo "✓ Docker 环境检查通过"

# 构建并启动服务
echo "1. 构建Docker镜像..."
docker-compose build

if [ $? -eq 0 ]; then
    echo "✓ Docker镜像构建成功"
else
    echo "❌ Docker镜像构建失败"
    exit 1
fi

echo "2. 启动服务..."
docker-compose up -d

if [ $? -eq 0 ]; then
    echo "✓ 服务启动成功"
    echo ""
    echo "访问地址："
    echo "- 网站: http://localhost:8080"
    echo "- API: http://localhost:8080/api"
    echo ""
    echo "默认管理员账户："
    echo "- 用户名: admin"
    echo "- 密码: password"
    echo ""
    echo "查看日志: docker-compose logs -f"
    echo "停止服务: docker-compose down"
else
    echo "❌ 服务启动失败"
    exit 1
fi