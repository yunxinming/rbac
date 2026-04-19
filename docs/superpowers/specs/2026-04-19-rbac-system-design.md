# RBAC权限管理系统设计文档

## 概述

构建一个独立的权限管理API服务，采用RBAC（基于角色的访问控制）模型，供其他系统调用。

### 项目定位
- **类型**：独立权限管理服务（纯后端API）
- **规模**：小型单租户
- **前端**：无

### 技术栈
- **框架**：Gin
- **认证**：JWT (HS256)
- **数据库**：SQLite
- **ORM**：GORM

---

## 功能范围

### 核心功能

| 模块 | 功能 | 说明 |
|------|------|------|
| 认证 | 登录/登出 | 获取JWT Token |
| 用户管理 | CRUD | 包含角色分配 |
| 角色管理 | CRUD | 包含权限分配 |
| 权限管理 | CRUD | 权限码定义 |
| 权限校验 | 中间件 | API级别访问控制 |

### 排除功能
- 用户注册（权限服务不公开注册）
- 菜单管理（无前端）
- 数据权限/行级权限
- 组织架构/多租户
- Refresh Token机制

---

## 数据模型

### ER图

```
┌──────────────┐       ┌──────────────┐       ┌──────────────┐
│    users     │       │    roles     │       │ permissions  │
├──────────────┤       ├──────────────┤       ├──────────────┤
│ id (PK)      │       │ id (PK)      │       │ id (PK)      │
│ username     │       │ name         │       │ name         │
│ password     │       │ code         │       │ code         │
│ email        │       │ created_at   │       │ created_at   │
│ status       │       │ updated_at   │       │ updated_at   │
│ created_at   │       └──────┬───────┘       └──────┬───────┘
│ updated_at   │              │                      │
└──────┬───────┘              │                      │
       │                      │                      │
       │    ┌─────────────────┴──────────────────────┘
       │    │
       ▼    ▼
┌──────────────┐       ┌───────────────────┐
│  user_roles  │       │ role_permissions  │
├──────────────┤       ├───────────────────┤
│ user_id (FK) │       │ role_id (FK)      │
│ role_id (FK) │       │ permission_id(FK) │
└──────────────┘       └───────────────────┘
```

### 表结构

#### users 用户表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键，自增 |
| username | VARCHAR(64) | 用户名，唯一 |
| password | VARCHAR(128) | 密码(bcrypt哈希) |
| email | VARCHAR(128) | 邮箱，唯一 |
| status | TINYINT | 状态：1启用 0禁用 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

#### roles 角色表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键，自增 |
| name | VARCHAR(64) | 角色名称 |
| code | VARCHAR(64) | 角色编码，唯一 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

#### permissions 权限表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键，自增 |
| name | VARCHAR(64) | 权限名称 |
| code | VARCHAR(64) | 权限编码，唯一，格式：`资源:操作` |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

#### user_roles 用户角色关联表
| 字段 | 类型 | 说明 |
|------|------|------|
| user_id | INTEGER | 用户ID，外键 |
| role_id | INTEGER | 角色ID，外键 |

#### role_permissions 角色权限关联表
| 字段 | 类型 | 说明 |
|------|------|------|
| role_id | INTEGER | 角色ID，外键 |
| permission_id | INTEGER | 权限ID，外键 |

### 权限编码规范

格式：`资源:操作`

示例：
- `user:create` - 创建用户
- `user:read` - 查看用户
- `user:update` - 更新用户
- `user:delete` - 删除用户
- `role:create` - 创建角色
- `role:read` - 查看角色
- ...

---

## API设计

### 基础信息
- 基础路径：`/api`
- 无版本管理
- 认证方式：Bearer Token (JWT)

### 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

错误响应：
```json
{
  "code": 40001,
  "message": "无权限访问",
  "data": null
}
```

### 接口清单

#### 认证接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/auth/login | 登录 | 否 |
| POST | /api/auth/logout | 登出 | 是 |

#### 用户接口

| 方法 | 路径 | 说明 | 权限码 |
|------|------|------|--------|
| GET | /api/users | 用户列表 | user:read |
| GET | /api/users/:id | 用户详情 | user:read |
| POST | /api/users | 创建用户 | user:create |
| PUT | /api/users/:id | 更新用户 | user:update |
| DELETE | /api/users/:id | 删除用户 | user:delete |
| PUT | /api/users/:id/roles | 分配角色 | user:update |

#### 角色接口

| 方法 | 路径 | 说明 | 权限码 |
|------|------|------|--------|
| GET | /api/roles | 角色列表 | role:read |
| GET | /api/roles/:id | 角色详情 | role:read |
| POST | /api/roles | 创建角色 | role:create |
| PUT | /api/roles/:id | 更新角色 | role:update |
| DELETE | /api/roles/:id | 删除角色 | role:delete |
| PUT | /api/roles/:id/permissions | 分配权限 | role:update |

#### 权限接口

