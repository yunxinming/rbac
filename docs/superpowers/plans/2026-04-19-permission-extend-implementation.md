# 权限表扩展实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 扩展权限表支持目录/菜单/操作三层结构，实现前端动态路由和API权限自动校验。

**Architecture:** 后端Permission模型增加类型、父级、路由、API等字段；Repository/Service层增加树形查询；中间件支持API路径匹配；前端实现动态路由和权限指令。

**Tech Stack:** Go + GORM (后端), Vue 3 + Vue Router (前端)

---

## 文件结构

### 后端修改
```
model/permission.go        # 添加新字段
repository/permission.go   # 添加树形查询方法
service/permission.go      # 添加树形构建逻辑
service/auth.go           # 添加获取用户菜单方法
handler/auth.go           # 添加获取菜单接口
handler/permission.go     # 添加获取权限树接口
middleware/permission.go  # 增强API权限匹配
repository/db.go          # 更新初始数据
```

### 前端修改
```
src/directives/permission.js    # 新增：权限指令
src/store/modules/user.js       # 修改：添加菜单状态和动态路由
src/router/index.js             # 修改：动态路由注册
src/components/layout/Sidebar.vue  # 修改：从store获取菜单
src/views/permission/List.vue   # 修改：树形权限管理
src/main.js                     # 修改：注册权限指令
```

---

## Task 1: 修改Permission模型

**Files:**
- Modify: `model/permission.go`

- [ ] **Step 1: 更新Permission模型**

将 `model/permission.go` 替换为：

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

// PermissionType 权限类型
type PermissionType string

const (
	PermissionTypeDirectory PermissionType = "directory" // 目录
	PermissionTypeMenu      PermissionType = "menu"      // 菜单
	PermissionTypeOperation PermissionType = "operation" // 操作
)

type Permission struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"size:64;not null" json:"name"`
	Code       string         `gorm:"size:128;uniqueIndex;not null" json:"code"`
	Type       PermissionType `gorm:"size:20;not null;default:operation" json:"type"`
	ParentID   uint           `gorm:"default:0;index" json:"parent_id"`
	Path       string         `gorm:"size:128" json:"path"`        // 前端路由路径
	Icon       string         `gorm:"size:64" json:"icon"`         // 图标
	Sort       int            `gorm:"default:0" json:"sort"`       // 排序
	ApiPath    string         `gorm:"size:255" json:"api_path"`    // API路径
	ApiMethod  string         `gorm:"size:10" json:"api_method"`   // API方法
	Component  string         `gorm:"size:128" json:"component"`   // 前端组件路径
	Status     int8           `gorm:"default:1" json:"status"`     // 状态：1启用 0禁用
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Children   []Permission   `gorm:"-" json:"children,omitempty"` // 子权限（非数据库字段）
}

func (Permission) TableName() string {
	return "permissions"
}
```

- [ ] **Step 2: 验证编译**

```bash
cd /code/go/hello && go build ./model/...
```

Expected: 无错误输出

---

## Task 2: 更新Repository层

**Files:**
- Modify: `repository/permission.go`

- [ ] **Step 1: 添加树形查询方法**

在 `repository/permission.go` 添加以下方法：

