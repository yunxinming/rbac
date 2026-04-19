# RBAC前端设计文档

## 概述

为RBAC权限管理系统构建Vue前端界面，对接已有后端API。

### 项目定位
- **位置**：`web/` 目录（与后端同项目）
- **类型**：完整管理后台
- **用户**：系统管理员

### 技术栈
- **框架**：Vue 3 (Composition API)
- **构建**：Vite
- **路由**：Vue Router 4
- **状态**：Pinia
- **UI库**：Element Plus
- **HTTP**：Axios
- **样式**：SCSS

---

## 项目结构

```
web/
├── index.html
├── package.json
├── vite.config.js
├── .env.development
├── .env.production
├── public/
│   └── favicon.ico
└── src/
    ├── main.js
    ├── App.vue
    ├── router/
    │   └── index.js
    ├── store/
    │   ├── index.js
    │   └── modules/
    │       ├── user.js
    │       └── permission.js
    ├── api/
    │   ├── index.js
    │   ├── auth.js
    │   ├── user.js
    │   ├── role.js
    │   └── permission.js
    ├── utils/
    │   └── auth.js
    ├── views/
    │   ├── Login.vue
    │   ├── Dashboard.vue
    │   ├── user/
    │   │   ├── List.vue
    │   │   └── Form.vue
    │   ├── role/
    │   │   ├── List.vue
    │   │   └── Form.vue
    │   └── permission/
    │       ├── List.vue
    │       └── Form.vue
    ├── components/
    │   └── layout/
    │       ├── Layout.vue
    │       ├── Sidebar.vue
    │       └── Header.vue
    └── styles/
        └── index.scss
```

---

## 页面设计

### 页面列表

| 页面 | 路由 | 说明 | 权限 |
|------|------|------|------|
| 登录 | /login | 用户登录 | 无 |
| 仪表盘 | /dashboard | 首页统计 | 登录即可 |
| 用户管理 | /users | 用户CRUD | user:read |
| 角色管理 | /roles | 角色CRUD | role:read |
| 权限管理 | /permissions | 权限CRUD | permission:read |

### 页面布局

```
┌─────────────────────────────────────────────────┐
│  Header (用户信息、退出按钮)                      │
├──────────┬──────────────────────────────────────┤
│          │                                      │
│ Sidebar  │           Main Content               │
│ (菜单)   │           (路由视图)                  │
│          │                                      │
│          │                                      │
└──────────┴──────────────────────────────────────┘
```

---

## 路由设计

### 路由配置

```javascript
const routes = [
  { path: '/login', component: Login, meta: { requiresAuth: false } },
  {
    path: '/',
    component: Layout,
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: '/dashboard' },
      { path: 'dashboard', component: Dashboard },
      { path: 'users', component: UserList, meta: { permission: 'user:read' } },
      { path: 'users/create', component: UserForm, meta: { permission: 'user:create' } },
      { path: 'users/:id/edit', component: UserForm, meta: { permission: 'user:update' } },
      { path: 'roles', component: RoleList, meta: { permission: 'role:read' } },
      { path: 'roles/create', component: RoleForm, meta: { permission: 'role:create' } },
      { path: 'roles/:id/edit', component: RoleForm, meta: { permission: 'role:update' } },
      { path: 'permissions', component: PermissionList, meta: { permission: 'permission:read' } },
      { path: 'permissions/create', component: PermissionForm, meta: { permission: 'permission:create' } },
      { path: 'permissions/:id/edit', component: PermissionForm, meta: { permission: 'permission:update' } },
    ]
  }
]
```

### 路由守卫

```javascript
router.beforeEach((to, from, next) => {
  const token = getToken()
  const userStore = useUserStore()

  if (to.meta.requiresAuth !== false && !token) {
    next('/login')
  } else if (to.meta.permission && !userStore.hasPermission(to.meta.permission)) {
    next('/403')
  } else {
    next()
  }
})
```

---

## API设计

### Axios实例

```javascript
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000
})

// 请求拦截器：添加Token
api.interceptors.request.use(config => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器：处理错误
api.interceptors.response.use(
  response => response.data,
  error => {
    if (error.response?.status === 401) {
      removeToken()
      router.push('/login')
    }
    return Promise.reject(error)
  }
)
```

### API模块

