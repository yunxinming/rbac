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
  },
  // 404 页面
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/dashboard'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 标记动态路由是否已注册
let dynamicRoutesAdded = false

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const token = getToken()

  if (to.meta.requiresAuth === false) {
    next()
    return
  }

  if (!token) {
    next('/login')
    return
  }

  // 确保动态路由已注册
  if (!dynamicRoutesAdded) {
    const userStore = useUserStore()
    userStore.initRoutes()
    dynamicRoutesAdded = true
    // 重新导航到目标路由
    next({ ...to, replace: true })
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