```go
// GetAllWithChildren 获取所有权限（含子权限）
func (r *PermissionRepository) GetAllWithChildren() ([]model.Permission, error) {
	var permissions []model.Permission
	err := DB.Where("parent_id = 0").Order("sort asc").Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	
	// 递归获取子权限
	for i := range permissions {
		r.loadChildren(&permissions[i])
	}
	return permissions, nil
}

// loadChildren 递归加载子权限
func (r *PermissionRepository) loadChildren(perm *model.Permission) {
	var children []model.Permission
	DB.Where("parent_id = ?", perm.ID).Order("sort asc").Find(&children)
	for i := range children {
		r.loadChildren(&children[i])
	}
	perm.Children = children
}

// GetUserMenuIDs 获取用户的菜单ID列表
func (r *PermissionRepository) GetUserMenuIDs(userID uint) ([]uint, error) {
	var menuIDs []uint
	err := DB.Table("permissions").
		Select("DISTINCT permissions.id").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.type IN ?", userID, []model.PermissionType{model.PermissionTypeDirectory, model.PermissionTypeMenu}).
		Pluck("permissions.id", &menuIDs).Error
	return menuIDs, err
}

// GetUserMenus 获取用户有权限的菜单树
func (r *PermissionRepository) GetUserMenus(userID uint) ([]model.Permission, error) {
	// 获取用户所有菜单ID
	menuIDs, err := r.GetUserMenuIDs(userID)
	if err != nil {
		return nil, err
	}
	
	if len(menuIDs) == 0 {
		return []model.Permission{}, nil
	}
	
	// 获取所有目录和菜单
	var allMenus []model.Permission
	err = DB.Where("type IN ?", []model.PermissionType{model.PermissionTypeDirectory, model.PermissionTypeMenu}).
		Order("sort asc").Find(&allMenus).Error
	if err != nil {
		return nil, err
	}
	
	// 构建菜单树并过滤权限
	return r.buildUserMenuTree(allMenus, menuIDs, 0), nil
}

// buildUserMenuTree 构建用户菜单树
func (r *PermissionRepository) buildUserMenuTree(allMenus []model.Permission, menuIDs []uint, parentID uint) []model.Permission {
	var result []model.Permission
	menuIDSet := make(map[uint]bool)
	for _, id := range menuIDs {
		menuIDSet[id] = true
	}
	
	for _, menu := range allMenus {
		if menu.ParentID == parentID {
			// 检查是否有权限（或子菜单有权限）
			if r.hasMenuAccess(menu, allMenus, menuIDSet) {
				menu.Children = r.buildUserMenuTree(allMenus, menuIDs, menu.ID)
				result = append(result, menu)
			}
		}
	}
	return result
}

// hasMenuAccess 检查用户是否有菜单访问权限
func (r *PermissionRepository) hasMenuAccess(menu model.Permission, allMenus []model.Permission, menuIDSet map[uint]bool) bool {
	// 直接有权限
	if menuIDSet[menu.ID] {
		return true
	}
	// 检查子菜单是否有权限
	for _, m := range allMenus {
		if m.ParentID == menu.ID {
			if r.hasMenuAccess(m, allMenus, menuIDSet) {
				return true
			}
		}
	}
	return false
}

// GetByType 按类型查询
func (r *PermissionRepository) GetByType(permType model.PermissionType) ([]model.Permission, error) {
	var permissions []model.Permission
	err := DB.Where("type = ?", permType).Order("sort asc").Find(&permissions).Error
	return permissions, err
}
```

- [ ] **Step 2: 验证编译**

```bash
cd /code/go/hello && go build ./repository/...
```

---

## Task 3: 更新Service层

**Files:**
- Modify: `service/permission.go`
- Modify: `service/auth.go`

- [ ] **Step 1: 更新Permission Service**

在 `service/permission.go` 添加：

```go
// 在文件开头添加请求结构体
type CreatePermissionRequest struct {
	Name       string `json:"name" binding:"required"`
	Code       string `json:"code" binding:"required"`
	Type       string `json:"type" binding:"required,oneof=directory menu operation"`
	ParentID   uint   `json:"parent_id"`
	Path       string `json:"path"`
	Icon       string `json:"icon"`
	Sort       int    `json:"sort"`
	ApiPath    string `json:"api_path"`
	ApiMethod  string `json:"api_method"`
	Component  string `json:"component"`
	Status     *int8  `json:"status"`
}

type UpdatePermissionRequest struct {
	Name       string `json:"name" binding:"required"`
	Code       string `json:"code" binding:"required"`
	Type       string `json:"type" binding:"required,oneof=directory menu operation"`
	ParentID   uint   `json:"parent_id"`
	Path       string `json:"path"`
	Icon       string `json:"icon"`
	Sort       int    `json:"sort"`
	ApiPath    string `json:"api_path"`
	ApiMethod  string `json:"api_method"`
	Component  string `json:"component"`
	Status     *int8  `json:"status"`
}

// 更新Create方法
func (s *PermissionService) Create(req *CreatePermissionRequest) (*model.Permission, error) {
	exists, err := s.repo.ExistsByCode(req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("权限编码已存在")
	}

	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}

	permission := &model.Permission{
		Name:      req.Name,
		Code:      req.Code,
		Type:      model.PermissionType(req.Type),
		ParentID:  req.ParentID,
		Path:      req.Path,
		Icon:      req.Icon,
		Sort:      req.Sort,
		ApiPath:   req.ApiPath,
		ApiMethod: req.ApiMethod,
		Component: req.Component,
		Status:    status,
	}
	if err := s.repo.Create(permission); err != nil {
		return nil, err
	}
	return permission, nil
}

// 更新Update方法
func (s *PermissionService) Update(id uint, req *UpdatePermissionRequest) (*model.Permission, error) {
	permission, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("权限不存在")
	}

	if permission.Code != req.Code {
		exists, err := s.repo.ExistsByCode(req.Code)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("权限编码已存在")
		}
	}

	permission.Name = req.Name
	permission.Code = req.Code
	permission.Type = model.PermissionType(req.Type)
	permission.ParentID = req.ParentID
	permission.Path = req.Path
	permission.Icon = req.Icon
	permission.Sort = req.Sort
	permission.ApiPath = req.ApiPath
	permission.ApiMethod = req.ApiMethod
	permission.Component = req.Component
	if req.Status != nil {
		permission.Status = *req.Status
	}
	
	if err := s.repo.Update(permission); err != nil {
		return nil, err
	}
	return permission, nil
}

// GetTree 获取权限树
func (s *PermissionService) GetTree() ([]model.Permission, error) {
	return s.repo.GetAllWithChildren()
}
```

