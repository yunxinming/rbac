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
            <el-icon><component :is="getIcon(menu.icon)" /></el-icon>
            <span>{{ menu.name }}</span>
          </template>
          <el-menu-item v-for="child in menu.children" :key="child.id" :index="child.path">
            <el-icon><component :is="getIcon(child.icon)" /></el-icon>
            <span>{{ child.name }}</span>
          </el-menu-item>
        </el-sub-menu>

        <el-menu-item v-else-if="menu.type === 'menu'" :index="menu.path">
          <el-icon><component :is="getIcon(menu.icon)" /></el-icon>
          <span>{{ menu.name }}</span>
        </el-menu-item>
      </template>
    </el-menu>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { House, Setting, User, Avatar, Key } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/modules/user'

const route = useRoute()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

// 图标映射
const iconMap = {
  Setting,
  User,
  Avatar,
  Key,
  House
}

const getIcon = (iconName) => {
  return iconMap[iconName] || Setting
}
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
