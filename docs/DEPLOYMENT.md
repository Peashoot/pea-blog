# Pea Blog 项目部署和使用指南

## 项目概述

Pea Blog 是一个现代化的个人博客系统，采用前后端分离架构开发。具备年轻化科技感的UI设计，完整的博客管理功能，并且完全适配移动端。

### 技术特性

- 🚀 **现代化技术栈**: Vue 3 + TypeScript + Golang
- 📱 **响应式设计**: 完美适配桌面端和移动端
- 🎨 **科技感UI**: 年轻化设计风格，毛玻璃效果，渐变色彩
- 🔐 **完整权限系统**: JWT身份认证，管理员权限控制
- ⚡ **高性能**: Vite构建工具，Go高性能后端
- 🔍 **智能搜索**: 支持标题、内容、标签多维度搜索

## 功能特性

### 核心功能
- ✅ 用户登录和身份认证
- ✅ 文章发布、编辑、删除（管理员）
- ✅ 文章浏览、点赞、评论
- ✅ 智能搜索和筛选
- ✅ 响应式设计，移动端适配

### 技术亮点
- 前端状态管理（Pinia）
- 后端分层架构（Repository-Service-Handler）
- 数据库设计（PostgreSQL）
- API接口规范（RESTful）
- 错误处理和日志记录

## 快速开始

### 环境要求

**前端环境:**
- Node.js 20.19.0+ 或 22.12.0+
- npm 或 yarn

**后端环境:**
- Go 1.21+
- PostgreSQL 12+

### 1. 克隆项目

```bash
git clone <your-repo-url>
cd pea-blog
```

### 2. 数据库设置

创建PostgreSQL数据库：

```sql
CREATE DATABASE pea_blog;
CREATE USER pea_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE pea_blog TO pea_user;
```

### 3. 后端启动

```bash
cd backend

# 复制环境配置
cp .env.example .env

# 编辑 .env 文件，设置数据库连接
# DATABASE_URL=postgres://pea_user:your_password@localhost/pea_blog?sslmode=disable

# 安装依赖
go mod tidy

# 启动服务
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动

### 4. 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端应用将在 `http://localhost:5173` 启动

### 5. 默认账户

系统会自动创建管理员账户：
- 用户名: `admin`
- 密码: `password`

## 项目结构

```
pea-blog/
├── frontend/                 # Vue 3 前端项目
│   ├── src/
│   │   ├── api/             # API 接口
│   │   ├── components/      # Vue 组件
│   │   ├── composables/     # 组合式函数
│   │   ├── stores/          # Pinia 状态管理
│   │   ├── types/           # TypeScript 类型定义
│   │   ├── utils/           # 工具函数
│   │   └── views/           # 页面组件
│   ├── package.json
│   └── vite.config.ts
├── backend/                  # Golang 后端项目
│   ├── cmd/server/          # 应用入口
│   ├── internal/
│   │   ├── config/          # 配置管理
│   │   ├── handler/         # HTTP 处理器
│   │   ├── middleware/      # 中间件
│   │   ├── model/           # 数据模型
│   │   ├── repository/      # 数据访问层
│   │   ├── service/         # 业务逻辑层
│   │   └── util/            # 工具函数
│   ├── pkg/
│   │   ├── database/        # 数据库连接
│   │   ├── logger/          # 日志工具
│   │   └── response/        # 响应格式
│   ├── go.mod
│   └── Makefile
└── README.md
```

## API 接口

### 认证接口
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 用户登出
- `GET /api/auth/me` - 获取当前用户信息

### 文章接口
- `GET /api/articles` - 获取文章列表
- `GET /api/articles/:id` - 获取文章详情
- `GET /api/articles/search` - 搜索文章
- `POST /api/articles` - 创建文章（管理员）
- `PUT /api/articles/:id` - 更新文章（管理员）
- `DELETE /api/articles/:id` - 删除文章（管理员）
- `POST /api/articles/:id/like` - 点赞文章
- `DELETE /api/articles/:id/like` - 取消点赞

### 评论接口
- `GET /api/articles/:id/comments` - 获取文章评论
- `POST /api/comments` - 创建评论
- `DELETE /api/comments/:id` - 删除评论

## 部署

### 前端部署

```bash
cd frontend
npm run build
```

构建产物在 `dist/` 目录，可部署到任何静态文件服务器。

### 后端部署

```bash
cd backend
make build
```

生成的可执行文件在 `bin/server`，设置好环境变量后可直接运行。

## 开发指南

### 前端开发

1. **组件开发**: 所有组件都采用 `<script setup>` 语法
2. **状态管理**: 使用 Pinia 进行状态管理
3. **API调用**: 统一使用 `src/api/` 中的接口
4. **样式规范**: 使用CSS变量，遵循响应式设计原则

### 后端开发

1. **分层架构**: Repository -> Service -> Handler
2. **错误处理**: 统一的错误响应格式
3. **日志记录**: 结构化日志记录
4. **数据验证**: 输入数据验证和类型转换

### 数据库设计

主要数据表：
- `users` - 用户表
- `articles` - 文章表
- `comments` - 评论表
- `likes` - 点赞表

## 特色功能

### 1. 科技感UI设计
- 毛玻璃效果（backdrop-filter）
- 渐变色彩和动画效果
- 现代化的卡片式布局
- 流畅的过渡动画

### 2. 智能搜索
- 支持标题、内容、摘要搜索
- 标签筛选功能
- 多种排序选项
- 实时搜索建议

### 3. 响应式设计
- 移动优先的设计原则
- 灵活的网格布局
- 触摸友好的交互
- 自适应导航菜单

### 4. 权限管理
- JWT令牌认证
- 角色基础的权限控制
- 安全的密码处理
- 自动令牌刷新

## 性能优化

### 前端优化
- 路由懒加载
- 组件按需导入
- 图片懒加载
- 无限滚动加载

### 后端优化
- 数据库索引优化
- 连接池管理
- 响应缓存
- 分页查询

## 安全特性

1. **密码安全**: bcrypt加密存储
2. **XSS防护**: 输入验证和转义
3. **CSRF防护**: CORS配置
4. **SQL注入防护**: 参数化查询
5. **身份验证**: JWT令牌机制

## 扩展建议

### 功能扩展
- [ ] 文章分类管理
- [ ] 用户个人资料
- [ ] 邮件通知系统
- [ ] 文章导出功能
- [ ] 多语言支持

### 技术优化
- [ ] Redis缓存
- [ ] CDN集成
- [ ] 全文搜索（Elasticsearch）
- [ ] 容器化部署（Docker）
- [ ] CI/CD流水线

## 常见问题

### Q: 如何修改默认管理员密码？
A: 登录后可以通过用户设置页面修改，或直接在数据库中更新用户记录。

### Q: 如何自定义主题颜色？
A: 修改 `frontend/src/assets/main.css` 中的CSS变量值。

### Q: 如何添加新的API接口？
A: 按照分层架构在对应的 repository、service、handler 层添加相应代码。

### Q: 数据库迁移失败怎么办？
A: 检查数据库连接配置，确保PostgreSQL服务正常运行，用户权限正确。

## 技术支持

如有问题或建议，请：
1. 查看项目文档
2. 检查issue列表
3. 提交新的issue
4. 联系开发团队

## 开源协议

本项目采用 MIT 协议开源。