- [ ] **Step 2: 更新Auth Service**

在 `service/auth.go` 添加：

```go
// 在LoginResponse结构体后添加Menus字段
type LoginResponse struct {
	Token string          `json:"token"`
	User  *model.User     `json:"user"`
	Menus []model.Permission `json:"menus"`
}

// 更新Login方法
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if !user.CheckPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 提取角色编码
	var roleCodes []string
	for _, role := range user.Roles {
		roleCodes = append(roleCodes, role.Code)
	}

	// 生成Token
	token, err := jwt.GenerateToken(user.ID, user.Username, roleCodes)
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	// 获取用户菜单
	permRepo := repository.NewPermissionRepository()
	menus, err := permRepo.GetUserMenus(user.ID)
	if err != nil {
		menus = []model.Permission{}
	}

	// 清空密码
	user.Password = ""

	return &LoginResponse{
		Token: token,
		User:  user,
		Menus: menus,
	}, nil
}
```

- [ ] **Step 3: 验证编译**

```bash
cd /code/go/hello && go build ./service/...
```

---

## Task 4: 更新Handler层

**Files:**
- Modify: `handler/auth.go`
- Modify: `handler/permission.go`

- [ ] **Step 1: 更新Permission Handler**

在 `handler/permission.go` 添加获取权限树接口：

```go
// 在文件中添加新方法
// Tree 获取权限树
func (h *PermissionHandler) Tree(c *gin.Context) {
	tree, err := h.permissionService.GetTree()
	if err != nil {
		response.InternalError(c, "获取权限树失败")
		return
	}
	response.Success(c, tree)
}
```

- [ ] **Step 2: 在main.go中注册新路由**

在 `setupRoutes` 函数的 permissions 路由组中添加：

```go
permissions.GET("/tree", middleware.RequirePermission("permission:read"), permissionHandler.Tree)
```

- [ ] **Step 3: 验证编译**

```bash
cd /code/go/hello && go build .
```

---

## Task 5: 更新初始数据

**Files:**
- Modify: `repository/db.go`

- [ ] **Step 1: 更新InitData函数**

将 `repository/db.go` 中的 `InitData` 函数替换为：