| 方法 | 路径 | 说明 | 权限码 |
|------|------|------|--------|
| GET | /api/permissions | 权限列表 | permission:read |
| GET | /api/permissions/:id | 权限详情 | permission:read |
| POST | /api/permissions | 创建权限 | permission:create |
| PUT | /api/permissions/:id | 更新权限 | permission:update |
| DELETE | /api/permissions/:id | 删除权限 | permission:delete |

### 关键接口详情

#### 登录
```
POST /api/auth/login
Request:
{
  "username": "admin",
  "password": "123456"
}

Response:
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "roles": ["admin"]
    }
  }
}
```

---

## 中间件设计

### 中间件链

```
Request → Logger → Recovery → CORS → JWT Auth → Permission → Handler
```

### JWT认证中间件

从请求头 `Authorization: Bearer <token>` 提取Token，解析并验证，将用户信息存入 `gin.Context`。

### 权限校验中间件

采用硬编码方式，使用方法：

```go
// 单个权限
r.GET("/api/users", middleware.RequirePermission("user:read"), handler.ListUsers)

// 多个权限（满足任一）
r.GET("/api/users", middleware.RequireAnyPermission("user:read", "admin:all"))

// 多个权限（需全部满足）
r.GET("/api/users", middleware.RequireAllPermissions("user:read", "user:list"))
```

权限校验流程：
1. 从Context获取当前用户ID
2. 查询用户的全部权限码（缓存优化）
3. 判断是否包含所需权限

---

## 安全设计

### JWT配置

| 配置项 | 值 |
|--------|-----|
| 签名算法 | HS256 |
| Secret | 环境变量 `JWT_SECRET` |
| Token有效期 | 24小时 |
| Payload | user_id, username, roles |

### 密码安全

- 哈希算法：bcrypt
- Cost因子：10
- 密码要求：最少6位（可配置）

### 必须防护

| 威胁 | 防护方式 |
|------|----------|
| SQL注入 | GORM参数化查询 |
| 越权访问 | 权限中间件 |
| 密码泄露 | bcrypt哈希 |
| Token泄露 | HTTPS传输 |

### 可省略防护

- CSRF（无Cookie认证）
- Refresh Token（简化场景）
- 复杂限流（可信调用方）

---

## 初始化设计

### 超级管理员

首次启动时自动初始化：

1. 检查是否存在admin用户
2. 不存在则创建：
   - 用户名：admin
   - 密码：admin123（首次登录强制修改）
   - 角色：超级管理员角色
3. 超级管理员角色拥有所有权限

### 初始数据

```go
// 初始化角色
adminRole := Role{Name: "超级管理员", Code: "admin"}

// 初始权限
permissions := []Permission{
    {Name: "创建用户", Code: "user:create"},
    {Name: "查看用户", Code: "user:read"},
    {Name: "更新用户", Code: "user:update"},
    {Name: "删除用户", Code: "user:delete"},
    // ... 其他权限
}
```

---

## 项目结构

```
hello/
├── main.go                 # 入口
├── config/
│   └── config.go           # 配置管理
├── model/
│   ├── user.go             # 用户模型
│   ├── role.go             # 角色模型
│   └── permission.go       # 权限模型
├── handler/
│   ├── auth.go             # 认证处理器
│   ├── user.go             # 用户处理器
│   ├── role.go             # 角色处理器
│   └── permission.go       # 权限处理器
├── middleware/
│   ├── jwt.go              # JWT认证中间件
│   └── permission.go       # 权限校验中间件
├── service/
│   ├── auth.go             # 认证业务逻辑
│   ├── user.go             # 用户业务逻辑
│   ├── role.go             # 角色业务逻辑
│   └── permission.go       # 权限业务逻辑
├── repository/
│   ├── user.go             # 用户数据访问
│   ├── role.go             # 角色数据访问
│   └── permission.go       # 权限数据访问
├── pkg/
│   ├── jwt/
│   │   └── jwt.go          # JWT工具函数
│   └── response/
│       └── response.go     # 统一响应工具
├── data/
│   └── rbac.db             # SQLite数据库文件
├── go.mod
└── go.sum
```

---

## 依赖清单

```go
require (
    github.com/gin-gonic/gin v1.9+
    github.com/golang-jwt/jwt/v5 v5.0+
    github.com/mattn/go-sqlite3 v1.14+
    golang.org/x/crypto v0.0+  // bcrypt
    gorm.io/gorm v1.25+
    gorm.io/driver/sqlite v1.5+
)
```

---

## 开发命令

```bash
# 运行服务
go run main.go

# 编译
go build -o rbac-server

# 运行测试
go test ./...

# 数据库迁移（自动执行）
# 首次启动时GORM自动创建表
```

---

## 配置项

通过环境变量配置：

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| PORT | 服务端口 | 8080 |
| JWT_SECRET | JWT签名密钥 | 必填 |
| DB_PATH | SQLite数据库路径 | ./data/rbac.db |
