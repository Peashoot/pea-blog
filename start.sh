#!/bin/bash

echo "=== Pea Blog 项目启动脚本 ==="

# 检查并安装后端依赖
echo "1. 检查后端依赖..."
cd backend
if ! go mod download; then
    echo "❌ 后端依赖安装失败"
    exit 1
fi
echo "✓ 后端依赖安装成功"

# 启动后端服务（后台运行）
echo "2. 启动后端服务..."
go run cmd/server/main.go &
BACKEND_PID=$!
echo "✓ 后端服务已启动 (PID: $BACKEND_PID)"

# 等待后端启动
sleep 3

# 检查后端是否正常运行
if ! curl -s http://localhost:8080/api/articles >/dev/null; then
    echo "❌ 后端服务启动失败"
    kill $BACKEND_PID 2>/dev/null
    exit 1
fi
echo "✓ 后端服务运行正常"

# 检查并安装前端依赖
echo "3. 检查前端依赖..."
cd ../frontend
if [ ! -d "node_modules" ]; then
    echo "安装前端依赖..."
    if ! npm install; then
        echo "❌ 前端依赖安装失败"
        kill $BACKEND_PID 2>/dev/null
        exit 1
    fi
fi
echo "✓ 前端依赖就绪"

# 启动前端服务
echo "4. 启动前端服务..."
echo "✓ 项目启动成功！"
echo ""
echo "访问地址："
echo "- 前端: http://localhost:5173"
echo "- 后端: http://localhost:8080"
echo ""
echo "默认管理员账户："
echo "- 用户名: admin"
echo "- 密码: password"
echo ""
echo "按 Ctrl+C 停止所有服务"

# 启动前端（前台运行）
npm run dev

# 清理后台进程
echo "停止后端服务..."
kill $BACKEND_PID 2>/dev/null
echo "✓ 所有服务已停止"