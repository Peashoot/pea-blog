@echo off
echo === Pea Blog Docker 部署脚本 ===

echo 1. 检查Docker环境...
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker 未安装，请先安装Docker Desktop
    pause
    exit /b 1
)

docker-compose --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker Compose 未安装，请先安装Docker Compose
    pause
    exit /b 1
)

echo ✓ Docker 环境检查通过

echo 2. 构建Docker镜像...
docker-compose build
if %errorlevel% neq 0 (
    echo ❌ Docker镜像构建失败
    pause
    exit /b 1
)
echo ✓ Docker镜像构建成功

echo 3. 启动服务...
docker-compose up -d
if %errorlevel% neq 0 (
    echo ❌ 服务启动失败
    pause
    exit /b 1
)

echo ✓ 服务启动成功
echo.
echo 访问地址：
echo - 网站: http://localhost:8080
echo - API: http://localhost:8080/api
echo.
echo 默认管理员账户：
echo - 用户名: admin
echo - 密码: password
echo.
echo 常用命令：
echo - 查看日志: docker-compose logs -f
echo - 停止服务: docker-compose down
echo - 重启服务: docker-compose restart
echo.
pause