```go
// InitData 初始化基础数据（超级管理员）
func InitData() error {
	// 检查是否已存在admin用户
	var count int64
	DB.Model(&model.User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		return nil
	}

	// 创建初始权限（树形结构）
	permissions := []model.Permission{
		// 目录
		{Name: "系统管理", Code: "system", Type: model.PermissionTypeDirectory, Icon: "Setting", Sort: 1},

		// 菜单
		{Name: "用户管理", Code: "system:user", Type: model.PermissionTypeMenu, ParentID: 1, Path: "/users", Icon: "User", Component: "user/List", Sort: 1},
		{Name: "角色管理", Code: "system:role", Type: model.PermissionTypeMenu, ParentID: 1, Path: "/roles", Icon: "Avatar", Component: "role/List", Sort: 2},
		{Name: "权限管理", Code: "system:permission", Type: model.PermissionTypeMenu, ParentID: 1, Path: "/permissions", Icon: "Key", Component: "permission/List", Sort: 3},

		// 用户操作
		{Name: "查看用户", Code: "system:user:read", Type: model.PermissionTypeOperation, ParentID: 2, ApiPath: "/api/users", ApiMethod: "GET", Sort: 1},
		{Name: "创建用户", Code: "system:user:create", Type: model.PermissionTypeOperation, ParentID: 2, ApiPath: "/api/users", ApiMethod: "POST", Sort: 2},
		{Name: "更新用户", Code: "system:user:update", Type: model.PermissionTypeOperation, ParentID: 2, ApiPath: "/api/users/:id", ApiMethod: "PUT", Sort: 3},
		{Name: "删除用户", Code: "system:user:delete", Type: model.PermissionTypeOperation, ParentID: 2, ApiPath: "/api/users/:id", ApiMethod: "DELETE", Sort: 4},

		// 角色操作
		{Name: "查看角色", Code: "system:role:read", Type: model.PermissionTypeOperation, ParentID: 3, ApiPath: "/api/roles", ApiMethod: "GET", Sort: 1},
		{Name: "创建角色", Code: "system:role:create", Type: model.PermissionTypeOperation, ParentID: 3, ApiPath: "/api/roles", ApiMethod: "POST", Sort: 2},
		{Name: "更新角色", Code: "system:role:update", Type: model.PermissionTypeOperation, ParentID: 3, ApiPath: "/api/roles/:id", ApiMethod: "PUT", Sort: 3},
		{Name: "删除角色", Code: "system:role:delete", Type: model.PermissionTypeOperation, ParentID: 3, ApiPath: "/api/roles/:id", ApiMethod: "DELETE", Sort: 4},

		// 权限操作
		{Name: "查看权限", Code: "system:permission:read", Type: model.PermissionTypeOperation, ParentID: 4, ApiPath: "/api/permissions", ApiMethod: "GET", Sort: 1},
		{Name: "创建权限", Code: "system:permission:create", Type: model.PermissionTypeOperation, ParentID: 4, ApiPath: "/api/permissions", ApiMethod: "POST", Sort: 2},
		{Name: "更新权限", Code: "system:permission:update", Type: model.PermissionTypeOperation, ParentID: 4, ApiPath: "/api/permissions/:id", ApiMethod: "PUT", Sort: 3},
		{Name: "删除权限", Code: "system:permission:delete", Type: model.PermissionTypeOperation, ParentID: 4, ApiPath: "/api/permissions/:id", ApiMethod: "DELETE", Sort: 4},
	}

	for i := range permissions {
		DB.FirstOrCreate(&permissions[i], model.Permission{Code: permissions[i].Code})
	}

	// 创建超级管理员角色（包含所有权限）
	adminRole := model.Role{
		Name:        "超级管理员",
		Code:        "admin",
		Permissions: permissions,
	}
	if err := DB.Create(&adminRole).Error; err != nil {
		return err
	}

	// 创建超级管理员用户
	adminUser := model.User{
		Username: "admin",
		Password: "admin123",
		Email:    "admin@example.com",
		Status:   1,
		Roles:    []model.Role{adminRole},
	}
	if err := DB.Create(&adminUser).Error; err != nil {
		return err
	}

	return nil
}
```

- [ ] **Step 2: 删除旧数据库重新初始化**

```bash
rm -f /code/go/hello/data/rbac.db
```

---

## Task 6: 前端权限指令

**Files:**
- Create: `web/src/directives/permission.js`

- [ ] **Step 1: 创建权限指令**

```javascript
// src/directives/permission.js
import { useUserStore } from '@/store/modules/user'

export const permission = {
  mounted(el, binding) {
    const userStore = useUserStore()
    const permissionCode = binding.value

    if (!permissionCode) return

    // 检查权限（支持数组）
    const codes = Array.isArray(permissionCode) ? permissionCode : [permissionCode]
    const hasPermission = codes.some(code => userStore.hasPermission(code))

    if (!hasPermission) {
      el.parentNode?.removeChild(el)
    }
  }
}

export function setupPermissionDirective(app) {
  app.directive('permission', permission)
}
```

- [ ] **Step 2: 在main.js中注册指令**

