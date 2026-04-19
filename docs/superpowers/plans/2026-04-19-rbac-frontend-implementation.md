# RBAC前端实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为RBAC权限管理系统构建Vue前端界面，实现登录、仪表盘、用户管理、角色管理、权限管理功能。

**Architecture:** Vue 3 + Vite构建，使用Pinia状态管理，Vue Router路由控制，Element Plus UI组件库，Axios请求后端API。

**Tech Stack:** Vue 3, Vite, Vue Router 4, Pinia, Element Plus, Axios, SCSS

---

## 文件结构

```
web/
├── index.html                    # HTML入口
├── package.json                  # 依赖配置
├── vite.config.js                # Vite配置
├── .env.development              # 开发环境变量
├── .env.production               # 生产环境变量
├── public/
│   └── favicon.ico
└── src/
    ├── main.js                   # Vue入口
    ├── App.vue                   # 根组件
    ├── router/
    │   └── index.js              # 路由配置
    ├── store/
    │   ├── index.js              # Pinia入口
    │   └── modules/
    │       └── user.js           # 用户状态
    ├── api/
    │   ├── index.js              # Axios实例
    │   ├── auth.js               # 认证API
    │   ├── user.js               # 用户API
    │   ├── role.js               # 角色API
    │   └── permission.js         # 权限API
    ├── utils/
    │   └── auth.js               # Token工具
    ├── views/
    │   ├── Login.vue             # 登录页
    │   ├── Dashboard.vue         # 仪表盘
    │   ├── user/
    │   │   ├── List.vue          # 用户列表
    │   │   └── Form.vue          # 用户表单
    │   ├── role/
    │   │   ├── List.vue          # 角色列表
    │   │   └── Form.vue          # 角色表单
    │   └── permission/
    │       ├── List.vue          # 权限列表
    │       └── Form.vue          # 权限表单
    ├── components/
    │   └── layout/
    │       ├── Layout.vue        # 主布局
    │       ├── Sidebar.vue       # 侧边栏
    │       └── Header.vue        # 顶部栏
    └── styles/
        └── index.scss            # 全局样式
```

---

## Task 1: 项目初始化

**Files:**
- Create: `web/package.json`
- Create: `web/vite.config.js`
- Create: `web/index.html`
- Create: `web/.env.development`
- Create: `web/.env.production`

- [ ] **Step 1: 创建项目目录**

```bash
mkdir -p web/public web/src/{router,store/modules,api,utils,views/{user,role,permission},components/layout,styles}
```

- [ ] **Step 2: 创建 package.json**

```json
{
  "name": "rbac-admin",
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "vue": "^3.4.21",
    "vue-router": "^4.3.0",
    "pinia": "^2.1.7",
    "element-plus": "^2.6.1",
    "axios": "^1.6.7"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.0.4",
    "vite": "^5.1.4",
    "sass": "^1.71.1"
  }
}
```

- [ ] **Step 3: 创建 vite.config.js**

```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    }
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
```

- [ ] **Step 4: 创建 index.html**

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>RBAC权限管理系统</title>
</head>
<body>
  <div id="app"></div>
  <script type="module" src="/src/main.js"></script>
</body>
</html>
```

- [ ] **Step 5: 创建环境变量文件**

`.env.development`:
```
VITE_API_BASE_URL=http://localhost:8080/api
```

`.env.production`:
```
VITE_API_BASE_URL=/api
```

- [ ] **Step 6: 安装依赖**

```bash
cd web && npm install
```

Expected: node_modules目录生成，依赖安装成功

---

## Task 2: Vue入口文件

**Files:**
- Create: `web/src/main.js`
- Create: `web/src/App.vue`
- Create: `web/src/styles/index.scss`

- [ ] **Step 1: 创建 main.js**

```javascript
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'

import App from './App.vue'
import router from './router'
import './styles/index.scss'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ElementPlus, { locale: zhCn })

app.mount('#app')
```

- [ ] **Step 2: 创建 App.vue**

```vue
<template>
  <router-view />
