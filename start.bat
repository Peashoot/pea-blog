@echo off
chcp 65001
echo === Pea Blog 项目启动脚本 ===

echo 1. 检查后端依赖...
cd backend
go mod download
if %errorlevel% neq 0 (
    echo ❌ 后端依赖安装失败
    pause
    exit /b 1
)
echo ✓ 后端依赖安装成功

echo 2. 启动后端服务...
start "Pea Blog Backend" /B go run cmd/server/main.go
if %errorlevel% neq 0 (
    echo ❌ 后端服务启动失败
    pause
    exit /b 1
)
echo ✓ 后端服务已启动

timeout /t 3 /nobreak > nul

echo 3. 检查前端依赖...
cd ..\frontend
if not exist "node_modules" (
    echo 安装前端依赖...
    call npm install
    if %errorlevel% neq 0 (
        echo ❌ 前端依赖安装失败
        pause
        exit /b 1
    )
)
echo ✓ 前端依赖就绪

echo 4. 启动前端服务...
echo ✓ 项目启动成功！
echo.
echo 访问地址：
echo - 前端: http://localhost:5173
echo - 后端: http://localhost:8080
echo.
echo 默认管理员账户：
echo - 用户名: admin
echo - 密码: password
echo.
echo 按 Ctrl+C 停止服务

npm run dev