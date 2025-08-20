# Pea Blog - 个人博客系统

一个现代化的个人博客系统，采用前后端分离架构，支持Docker部署。

## 技术栈

**前端:**
- Vue 3 + TypeScript
- Vite 构建工具
- Element Plus UI 组件库
- Pinia 状态管理
- Vue Router 路由管理

**后端:**
- Go 1.23 + Gin Web框架
- SQLite/PostgreSQL 数据库支持
- JWT 身份验证
- RESTful API 设计

## 功能特性

- ✅ 管理员登录和权限管理
- ✅ 文章创建、编辑、删除和发布
- ✅ 文章预览和详情查看
- ✅ 文章点赞和评论系统
- ✅ 文章搜索和标签过滤
- ✅ 移动端响应式适配
- ✅ 现代化科技感UI设计
- ✅ **自动前端构建** - 服务器启动时自动检测并构建前端项目
- ✅ **可配置构建** - 支持开启/关闭自动构建，自定义构建命令

## 快速开始

### 方式一：Docker 部署（推荐）

1. 确保已安装 Docker 和 Docker Compose

2. 克隆项目并进入目录
```bash
git clone <repository-url>
cd pea-blog
```

3. 使用部署脚本启动服务
```bash
# Linux/Mac
chmod +x deploy.sh
./deploy.sh

# Windows
deploy.bat
```

4. 访问 http://localhost:8080

### 方式二：本地开发

#### 后端启动
```bash
cd backend
go mod download
go run cmd/server/main.go
```

#### 前端启动（开发模式）
```bash
cd frontend
npm install
npm run dev
```

### 方式三：集成模式启动

已集成前端到后端服务中，只需启动后端即可：

```bash
# 构建前端（可选，服务器会自动构建）
cd frontend
npm run build-only

# 启动后端（包含前端静态文件服务）
cd ../backend
go run cmd/server/main.go
```

访问 http://localhost:8080 即可使用完整的博客系统。

**优势：**
- 🌐 **统一域名**: 前后端使用相同域名，无跨域问题
- 🔄 **自动构建**: 后端服务启动时自动检测并构建前端
- 📡 **API代理**: 前端API请求直接使用相对路径 `/api`
- 🚀 **一键启动**: 只需启动一个服务即可运行完整系统

## 默认账户

- 用户名: `admin`
- 密码: `password`

**⚠️ 生产环境请务必修改默认密码！**

## 环境配置

### 开发环境
后端配置文件位于 `backend/.env`:
```env
PORT=8080
DATABASE_URL=./pea_blog.db
JWT_SECRET=development-jwt-secret-key
JWT_EXPIRE_HOURS=24
ENVIRONMENT=development

# 前端自动构建配置
FRONTEND_AUTO_BUILD=true
FRONTEND_BUILD_COMMAND=npm run build
FRONTEND_DIST_PATH=../frontend/dist
FRONTEND_SOURCE_PATH=../frontend
```

### 生产环境
生产环境配置示例文件 `.env.production`:
```env
PORT=8080
DATABASE_URL=/app/data/pea_blog.db
JWT_SECRET=your-secure-jwt-secret-here
JWT_EXPIRE_HOURS=24
ENVIRONMENT=production

# 前端自动构建配置（生产环境建议关闭）
FRONTEND_AUTO_BUILD=false
FRONTEND_BUILD_COMMAND=npm run build
FRONTEND_DIST_PATH=./frontend/dist
FRONTEND_SOURCE_PATH=./frontend
```

## Docker 配置

### docker-compose.yml
```yaml
version: '3.8'
services:
  pea-blog:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - pea_blog_data:/app/data
    restart: unless-stopped
```

### 常用 Docker 命令
```bash
# 构建和启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 重启服务
docker-compose restart
```

## API 接口

### 认证接口
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 用户登出
- `GET /api/auth/me` - 获取当前用户信息

### 文章接口
- `GET /api/articles` - 获取文章列表
- `GET /api/articles/:id` - 获取文章详情
- `POST /api/articles` - 创建文章（需要管理员权限）
- `PUT /api/articles/:id` - 更新文章（需要管理员权限）
- `DELETE /api/articles/:id` - 删除文章（需要管理员权限）
- `POST /api/articles/:id/like` - 点赞文章
- `DELETE /api/articles/:id/like` - 取消点赞

### 评论接口
- `GET /api/articles/:id/comments` - 获取文章评论
- `POST /api/comments` - 创建评论
- `DELETE /api/comments/:id` - 删除评论

### 系统接口
- `GET /api/system/build-status` - 获取构建状态
- `POST /api/system/rebuild-frontend` - 手动重新构建前端（需要管理员权限）

## 项目结构

```
pea-blog/
├── backend/                # 后端代码
│   ├── cmd/server/        # 服务入口
│   ├── internal/          # 内部包
│   │   ├── config/        # 配置管理
│   │   ├── handler/       # HTTP处理器
│   │   ├── middleware/    # 中间件
│   │   ├── model/         # 数据模型
│   │   ├── repository/    # 数据访问层
│   │   └── service/       # 业务逻辑层
│   └── pkg/               # 公共包
├── frontend/              # 前端代码
│   ├── src/
│   │   ├── components/    # 组件
│   │   ├── views/         # 页面
│   │   ├── router/        # 路由配置
│   │   ├── stores/        # 状态管理
│   │   └── utils/         # 工具函数
│   └── dist/              # 构建输出
├── Dockerfile             # Docker镜像构建文件
├── docker-compose.yml     # Docker编排文件
└── deploy.sh/.bat         # 部署脚本
```

## 开发说明

### 自动构建功能
后端服务集成了智能的前端自动构建功能：

**功能特性：**
- 🔍 **智能检测**: 启动时自动检查前端源码是否有更新
- ⚡ **按需构建**: 只有在源码变更时才重新构建，避免不必要的构建
- 🎛️ **可配置控制**: 支持通过环境变量开启/关闭自动构建
- 📋 **详细日志**: 提供构建过程的详细日志输出
- 🔧 **手动触发**: 提供API接口手动触发重新构建

**配置说明：**
```bash
# 开启/关闭自动构建
FRONTEND_AUTO_BUILD=true

# 自定义构建命令
FRONTEND_BUILD_COMMAND="npm run build"

# 配置路径
FRONTEND_DIST_PATH="../frontend/dist"
FRONTEND_SOURCE_PATH="../frontend"
```

**使用场景：**
- **开发环境**: 开启自动构建，源码变更时自动重新构建
- **生产环境**: 关闭自动构建，使用预构建的静态文件
- **Docker环境**: 构建时已包含前端文件，无需自动构建

### 添加新功能
1. 后端：在相应的 handler、service、repository 层添加代码
2. 前端：在 components、views、stores 中添加相应组件和逻辑
3. 更新路由配置和API接口

### 数据库迁移
数据库迁移会在服务启动时自动执行，迁移文件位于 `backend/pkg/database/` 目录。

## 部署建议

1. **安全性**: 修改默认JWT密钥和管理员密码
2. **数据库**: 生产环境建议使用PostgreSQL
3. **反向代理**: 建议使用Nginx作为反向代理
4. **HTTPS**: 生产环境请配置SSL证书
5. **监控**: 建议添加日志监控和健康检查

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

MIT License