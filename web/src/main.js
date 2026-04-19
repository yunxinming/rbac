import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'

import App from './App.vue'
import router from './router'
import './styles/index.scss'
import { setupPermissionDirective } from './directives/permission'
import { useUserStore } from './store/modules/user'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(ElementPlus, { locale: zhCn })
setupPermissionDirective(app)

// 初始化动态路由
const userStore = useUserStore()
userStore.initRoutes()

app.mount('#app')
