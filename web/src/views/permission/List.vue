<template>
  <div class="permission-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>权限管理</span>
          <el-button type="primary" @click="openDialog()" v-if="canCreate">
            新增权限
          </el-button>
        </div>
      </template>
      <el-table
        :data="permissions"
        v-loading="loading"
        stripe
        row-key="id"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        default-expand-all
      >
        <el-table-column prop="name" label="权限名称" width="200" />
        <el-table-column prop="code" label="权限编码" width="200" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.type === 'directory'" type="warning">目录</el-tag>
            <el-tag v-else-if="row.type === 'menu'" type="success">菜单</el-tag>
            <el-tag v-else type="info">操作</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="icon" label="图标" width="100" />
        <el-table-column prop="path" label="路由路径" width="150" />
        <el-table-column prop="component" label="组件" width="150" />
        <el-table-column prop="api_path" label="API路径" show-overflow-tooltip />
        <el-table-column prop="api_method" label="API方法" width="100" />
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑权限' : '新增权限'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="权限类型" prop="type">
          <el-radio-group v-model="form.type" :disabled="isEdit">
            <el-radio value="directory">目录</el-radio>
            <el-radio value="menu">菜单</el-radio>
            <el-radio value="operation">操作</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="父级权限" prop="parent_id">
          <el-tree-select
            v-model="form.parent_id"
            :data="parentOptions"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            check-strictly
            clearable
            placeholder="请选择父级权限"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="权限名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="权限编码" prop="code">
          <el-input v-model="form.code" placeholder="如：system:user:create" />
        </el-form-item>

        <!-- 目录/菜单字段 -->
        <template v-if="form.type === 'directory' || form.type === 'menu'">
          <el-form-item label="图标" prop="icon">
            <el-input v-model="form.icon" placeholder="如：User" />
          </el-form-item>
          <el-form-item label="排序" prop="sort">
            <el-input-number v-model="form.sort" :min="0" />
          </el-form-item>
        </template>

        <!-- 菜单专属字段 -->
        <template v-if="form.type === 'menu'">
          <el-form-item label="路由路径" prop="path">
            <el-input v-model="form.path" placeholder="如：/users" />
          </el-form-item>
          <el-form-item label="组件路径" prop="component">
            <el-input v-model="form.component" placeholder="如：user/List" />
          </el-form-item>
        </template>

        <!-- 操作专属字段 -->
        <template v-if="form.type === 'operation'">
          <el-form-item label="API路径" prop="api_path">
            <el-input v-model="form.api_path" placeholder="如：/api/users" />
          </el-form-item>
          <el-form-item label="API方法" prop="api_method">
            <el-select v-model="form.api_method" placeholder="请选择">
              <el-option label="GET" value="GET" />
              <el-option label="POST" value="POST" />
              <el-option label="PUT" value="PUT" />
              <el-option label="DELETE" value="DELETE" />
            </el-select>
          </el-form-item>
          <el-form-item label="排序" prop="sort">
            <el-input-number v-model="form.sort" :min="0" />
          </el-form-item>
        </template>
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
import { getPermissionTree, createPermission, updatePermission, deletePermission } from '@/api/permission'
import { calculateActionWidth, confirmDelete } from '@/utils/table'

const userStore = useUserStore()

const loading = ref(false)
const submitting = ref(false)
const permissions = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const formRef = ref(null)

// 缓存权限检查结果
const canCreate = computed(() => userStore.hasPermission('system:permission:create'))
const canUpdate = computed(() => userStore.hasPermission('system:permission:update'))
const canDelete = computed(() => userStore.hasPermission('system:permission:delete'))

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
  name: '',
  code: '',
  type: 'operation',
  parent_id: 0,
  path: '',
  icon: '',
  sort: 0,
  api_path: '',
  api_method: '',
  component: ''
})

const rules = {
  name: [{ required: true, message: '请输入权限名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入权限编码', trigger: 'blur' }],
  type: [{ required: true, message: '请选择权限类型', trigger: 'change' }]
}

// 父级权限选项
const parentOptions = computed(() => {
  const options = [{ id: 0, name: '顶级权限', children: [] }]
  // 只允许选择目录和菜单作为父级
  const filterMenus = (items) => {
    return items
      .filter(item => item.type === 'directory' || item.type === 'menu')
      .map(item => ({
        ...item,
        children: item.children ? filterMenus(item.children) : []
      }))
  }
  options[0].children = filterMenus(permissions.value)
  return options
})

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getPermissionTree()
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
    form.type = row.type
    form.parent_id = row.parent_id || 0
    form.path = row.path || ''
    form.icon = row.icon || ''
    form.sort = row.sort || 0
    form.api_path = row.api_path || ''
    form.api_method = row.api_method || ''
    form.component = row.component || ''
  } else {
    form.name = ''
    form.code = ''
    form.type = 'operation'
    form.parent_id = 0
    form.path = ''
    form.icon = ''
    form.sort = 0
    form.api_path = ''
    form.api_method = ''
    form.component = ''
  }

  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const data = { ...form }
    // 根据类型清理不需要的字段
    if (data.type === 'directory') {
      data.path = ''
      data.component = ''
      data.api_path = ''
      data.api_method = ''
    } else if (data.type === 'menu') {
      data.api_path = ''
      data.api_method = ''
    } else {
      data.path = ''
      data.component = ''
      data.icon = ''
    }

    if (isEdit.value) {
      await updatePermission(currentId.value, data)
      ElMessage.success('更新成功')
    } else {
      await createPermission(data)
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
  if (row.children && row.children.length > 0) {
    ElMessage.warning('请先删除子权限')
    return
  }
  await confirmDelete('确定要删除该权限吗？', async () => {
    await deletePermission(row.id)
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
