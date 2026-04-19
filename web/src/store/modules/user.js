import { defineStore } from 'pinia'
import { getToken, setToken, removeToken } from '@/utils/auth'
import { login as loginApi, logout as logoutApi } from '@/api/auth'
import router from '@/router'

const MENUS_KEY = 'rbac_menus'
const USER_KEY = 'rbac_user'
const PERMISSIONS_KEY = 'rbac_permissions'

// 组件映射（Vite 动态 import 需要静态分析）
const componentModules = import.meta.glob('@/views/**/*.vue')

export const useUserStore = defineStore('user', {
  state: () => ({
    token: getToken() || '',
    user: JSON.parse(localStorage.getItem(USER_KEY) || 'null'),
    permissions: JSON.parse(localStorage.getItem(PERMISSIONS_KEY) || '[]'),
    menus: JSON.parse(localStorage.getItem(MENUS_KEY) || '[]')
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

      // 保存到 localStorage
      localStorage.setItem(USER_KEY, JSON.stringify(this.user))
      localStorage.setItem(MENUS_KEY, JSON.stringify(this.menus))
      localStorage.setItem(PERMISSIONS_KEY, JSON.stringify(this.permissions))

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
      localStorage.removeItem(USER_KEY)
      localStorage.removeItem(MENUS_KEY)
      localStorage.removeItem(PERMISSIONS_KEY)
    },

    hasPermission(code) {
      if (this.isAdmin) return true
      return this.permissions.includes(code)
    },

    // 初始化时恢复动态路由
    initRoutes() {
      if (this.menus.length > 0) {
        this.generateRoutes()
      }
    },

    generateRoutes() {
      const routes = []

      // 递归提取所有菜单类型的路由
      const extractRoutes = (menus) => {
        menus.forEach(menu => {
          if (menu.type === 'menu' && menu.component) {
            // 去掉路径前导斜杠，作为子路由
            let path = menu.path
            if (path.startsWith('/')) {
              path = path.slice(1)
            }
            routes.push({
              path: path,
              name: menu.code,
              component: componentModules[`/src/views/${menu.component}.vue`],
              meta: { title: menu.name, icon: menu.icon }
            })
          }
          // 递归处理子菜单
          if (menu.children && menu.children.length > 0) {
            extractRoutes(menu.children)
          }
        })
      }

      extractRoutes(this.menus)
      console.log('Generated routes:', routes)

      // 动态添加路由
      routes.forEach(route => {
        if (!router.hasRoute(route.name)) {
          router.addRoute('main', route)
          console.log('Added route:', route.path, route.name)
        }
      })

      console.log('All routes:', router.getRoutes().map(r => r.path))
    }
  }
})