</template>

<script setup>
</script>

<style>
#app {
  height: 100vh;
}
</style>
```

- [ ] **Step 3: 创建全局样式 index.scss**

```scss
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
}

a {
  text-decoration: none;
  color: inherit;
}
```

- [ ] **Step 4: 验证项目启动**

```bash
cd web && npm run dev
```

Expected: 服务启动在 http://localhost:5173

---

## Task 3: 工具函数

**Files:**
- Create: `web/src/utils/auth.js`

- [ ] **Step 1: 创建 Token 工具**

```javascript
const TOKEN_KEY = 'rbac_token'

export function getToken() {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function removeToken() {
  localStorage.removeItem(TOKEN_KEY)
}
```

---

## Task 4: API层

**Files:**
- Create: `web/src/api/index.js`
- Create: `web/src/api/auth.js`
- Create: `web/src/api/user.js`
- Create: `web/src/api/role.js`
- Create: `web/src/api/permission.js`

- [ ] **Step 1: 创建 Axios 实例 api/index.js**

```javascript
import axios from 'axios'
import { getToken, removeToken } from '@/utils/auth'
import { ElMessage } from 'element-plus'
import router from '@/router'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => Promise.reject(error)
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 0) {
      ElMessage.error(res.message || '请求失败')
      if (res.code === 40001 || res.code === 401) {
        removeToken()
        router.push('/login')
      }
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return res
  },
  error => {
    ElMessage.error(error.message || '网络错误')
    if (error.response?.status === 401) {
      removeToken()
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

export default api
```

- [ ] **Step 2: 创建认证API api/auth.js**

```javascript
import api from './index'

export const login = (data) => api.post('/auth/login', data)
export const logout = () => api.post('/auth/logout')
```

- [ ] **Step 3: 创建用户API api/user.js**

```javascript
import api from './index'

export const getUsers = () => api.get('/users')
export const getUser = (id) => api.get(`/users/${id}`)
export const createUser = (data) => api.post('/users', data)
export const updateUser = (id, data) => api.put(`/users/${id}`, data)
export const deleteUser = (id) => api.delete(`/users/${id}`)
export const assignRoles = (id, roleIds) => api.put(`/users/${id}/roles`, { role_ids: roleIds })
```

- [ ] **Step 4: 创建角色API api/role.js**

```javascript
import api from './index'

export const getRoles = () => api.get('/roles')
export const getRole = (id) => api.get(`/roles/${id}`)
export const createRole = (data) => api.post('/roles', data)
export const updateRole = (id, data) => api.put(`/roles/${id}`, data)
export const deleteRole = (id) => api.delete(`/roles/${id}`)
export const assignPermissions = (id, permissionIds) => api.put(`/roles/${id}/permissions`, { permission_ids: permissionIds })
```

- [ ] **Step 5: 创建权限API api/permission.js**

```javascript
import api from './index'

export const getPermissions = () => api.get('/permissions')
export const getPermission = (id) => api.get(`/permissions/${id}`)
export const createPermission = (data) => api.post('/permissions', data)
export const updatePermission = (id, data) => api.put(`/permissions/${id}`, data)
export const deletePermission = (id) => api.delete(`/permissions/${id}`)
```

---

## Task 5: 状态管理

**Files:**
- Create: `web/src/store/index.js`
- Create: `web/src/store/modules/user.js`

- [ ] **Step 1: 创建 Pinia 入口 store/index.js**

```javascript
import { createPinia } from 'pinia'

const pinia = createPinia()

export default pinia
```

- [ ] **Step 2: 创建用户状态 store/modules/user.js**

```javascript
import { defineStore } from 'pinia'
import { getToken, setToken, removeToken } from '@/utils/auth'
import { login as loginApi, logout as logoutApi } from '@/api/auth'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: getToken() || '',
    user: null,
    permissions: []
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.user?.roles?.some(r => r.code === 'admin') || false
  },

  actions: {
    async login(credentials) {
      const res = await loginApi(credentials)
      this.token = res.data.token
      this.user = res.data.user
      setToken(res.data.token)

      // 提取权限码
      if (this.user?.roles) {
        this.permissions = this.user.roles.flatMap(r => r.permissions?.map(p => p.code) || [])
      }
    },

    async logout() {
      try {
        await logoutApi()
      } catch (e) {
        // ignore
      }
      this.token = ''
      this.user = null
      this.permissions = []
      removeToken()
    },

    hasPermission(code) {
      if (this.isAdmin) return true
      return this.permissions.includes(code)
    }
  }
})
```

---

## Task 6: 路由配置

**Files:**
- Create: `web/src/router/index.js`

- [ ] **Step 1: 创建路由配置**

```javascript
import { createRouter, createWebHistory } from 'vue-router'
import { getToken } from '@/utils/auth'
import { useUserStore } from '@/store/modules/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/components/layout/Layout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/dashboard'
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue')
      },
      {
        path: 'users',
        name: 'UserList',
        component: () => import('@/views/user/List.vue'),
        meta: { permission: 'user:read' }
      },
      {
        path: 'roles',
        name: 'RoleList',
        component: () => import('@/views/role/List.vue'),
        meta: { permission: 'role:read' }
      },
      {
        path: 'permissions',
        name: 'PermissionList',
        component: () => import('@/views/permission/List.vue'),
        meta: { permission: 'permission:read' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = getToken()

  if (to.meta.requiresAuth === false) {
    next()
    return
  }

  if (!token) {
    next('/login')
    return
  }

  // 权限检查
  if (to.meta.permission) {
    const userStore = useUserStore()
    if (!userStore.hasPermission(to.meta.permission)) {
      next('/dashboard')
      return
    }
  }

  next()
})

