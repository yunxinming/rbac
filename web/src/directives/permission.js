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
