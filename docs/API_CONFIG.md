# API 地址配置说明

## 概述

由于前端项目现在由后端服务转发，前端和后端使用相同的域名和端口，因此API请求地址已从绝对路径改为相对路径。

## 配置变更

### 之前（分离部署）
```typescript
// frontend/src/api/client.ts
const API_BASE_URL = 'http://localhost:8080/api'

// 前端: http://localhost:5173
// 后端: http://localhost:8080
// API:  http://localhost:8080/api
```

### 现在（集成部署）
```typescript
// frontend/src/api/client.ts
const API_BASE_URL = '/api'

// 统一访问地址: http://localhost:8080
// 前端页面: http://localhost:8080/
// API接口: http://localhost:8080/api
```

## 环境变量配置

### 开发环境 (.env.development)
```env
VITE_API_URL=/api
```

### 生产环境 (.env.production)
```env
VITE_API_URL=/api
```

## 优势

1. **无跨域问题**: 前后端同源，无需处理CORS
2. **部署简化**: 只需要部署一个服务
3. **配置统一**: 开发和生产环境配置一致
4. **网络优化**: 减少DNS查询和连接建立时间

## 路由处理

后端服务器的路由处理逻辑：

1. **静态文件**: `/assets/*` → 前端静态资源
2. **API接口**: `/api/*` → 后端API处理
3. **前端路由**: `/*` → 返回 index.html，由前端路由处理

## 测试验证

启动服务后，可以通过以下方式验证：

```bash
# 访问首页
curl http://localhost:8080/

# 访问API
curl http://localhost:8080/api/articles

# 访问静态资源
curl http://localhost:8080/assets/index.css
```

## 注意事项

- 开发模式下仍可以分别启动前后端服务
- 集成模式下，前端的API请求会自动路由到后端
- 如需修改API地址，请同时更新环境变量文件和client.ts