export default router
```

---

## Task 7: 布局组件

**Files:**
- Create: `web/src/components/layout/Layout.vue`
- Create: `web/src/components/layout/Sidebar.vue`
- Create: `web/src/components/layout/Header.vue`

- [ ] **Step 1: 创建主布局 Layout.vue**

```vue
<template>
  <el-container class="layout-container">
    <el-aside width="200px">
      <Sidebar />
    </el-aside>
    <el-container>
      <el-header>
        <Header />
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import Sidebar from './Sidebar.vue'
import Header from './Header.vue'
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.el-aside {
  background-color: #304156;
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #eee;
  padding: 0 20px;
}

.el-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>
```

- [ ] **Step 2: 创建侧边栏 Sidebar.vue**

```vue
<template>
  <div class="sidebar">
    <div class="logo">
      <h2>RBAC管理</h2>
    </div>
    <el-menu
      :default-active="activeMenu"
      class="sidebar-menu"
      background-color="#304156"
      text-color="#bfcbd9"
      active-text-color="#409EFF"
      router
    >
      <el-menu-item index="/dashboard">
        <el-icon><House /></el-icon>
        <span>仪表盘</span>
      </el-menu-item>
      <el-menu-item index="/users" v-if="hasPermission('user:read')">
        <el-icon><User /></el-icon>
        <span>用户管理</span>
      </el-menu-item>
      <el-menu-item index="/roles" v-if="hasPermission('role:read')">
        <el-icon><Avatar /></el-icon>
        <span>角色管理</span>
      </el-menu-item>
      <el-menu-item index="/permissions" v-if="hasPermission('permission:read')">
        <el-icon><Key /></el-icon>
        <span>权限管理</span>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { House, User, Avatar, Key } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/modules/user'

const route = useRoute()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const hasPermission = (code) => userStore.hasPermission(code)
</script>

<style scoped>
.sidebar {
  height: 100%;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #263445;
}

.logo h2 {
  color: #fff;
  font-size: 18px;
  margin: 0;
}

