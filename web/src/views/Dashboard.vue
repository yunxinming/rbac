<template>
  <div class="dashboard">
    <h2>欢迎回来，{{ userStore.user?.username }}</h2>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="8" v-if="userStore.hasPermission('system:user:read')">
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
      <el-col :span="8" v-if="userStore.hasPermission('system:role:read')">
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
      <el-col :span="8" v-if="userStore.hasPermission('system:permission:read')">
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
  const promises = []

  if (userStore.hasPermission('system:user:read')) {
    promises.push(
      getUsers().then(res => { stats.value.users = res.data?.length || 0 }).catch(() => {})
    )
  }

  if (userStore.hasPermission('system:role:read')) {
    promises.push(
      getRoles().then(res => { stats.value.roles = res.data?.length || 0 }).catch(() => {})
    )
  }

  if (userStore.hasPermission('system:permission:read')) {
    promises.push(
      getPermissions().then(res => { stats.value.permissions = res.data?.length || 0 }).catch(() => {})
    )
  }

  if (promises.length > 0) {
    await Promise.all(promises)
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
