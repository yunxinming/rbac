# 权限表扩展设计文档

## 概述

扩展现有权限表，支持目录、菜单、按钮、API三种权限类型，实现前端动态路由和细粒度权限控制。

### 设计目标
- 支持三层权限结构：目录 → 菜单 → 操作（按钮/API）
- 权限编码采用层级格式
- 前端根据用户权限动态生成路由
- 后端API权限自动校验

---

## 权限类型层级

```
目录 (directory)
└── 菜单 (menu)
    └── 操作 (operation) - 按钮/API
```

**示例结构**：
```
系统管理 [目录]
├── 用户管理 [菜单]
│   ├── 查看用户 [操作-按钮]
│   ├── 创建用户 [操作-按钮]
│   ├── 编辑用户 [操作-按钮]
│   └── 删除用户 [操作-按钮]
├── 角色管理 [菜单]
│   └── ...
└── 权限管理 [菜单]
    └── ...
```

---

## 权限编码规范

### 格式
`目录:菜单:操作`

### 示例

| 权限 | 编码 | 类型 |
|------|------|------|
| 系统管理 | `system` | 目录 |
| 用户管理 | `system:user` | 菜单 |
| 创建用户 | `system:user:create` | 操作 |
| 用户列表API | `system:user:list` | 操作 |

---

## 数据模型

### Permission 表结构

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| id | uint | 主键，自增 | 1 |
| name | varchar(64) | 权限名称 | 用户管理 |
| code | varchar(128) | 权限编码，唯一 | system:user |
| type | varchar(20) | 类型：directory/menu/operation | menu |
| parent_id | uint | 父级ID，0为顶级 | 1 |
| path | varchar(128) | 前端路由路径 | /users |
| icon | varchar(64) | 图标名称 | User |
| sort | int | 排序（升序） | 1 |
| api_path | varchar(255) | API路径（操作类型） | /api/users |
| api_method | varchar(10) | API方法（操作类型） | GET |
| component | varchar(128) | 前端组件路径 | user/List |
| status | tinyint | 状态：1启用 0禁用 | 1 |
| created_at | datetime | 创建时间 | |
| updated_at | datetime | 更新时间 | |
| deleted_at | datetime | 删除时间（软删除） | |

### 字段说明

**按类型区分必填字段**：

| 字段 | directory | menu | operation |
|------|-----------|------|-----------|
| name | ✅ | ✅ | ✅ |
| code | ✅ | ✅ | ✅ |
| type | ✅ | ✅ | ✅ |
| parent_id | 0 | 目录ID | 菜单ID |
| path | - | ✅ | - |
| icon | 可选 | 可选 | - |
| sort | ✅ | ✅ | ✅ |
| api_path | - | - | ✅ |
| api_method | - | - | ✅ |
| component | - | ✅ | - |

---

## API设计

### 新增接口

#### 获取用户菜单树

```
GET /api/auth/menus

Response:
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "系统管理",
      "code": "system",
      "type": "directory",
      "icon": "Setting",
      "sort": 1,
      "children": [
        {
          "id": 2,
          "name": "用户管理",
          "code": "system:user",
          "type": "menu",
          "path": "/users",
          "icon": "User",
          "component": "user/List",
          "sort": 1,
          "children": []
        }
      ]
    }
  ]
}
```

#### 获取权限树（管理用）

```
GET /api/permissions/tree
```

#### 创建权限

```
POST /api/permissions
{
  "name": "用户管理",
  "code": "system:user",
  "type": "menu",
  "parent_id": 1,
  "path": "/users",
  "icon": "User",
  "component": "user/List",
  "sort": 1
}
```

#### 更新权限

```
PUT /api/permissions/:id
```

#### 删除权限

```
DELETE /api/permissions/:id
```

---

## 后端变更

### 1. Model 修改

`model/permission.go` 添加新字段：
- Type、ParentID、Path、Icon、Sort、ApiPath、ApiMethod、Component、Status

