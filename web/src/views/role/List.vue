<template>
  <div class="role-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>角色列表</span>
          <el-button type="primary" @click="openDialog()" v-if="canCreate">
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
        <el-table-column v-if="showActions" label="操作" :width="actionWidth" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDialog(row)" v-if="canUpdate">
              编辑
            </el-button>
            <el-button link type="primary" @click="openPermissionDialog(row)" v-if="canAssignPermissions">
              分配权限
            </el-button>
            <el-button link type="danger" @click="handleDelete(row)" v-if="canDelete">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

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

    <el-dialog v-model="permissionDialogVisible" title="分配权限" width="500px">
      <el-tree
        ref="permissionTreeRef"
        :data="permissionTree"
        :props="{ label: 'name', children: 'children' }"
        show-checkbox
        node-key="id"
        default-expand-all
      />
      <template #footer>
        <el-button @click="permissionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAssignPermissions" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/modules/user'
import { getRoles, createRole, updateRole, deleteRole, assignPermissions } from '@/api/role'
import { getPermissionTree } from '@/api/permission'
import { calculateActionWidth, confirmDelete } from '@/utils/table'

const userStore = useUserStore()

const loading = ref(false)
const submitting = ref(false)
const roles = ref([])
const permissionTree = ref([])
const dialogVisible = ref(false)
const permissionDialogVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const formRef = ref(null)
const permissionTreeRef = ref(null)

// 缓存权限检查结果
const canCreate = computed(() => userStore.hasPermission('system:role:create'))
const canUpdate = computed(() => userStore.hasPermission('system:role:update'))
const canDelete = computed(() => userStore.hasPermission('system:role:delete'))
const canReadPermissions = computed(() => userStore.hasPermission('system:permission:read'))

// 是否可以分配权限
const canAssignPermissions = computed(() => canUpdate.value && canReadPermissions.value)

// 是否显示操作栏
const showActions = computed(() => canUpdate.value || canAssignPermissions.value || canDelete.value)

// 动态计算操作栏宽度
const actionWidth = computed(() => {
  const buttons = []
  if (canUpdate.value) buttons.push('编辑')
  if (canAssignPermissions.value) buttons.push('分配权限')
  if (canDelete.value) buttons.push('删除')
  return calculateActionWidth(buttons)
})

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
    // 并行获取角色列表和权限树
    const promises = [getRoles()]
    if (canReadPermissions.value) {
      promises.push(getPermissionTree())
    }
    const results = await Promise.all(promises)
    roles.value = results[0].data || []
    if (results[1]) {
      permissionTree.value = results[1].data || []
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
    form.name = row.name
    form.code = row.code
  } else {
    form.name = ''
    form.code = ''
  }

  dialogVisible.value = true
}

const openPermissionDialog = async (row) => {
  currentId.value = row.id
  permissionDialogVisible.value = true
  // 使用 nextTick 等待对话框渲染后设置选中状态
  await nextTick()
  const checkedIds = row.permissions?.map(p => p.id) || []
  permissionTreeRef.value?.setCheckedKeys(checkedIds)
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
    const checkedIds = permissionTreeRef.value?.getCheckedKeys() || []
    await assignPermissions(currentId.value, checkedIds)
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
  await confirmDelete('确定要删除该角色吗？', async () => {
    await deleteRole(row.id)
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
