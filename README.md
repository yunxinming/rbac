# RBAC 权限管理系统

基于 Go + Vue 3 的角色权限控制系统，支持动态路由和细粒度权限管理。

## 功能特性

- 🔐 **用户认证**：JWT Token 认证，支持登录/登出
- 👥 **用户管理**：用户增删改查、状态管理、角色分配
- 🎭 **角色管理**：角色增删改查、权限分配
- 🔑 **权限管理**：支持目录/菜单/操作三级权限结构
- 🛣️ **动态路由**：根据用户权限动态生成菜单和路由
- 🎯 **细粒度控制**：按钮级权限控制，无权限自动隐藏

## 技术栈

### 后端
- Go 1.26
- Gin (Web 框架)
- GORM (ORM)
- JWT (认证)
- SQLite (数据库)

### 前端
- Vue 3
- Element Plus (UI 组件库)
- Pinia (状态管理)
- Vue Router (路由)
- Vite (构建工具)

## 项目结构

```
.
├── config/          # 配置
├── handler/         # HTTP 处理器
├── middleware/      # 中间件 (JWT、权限校验)
├── model/           # 数据模型
├── pkg/             # 公共包 (JWT、响应)
├── repository/      # 数据访问层
├── service/         # 业务逻辑层
├── web/             # 前端项目
│   ├── src/
│   │   ├── api/           # API 接口
│   │   ├── components/    # 组件
│   │   ├── directives/    # 自定义指令
│   │   ├── router/        # 路由配置
│   │   ├── store/         # 状态管理
│   │   ├── utils/         # 工具函数
│   │   └── views/         # 页面组件
│   └── ...
├── main.go           # 入口文件
└── go.mod
```

## 快速开始

### 环境要求

- Go 1.26+
- Node.js 18+
- npm 或 yarn

### 后端启动

```bash
# 安装依赖
go mod tidy

# 启动服务
go run main.go
```

后端服务运行在 http://localhost:8080

### 前端启动

```bash
cd web

# 安装依赖
npm install

# 开发模式
npm run dev

# 生产构建
npm run build
```

前端服务运行在 http://localhost:5173

### 默认账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | admin123 | 管理员 |
| user | user123 | 普通用户 |

## 权限设计

### 三级权限结构

```
目录 (directory)
├── 菜单 (menu)
│   └── 操作 (operation)
```

### 权限编码规范

格式：`目录:菜单:操作`

示例：
- `system:user:read` - 用户列表查看
- `system:user:create` - 用户创建
- `system:user:update` - 用户编辑
- `system:user:delete` - 用户删除

### 动态路由

- 用户登录后，后端返回该用户有权限的菜单树
- 前端根据菜单树动态注册路由
- 只有拥有 `:read` 权限的菜单才会显示在侧边栏

## API 接口

### 认证
- `POST /api/auth/login` - 登录
- `POST /api/auth/logout` - 登出
- `GET /api/auth/user` - 获取当前用户信息

### 用户管理
- `GET /api/users` - 获取用户列表
- `POST /api/users` - 创建用户
- `PUT /api/users/:id` - 更新用户
- `DELETE /api/users/:id` - 删除用户
- `POST /api/users/:id/roles` - 分配角色

### 角色管理
- `GET /api/roles` - 获取角色列表
- `POST /api/roles` - 创建角色
- `PUT /api/roles/:id` - 更新角色
- `DELETE /api/roles/:id` - 删除角色
- `POST /api/roles/:id/permissions` - 分配权限

### 权限管理
- `GET /api/permissions` - 获取权限树
- `GET /api/permissions/tree` - 获取权限树（菜单用）
- `POST /api/permissions` - 创建权限
- `PUT /api/permissions/:id` - 更新权限
- `DELETE /api/permissions/:id` - 删除权限

## 前端使用

### 权限指令

```vue
<!-- 按钮级权限控制 -->
<el-button v-permission="'system:user:create'">新增用户</el-button>
```

### 权限判断

```javascript
import { useUserStore } from '@/store/modules/user'

const userStore = useUserStore()

// 判断是否有权限
if (userStore.hasPermission('system:user:delete')) {
  // 有权限
}
```

## License

MIT
