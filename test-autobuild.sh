#!/bin/bash

echo "=== 前端自动构建功能测试 ==="

echo "1. 测试关闭自动构建..."
cd backend
export FRONTEND_AUTO_BUILD=false
echo "✓ 设置 FRONTEND_AUTO_BUILD=false"

echo "2. 测试开启自动构建..."
export FRONTEND_AUTO_BUILD=true
echo "✓ 设置 FRONTEND_AUTO_BUILD=true"

echo "3. 删除构建文件来测试自动构建..."
rm -rf ../frontend/dist
echo "✓ 删除 frontend/dist 目录"

echo "4. 启动服务器测试自动构建..."
echo "注意：服务器会自动检测前端文件变化并重新构建"
echo "如果构建成功，应该会看到 'Frontend build completed successfully' 消息"

echo ""
echo "配置说明："
echo "- FRONTEND_AUTO_BUILD=true/false  : 启用/禁用自动构建"
echo "- FRONTEND_BUILD_COMMAND         : 构建命令 (默认: npm run build)"
echo "- FRONTEND_DIST_PATH             : 构建输出路径 (默认: ../frontend/dist)"
echo "- FRONTEND_SOURCE_PATH           : 源码路径 (默认: ../frontend)"
echo ""
echo "手动重建API："
echo "- POST /api/system/rebuild-frontend (需要管理员权限)"
echo "- GET  /api/system/build-status"