```javascript
// api/auth.js
export const login = (data) => api.post('/auth/login', data)
export const logout = () => api.post('/auth/logout')

// api/user.js
export const getUsers = () => api.get('/users')
export const getUser = (id) => api.get(`/users/${id}`)
export const createUser = (data) => api.post('/users', data)
export const updateUser = (id, data) => api.put(`/users/${id}`, data)
export const deleteUser = (id) => api.delete(`/users/${id}`)
export const assignRoles = (id, roleIds) => api.put(`/users/${id}/roles`, { role_ids: roleIds })

// api/role.js
export const getRoles = () => api.get('/roles')
export const getRole = (id) => api.get(`/roles/${id}`)
export const createRole = (data) => api.post('/roles', data)
export const updateRole = (id, data) => api.put(`/roles/${id}`, data)
export const deleteRole = (id) => api.delete(`/roles/${id}`)
export const assignPermissions = (id, permissionIds) => api.put(`/roles/${id}/permissions`, { permission_ids: permissionIds })

// api/permission.js
export const getPermissions = () => api.get('/permissions')
export const getPermission = (id) => api.get(`/permissions/${id}`)
export const createPermission = (data) => api.post('/permissions', data)
export const updatePermission = (id, data) => api.put(`/permissions/${id}`, data)
export const deletePermission = (id) => api.delete(`/permissions/${id}`)
```

---

## 状态管理

### User Store

```javascript
export const useUserStore = defineStore('user', {
  state: () => ({
    token: getToken(),
    user: null,
    permissions: []
  }),
  actions: {
    async login(credentials) {
      const res = await login(credentials)
      this.token = res.data.token
      this.user = res.data.user
      setToken(res.data.token)
    },
    async logout() {
      await logout()
      this.token = null
      this.user = null
      this.permissions = []
      removeToken()
    },
    hasPermission(code) {
      if (this.user?.roles?.includes('admin')) return true
      return this.permissions.includes(code)
    }
  }
})
```

---

## 组件设计

### Layout组件

主布局组件，包含Header、Sidebar、Main三个区域。

**Header**：
- Logo/系统名称
- 当前用户信息
- 退出按钮

**Sidebar**：
- 动态菜单（根据权限显示）
- 菜单项：仪表盘、用户管理、角色管理、权限管理

### 列表页面通用结构

```
┌─────────────────────────────────────┐
│  搜索栏 + 新增按钮                    │
├─────────────────────────────────────┤
│                                     │
│  Element Plus Table                 │
│  (数据表格 + 操作列)                  │
│                                     │
├─────────────────────────────────────┤
│  分页器                              │
└─────────────────────────────────────┘
```

### 表单页面通用结构

```
┌─────────────────────────────────────┐
│  页面标题                            │
├─────────────────────────────────────┤
│                                     │
│  Element Plus Form                  │
│  (表单字段 + 验证规则)                │
│                                     │
├─────────────────────────────────────┤
│  提交按钮 + 取消按钮                  │
└─────────────────────────────────────┘
```

---

## 页面详情

### 登录页 (Login.vue)

- 用户名输入框
- 密码输入框
- 登录按钮
- 登录成功后跳转到仪表盘

### 仪表盘 (Dashboard.vue)

- 欢迎信息
- 统计卡片：用户数、角色数、权限数
- 最近活动（可选）

### 用户管理 (user/List.vue)

- 用户列表表格：用户名、邮箱、状态、角色、操作
- 操作按钮：编辑、删除、分配角色
- 创建用户弹窗/页面

### 角色管理 (role/List.vue)

- 角色列表表格：角色名、编码、权限数、操作
- 操作按钮：编辑、删除、分配权限
- 创建角色弹窗/页面

### 权限管理 (permission/List.vue)

- 权限列表表格：权限名、编码、操作
- 操作按钮：编辑、删除
- 创建权限弹窗/页面

---

## 环境配置

### .env.development

```
VITE_API_BASE_URL=http://localhost:8080/api
```

### .env.production

```
VITE_API_BASE_URL=/api
```

---

## 开发命令

```bash
# 安装依赖
cd web && npm install

# 开发模式
npm run dev

# 构建生产版本
npm run build

# 预览生产版本
npm run preview
```

---

## 后端CORS配置

需要在后端添加CORS中间件，允许前端跨域访问：

```go
// main.go 中添加
r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
}))
```

---

## 依赖清单

```json
{
  "dependencies": {
    "vue": "^3.4",
    "vue-router": "^4.3",
    "pinia": "^2.1",
    "element-plus": "^2.6",
    "axios": "^1.6"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.0",
    "vite": "^5.1",
    "sass": "^1.71"
  }
}
```