### 2. Repository 修改

`repository/permission.go` 添加：
- GetMenuTree() - 获取完整权限树
- GetUserMenus(userID) - 获取用户有权限的菜单树
- GetByType(permType) - 按类型查询

### 3. Service 修改

`service/permission.go` 添加：
- GetMenuTree() - 构建树形结构
- GetUserMenus(userID) - 根据用户权限过滤菜单

`service/auth.go` 添加：
- GetUserMenus() - 登录后返回菜单树

### 4. 中间件修改

`middleware/permission.go` 增强：
- 支持API权限自动匹配
- 根据请求路径和方法自动校验权限

---

## 前端变更

### 1. 路由改造

`router/index.js`:
- 移除静态路由定义
- 添加动态路由注册逻辑
- 从后端获取菜单并生成路由

### 2. 状态管理

`store/modules/user.js`:
- 添加 menus 状态
- 登录时获取菜单树
- 动态生成路由

### 3. 菜单组件

`components/layout/Sidebar.vue`:
- 改为从 store 获取菜单
- 递归渲染多级菜单

### 4. 权限指令

新增 `directives/permission.js`:
- `v-permission="'system:user:create'"`
- 根据权限控制按钮显示

---

## 初始数据

```go
permissions := []model.Permission{
    // 目录
    {Name: "系统管理", Code: "system", Type: "directory", Icon: "Setting", Sort: 1},
    
    // 菜单
    {Name: "用户管理", Code: "system:user", Type: "menu", ParentID: 1, Path: "/users", Icon: "User", Component: "user/List", Sort: 1},
    {Name: "角色管理", Code: "system:role", Type: "menu", ParentID: 1, Path: "/roles", Icon: "Avatar", Component: "role/List", Sort: 2},
    {Name: "权限管理", Code: "system:permission", Type: "menu", ParentID: 1, Path: "/permissions", Icon: "Key", Component: "permission/List", Sort: 3},
    
    // 用户操作
    {Name: "查看用户", Code: "system:user:read", Type: "operation", ParentID: 2, ApiPath: "/api/users", ApiMethod: "GET", Sort: 1},
    {Name: "创建用户", Code: "system:user:create", Type: "operation", ParentID: 2, ApiPath: "/api/users", ApiMethod: "POST", Sort: 2},
    {Name: "更新用户", Code: "system:user:update", Type: "operation", ParentID: 2, ApiPath: "/api/users/:id", ApiMethod: "PUT", Sort: 3},
    {Name: "删除用户", Code: "system:user:delete", Type: "operation", ParentID: 2, ApiPath: "/api/users/:id", ApiMethod: "DELETE", Sort: 4},
    
    // 角色操作
    // ... 类似
    
    // 权限操作
    // ... 类似
}
```

---

## 权限校验流程

### 前端流程

```
1. 用户登录 → 获取Token和菜单树
2. 根据菜单树动态注册路由
3. 侧边栏根据菜单树渲染
4. 按钮通过 v-permission 指令控制显示
```

### 后端流程

```
1. 请求到达 → JWT中间件验证Token
2. 权限中间件获取请求路径和方法
3. 匹配用户权限中的 api_path + api_method
4. 匹配成功放行，失败返回403
```

---

## 实现步骤

### 后端
1. 修改 Permission 模型
2. 数据库迁移
3. 修改 Repository 层
4. 修改 Service 层
5. 修改 Handler 层
6. 增强权限中间件
7. 更新初始数据

### 前端
1. 添加权限指令
2. 改造路由为动态路由
3. 改造状态管理
4. 改造侧边栏组件
5. 添加权限管理页面（树形）

---

## 兼容性说明

- 现有的 `user:read` 格式权限编码仍然有效
- 新的 `system:user:read` 格式向后兼容
- 前端 `hasPermission()` 方法支持两种格式