修改 `src/main.js`：

```javascript
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'

import App from './App.vue'
import router from './router'
import './styles/index.scss'
import { setupPermissionDirective } from './directives/permission'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ElementPlus, { locale: zhCn })
setupPermissionDirective(app)

app.mount('#app')
```

---

## Task 7: 前端状态管理更新

**Files:**
- Modify: `web/src/store/modules/user.js`

- [ ] **Step 1: 更新用户状态管理**

```javascript
import { defineStore } from 'pinia'
import { getToken, setToken, removeToken } from '@/utils/auth'
import { login as loginApi, logout as logoutApi } from '@/api/auth'
import router from '@/router'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: getToken() || '',
    user: null,
    permissions: [],
    menus: []
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
      this.menus = res.data.menus || []
      setToken(res.data.token)

      // 提取权限码
      if (this.user?.roles) {
        this.permissions = this.user.roles.flatMap(r => r.permissions?.map(p => p.code) || [])
      }

      // 生成动态路由
      this.generateRoutes()
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
      this.menus = []
      removeToken()
    },

    hasPermission(code) {
      if (this.isAdmin) return true
      return this.permissions.includes(code)
    },

    generateRoutes() {
      // 递归生成路由
      const generate = (menus) => {
        return menus
          .filter(m => m.type === 'menu' && m.component)
          .map(menu => ({
            path: menu.path,
            name: menu.code,
            component: () => import(`@/views/${menu.component}.vue`),
            meta: { title: menu.name, icon: menu.icon },
            children: menu.children ? generate(menu.children) : []
          }))
      }

      const routes = generate(this.menus)
      
      // 动态添加路由
      routes.forEach(route => {
        router.addRoute('main', route)
      })
    }
  }
})
```

---

## Task 8: 前端路由改造

**Files:**
- Modify: `web/src/router/index.js`

- [ ] **Step 1: 改造路由配置**

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
    name: 'main',
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
      }
      // 动态路由将通过 addRoute 添加
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

## Task 9: 前端侧边栏改造

**Files:**
- Modify: `web/src/components/layout/Sidebar.vue`

- [ ] **Step 1: 改造侧边栏组件**

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
      
      <!-- 动态菜单 -->
      <template v-for="menu in userStore.menus" :key="menu.id">
        <el-sub-menu v-if="menu.type === 'directory' && menu.children?.length" :index="menu.code">
          <template #title>
            <el-icon><component :is="menu.icon" /></el-icon>
            <span>{{ menu.name }}</span>
          </template>
          <el-menu-item v-for="child in menu.children" :key="child.id" :index="child.path">
            <el-icon><component :is="child.icon" /></el-icon>
            <span>{{ child.name }}</span>
          </el-menu-item>
        </el-sub-menu>
        
        <el-menu-item v-else-if="menu.type === 'menu'" :index="menu.path">
          <el-icon><component :is="menu.icon" /></el-icon>
          <span>{{ menu.name }}</span>
        </el-menu-item>
      </template>
    </el-menu>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { House } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/modules/user'

const route = useRoute()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)
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

---

## Task 10: 权限管理页面改造

**Files:**
- Modify: `web/src/views/permission/List.vue`

- [ ] **Step 1: 改造权限管理页面为树形展示**

将 `src/views/permission/List.vue` 改造为支持树形结构的权限管理页面（支持创建目录、菜单、操作三种类型）。

主要改动：
1. 添加权限类型选择
2. 添加父级权限选择（下拉树）
3. 根据类型动态显示不同表单字段
4. 表格改为树形表格展示

---

## Task 11: 重新编译测试

- [ ] **Step 1: 删除旧数据库**

```bash
rm -f /code/go/hello/data/rbac.db
```

- [ ] **Step 2: 重新编译后端**

```bash
cd /code/go/hello && go build -o rbac-server .
```

- [ ] **Step 3: 启动服务测试**

```bash
# 启动后端
JWT_SECRET=rbac-secret-key-2024 ./rbac-server &

# 启动前端
cd /code/go/hello/web && npm run dev &
```

- [ ] **Step 4: 验证功能**

1. 访问 http://localhost:5173
2. 使用 admin/admin123 登录
3. 验证菜单树是否正确显示
4. 验证权限管理页面是否支持树形结构
