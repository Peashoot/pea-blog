# Pea Blog Backend

使用 Golang 构建的博客系统后端服务。

## 技术栈

- **框架**: Gin Web Framework
- **数据库**: PostgreSQL
- **身份认证**: JWT
- **密码加密**: bcrypt
- **日志**: slog (Go 1.21+)

## 项目结构

```
backend/
├── cmd/server/          # 应用程序入口点
├── internal/
│   ├── config/         # 配置管理
│   ├── handler/        # HTTP 请求处理器
│   ├── middleware/     # 中间件
│   ├── model/          # 数据模型
│   ├── repository/     # 数据访问层
│   ├── service/        # 业务逻辑层
│   └── util/           # 工具函数
├── pkg/
│   ├── database/       # 数据库连接和迁移
│   ├── logger/         # 日志工具
│   └── response/       # 统一响应格式
├── docs/               # API 文档
├── migrations/         # 数据库迁移文件
├── .env.example        # 环境变量示例
└── Makefile           # 构建脚本
```

## 快速开始

### 1. 环境要求

- Go 1.21+
- PostgreSQL 12+

### 2. 安装依赖

```bash
make deps
```

### 3. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，设置数据库连接等配置
```

### 4. 运行应用

开发模式：
```bash
make dev
```

生产模式：
```bash
make build
make run
```

## API 接口

### 认证接口

- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 用户登出
- `GET /api/auth/me` - 获取当前用户信息
- `POST /api/auth/refresh` - 刷新 token

### 文章接口

- `GET /api/articles` - 获取文章列表
- `GET /api/articles/:id` - 获取文章详情
- `GET /api/articles/search` - 搜索文章
- `POST /api/articles` - 创建文章 (需要管理员权限)
- `PUT /api/articles/:id` - 更新文章 (需要管理员权限)
- `DELETE /api/articles/:id` - 删除文章 (需要管理员权限)
- `POST /api/articles/:id/like` - 点赞文章
- `DELETE /api/articles/:id/like` - 取消点赞

### 评论接口

- `GET /api/articles/:id/comments` - 获取文章评论
- `POST /api/comments` - 创建评论
- `DELETE /api/comments/:id` - 删除评论

## 数据库

### 默认管理员账户

- 用户名: `admin`
- 邮箱: `admin@example.com`  
- 密码: `password` (请在生产环境中修改)

### 数据库迁移

项目使用代码中的迁移脚本自动创建表结构，无需额外的迁移工具。

## 开发

### 代码格式化

```bash
make fmt
```

### 代码检查

```bash
make vet
```

### 运行测试

```bash
make test
```

## 部署

1. 构建应用：`make build`
2. 设置生产环境变量
3. 运行：`./bin/server`