.sidebar-menu {
  border-right: none;
}
</style>
```

- [ ] **Step 3: 创建顶部栏 Header.vue**

```vue
<template>
  <div class="header">
    <div class="left">
      <span class="title">RBAC权限管理系统</span>
    </div>
    <div class="right">
      <el-dropdown @command="handleCommand">
        <span class="user-info">
          <el-icon><User /></el-icon>
          {{ userStore.user?.username }}
          <el-icon class="el-icon--right"><ArrowDown /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="logout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { User, ArrowDown } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/modules/user'
import { ElMessageBox } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const handleCommand = async (command) => {
  if (command === 'logout') {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await userStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.header {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.title {
  font-size: 18px;
  font-weight: 500;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 5px;
}
</style>
```

---

## Task 8: 登录页

**Files:**
- Create: `web/src/views/Login.vue`

- [ ] **Step 1: 创建登录页 Login.vue**

```vue
<template>
  <div class="login-container">
    <div class="login-box">
      <h2>RBAC权限管理系统</h2>
      <el-form :model="form" :rules="rules" ref="formRef" @keyup.enter="handleLogin">
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="用户名" prefix-icon="User" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="密码" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleLogin" style="width: 100%">
            登 录
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/modules/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await userStore.login(form)
    ElMessage.success('登录成功')
    router.push('/dashboard')
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-box {
  width: 400px;
  padding: 40px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.login-box h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}
</style>
```

---

## Task 9: 仪表盘页

**Files:**
- Create: `web/src/views/Dashboard.vue`

- [ ] **Step 1: 创建仪表盘 Dashboard.vue**

```vue
<template>
  <div class="dashboard">
    <h2>欢迎回来，{{ userStore.user?.username }}</h2>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #409EFF">
              <el-icon :size="30"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.users }}</div>
              <div class="stat-label">用户数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #67C23A">
              <el-icon :size="30"><Avatar /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.roles }}</div>
              <div class="stat-label">角色数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #E6A23C">
              <el-icon :size="30"><Key /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.permissions }}</div>
              <div class="stat-label">权限数</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { User, Avatar, Key } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/modules/user'
import { getUsers } from '@/api/user'
import { getRoles } from '@/api/role'
import { getPermissions } from '@/api/permission'

const userStore = useUserStore()

const stats = ref({
  users: 0,
  roles: 0,
  permissions: 0
})

onMounted(async () => {
  try {
    const [usersRes, rolesRes, permsRes] = await Promise.all([
      getUsers(),
      getRoles(),
      getPermissions()
    ])
    stats.value.users = usersRes.data?.length || 0
    stats.value.roles = rolesRes.data?.length || 0
    stats.value.permissions = permsRes.data?.length || 0
  } catch (error) {
    console.error(error)
  }
})
</script>

