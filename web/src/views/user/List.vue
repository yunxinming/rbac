<template>
  <div class="user-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户列表</span>
          <el-button type="primary" @click="openDialog()" v-if="canCreate">
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
        <el-table-column v-if="showActions" label="操作" :width="actionWidth" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDialog(row)" v-if="canUpdate">
              编辑
            </el-button>
            <el-button link type="danger" @click="handleDelete(row)" v-if="canDelete">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

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
        <el-form-item label="角色" prop="role_ids" v-if="canReadRoles">
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
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/modules/user'
import { getUsers, createUser, updateUser, deleteUser, assignRoles } from '@/api/user'
import { getRoles } from '@/api/role'
import { calculateActionWidth, confirmDelete } from '@/utils/table'

const userStore = useUserStore()

const loading = ref(false)
const submitting = ref(false)
const users = ref([])
const roles = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const formRef = ref(null)

// 缓存权限检查结果
const canCreate = computed(() => userStore.hasPermission('system:user:create'))
const canUpdate = computed(() => userStore.hasPermission('system:user:update'))
const canDelete = computed(() => userStore.hasPermission('system:user:delete'))
const canReadRoles = computed(() => userStore.hasPermission('system:role:read'))

// 是否显示操作栏
const showActions = computed(() => canUpdate.value || canDelete.value)

// 动态计算操作栏宽度
const actionWidth = computed(() => {
  const buttons = []
  if (canUpdate.value) buttons.push('编辑')
  if (canDelete.value) buttons.push('删除')
  return calculateActionWidth(buttons)
})

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
    // 并行获取用户列表和角色列表
    const promises = [getUsers()]
    if (canReadRoles.value) {
      promises.push(getRoles())
    }
    const results = await Promise.all(promises)
    users.value = results[0].data || []
    if (results[1]) {
      roles.value = results[1].data || []
    }
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
  await confirmDelete('确定要删除该用户吗？', async () => {
    await deleteUser(row.id)
    fetchData()
  })
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