<style scoped>
.dashboard h2 {
  margin: 0;
  color: #333;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 10px 0;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.stat-info {
  margin-left: 20px;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #999;
  margin-top: 5px;
}
</style>
```

---

## Task 10: 用户管理页

**Files:**
- Create: `web/src/views/user/List.vue`
- Create: `web/src/views/user/Form.vue`

- [ ] **Step 1: 创建用户列表 user/List.vue**

```vue
<template>
  <div class="user-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户列表</span>
          <el-button type="primary" @click="openDialog()" v-if="userStore.hasPermission('user:create')">
            新增用户
          </el-button>
        </div>
      </template>
      <el-table :data="users" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="roles" label="角色">
          <template #default="{ row }">
            <el-tag v-for="role in row.roles" :key="role.id" style="margin-right: 5px">
              {{ role.name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDialog(row)" v-if="userStore.hasPermission('user:update')">
              编辑
            </el-button>
            <el-button link type="danger" @click="handleDelete(row)" v-if="userStore.hasPermission('user:delete')">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 用户表单对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑用户' : '新增用户'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input v-model="form.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="角色" prop="role_ids">
          <el-select v-model="form.role_ids" multiple placeholder="选择角色" style="width: 100%">
            <el-option v-for="role in roles" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/store/modules/user'
import { getUsers, createUser, updateUser, deleteUser, assignRoles } from '@/api/user'
import { getRoles } from '@/api/role'

const userStore = useUserStore()

const loading = ref(false)
const submitting = ref(false)
const users = ref([])
const roles = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const formRef = ref(null)

const form = reactive({
  username: '',
  password: '',
  email: '',
  status: 1,
  role_ids: []
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }, { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const [usersRes, rolesRes] = await Promise.all([getUsers(), getRoles()])
    users.value = usersRes.data || []
    roles.value = rolesRes.data || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const openDialog = (row = null) => {
  isEdit.value = !!row
  currentId.value = row?.id || null

  if (row) {
    form.username = row.username
    form.email = row.email
    form.status = row.status
    form.role_ids = row.roles?.map(r => r.id) || []
  } else {
    form.username = ''
    form.password = ''
    form.email = ''
    form.status = 1
    form.role_ids = []
  }

  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      await updateUser(currentId.value, {
        username: form.username,
        email: form.email,
        status: form.status
      })
      await assignRoles(currentId.value, form.role_ids)
      ElMessage.success('更新成功')
    } else {
      const res = await createUser({
        username: form.username,
        password: form.password,
        email: form.email,
        status: form.status
      })
      if (form.role_ids.length > 0) {
        await assignRoles(res.data.id, form.role_ids)
      }
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定要删除该用户吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
  await deleteUser(row.id)
  ElMessage.success('删除成功')
  fetchData()
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
```

---

## Task 11: 角色管理页

**Files:**
- Create: `web/src/views/role/List.vue`

- [ ] **Step 1: 创建角色列表 role/List.vue**

```vue
<template>
  <div class="role-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>角色列表</span>
          <el-button type="primary" @click="openDialog()" v-if="userStore.hasPermission('role:create')">
            新增角色
          </el-button>
        </div>
      </template>
      <el-table :data="roles" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" />
        <el-table-column prop="code" label="角色编码" />
        <el-table-column label="权限数" width="100">
          <template #default="{ row }">
            {{ row.permissions?.length || 0 }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDialog(row)" v-if="userStore.hasPermission('role:update')">
              编辑
            </el-button>
            <el-button link type="primary" @click="openPermissionDialog(row)" v-if="userStore.hasPermission('role:update')">
              分配权限
            </el-button>
            <el-button link type="danger" @click="handleDelete(row)" v-if="userStore.hasPermission('role:delete')">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 角色表单对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑角色' : '新增角色'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="角色编码" prop="code">
          <el-input v-model="form.code" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 分配权限对话框 -->
    <el-dialog v-model="permissionDialogVisible" title="分配权限" width="500px">
      <el-checkbox-group v-model="selectedPermissions">
        <el-checkbox v-for="perm in permissions" :key="perm.id" :value="perm.id" style="display: block; margin-bottom: 10px">
          {{ perm.name }} ({{ perm.code }})
        </el-checkbox>
      </el-checkbox-group>
      <template #footer>
        <el-button @click="permissionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAssignPermissions" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/store/modules/user'
import { getRoles, createRole, updateRole, deleteRole, assignPermissions } from '@/api/role'
import { getPermissions } from '@/api/permission'

const userStore = useUserStore()

const loading = ref(false)
const submitting = ref(false)
const roles = ref([])
const permissions = ref([])
const dialogVisible = ref(false)
const permissionDialogVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const formRef = ref(null)
const selectedPermissions = ref([])

const form = reactive({
  name: '',
  code: ''
})

const rules = {
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入角色编码', trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const [rolesRes, permsRes] = await Promise.all([getRoles(), getPermissions()])
    roles.value = rolesRes.data || []
    permissions.value = permsRes.data || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const openDialog = (row = null) => {
  isEdit.value = !!row
  currentId.value = row?.id || null

  if (row) {
    form.name = row.name
    form.code = row.code
  } else {
    form.name = ''
    form.code = ''
  }

  dialogVisible.value = true
}

const openPermissionDialog = (row) => {
  currentId.value = row.id
  selectedPermissions.value = row.permissions?.map(p => p.id) || []
  permissionDialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      await updateRole(currentId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createRole(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const handleAssignPermissions = async () => {
  submitting.value = true
  try {
    await assignPermissions(currentId.value, selectedPermissions.value)
    ElMessage.success('分配成功')
    permissionDialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定要删除该角色吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
  await deleteRole(row.id)
  ElMessage.success('删除成功')
  fetchData()
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
```

---

## Task 12: 权限管理页

**Files:**
- Create: `web/src/views/permission/List.vue`

- [ ] **Step 1: 创建权限列表 permission/List.vue**

```vue
<template>
  <div class="permission-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>权限列表</span>
          <el-button type="primary" @click="openDialog()" v-if="userStore.hasPermission('permission:create')">
            新增权限
          </el-button>
        </div>
      </template>
      <el-table :data="permissions" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="权限名称" />
        <el-table-column prop="code" label="权限编码" />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDialog(row)" v-if="userStore.hasPermission('permission:update')">
              编辑
            </el-button>
            <el-button link type="danger" @click="handleDelete(row)" v-if="userStore.hasPermission('permission:delete')">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 权限表单对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑权限' : '新增权限'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="权限名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="权限编码" prop="code">
          <el-input v-model="form.code" placeholder="格式：资源:操作，如 user:create" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/store/modules/user'
import { getPermissions, createPermission, updatePermission, deletePermission } from '@/api/permission'

const userStore = useUserStore()

const loading = ref(false)
const submitting = ref(false)
const permissions = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const formRef = ref(null)

const form = reactive({
  name: '',
  code: ''
})

const rules = {
  name: [{ required: true, message: '请输入权限名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入权限编码', trigger: 'blur' }]
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getPermissions()
    permissions.value = res.data || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const openDialog = (row = null) => {
  isEdit.value = !!row
  currentId.value = row?.id || null

  if (row) {
    form.name = row.name
    form.code = row.code
  } else {
    form.name = ''
    form.code = ''
  }

  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      await updatePermission(currentId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createPermission(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定要删除该权限吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
  await deletePermission(row.id)
  ElMessage.success('删除成功')
  fetchData()
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
```

---

## Task 13: 后端CORS配置

**Files:**
- Modify: `main.go`

- [ ] **Step 1: 添加CORS中间件**

在 `main.go` 中添加CORS配置：

```go
// 在 r.Use(gin.Recovery()) 后添加
r.Use(func(c *gin.Context) {
    c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
    c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
    c.Header("Access-Control-Expose-Headers", "Content-Length")
    c.Header("Access-Control-Allow-Credentials", "true")

    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }

    c.Next()
})
```

- [ ] **Step 2: 重启后端服务**

```bash
# 在项目根目录
go build -o rbac-server . && JWT_SECRET=rbac-secret-key-2024 ./rbac-server
```

---

## Task 14: 启动测试

**Files:**
- 无新文件

- [ ] **Step 1: 启动前端开发服务**

```bash
cd web && npm run dev
```

Expected: 服务启动在 http://localhost:5173

- [ ] **Step 2: 测试登录**

1. 访问 http://localhost:5173
2. 使用 admin/admin123 登录
3. 验证跳转到仪表盘

- [ ] **Step 3: 测试功能**

1. 测试用户管理：创建、编辑、删除用户
2. 测试角色管理：创建角色、分配权限
3. 测试权限管理：创建、编辑、删除权限
4. 测试退出登录

---

## 实现顺序总结

1. **项目初始化** (Task 1-2): 目录、配置、入口文件
2. **基础设施** (Task 3-6): 工具函数、API、状态管理、路由
3. **布局组件** (Task 7): Layout、Sidebar、Header
4. **页面开发** (Task 8-12): 登录、仪表盘、用户、角色、权限
5. **集成测试** (Task 13-14): CORS配置、启动测试
