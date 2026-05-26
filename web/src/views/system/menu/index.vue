<template>
  <div class="page-card">
    <div class="mb-4">
      <el-button type="primary" @click="openCreate()">新增根菜单</el-button>
    </div>
    <el-table v-loading="loading" :data="tree" row-key="id" default-expand-all>
      <el-table-column prop="title" label="标题" min-width="120" />
      <el-table-column prop="menuType" label="类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="menuTypeTag[row.menuType] || 'info'" size="small">
            {{ menuTypeLabel[row.menuType] || '-' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="图标" width="80" align="center">
        <template #default="{ row }">
          <el-icon v-if="row.menuType !== 2 && row.icon" :size="18">
            <component :is="resolveIcon(row.icon)" />
          </el-icon>
          <span v-else class="text-gray-400">-</span>
        </template>
      </el-table-column>

      <el-table-column prop="path" label="路径" min-width="120" show-overflow-tooltip />
      <el-table-column prop="name" label="路由名" min-width="120" show-overflow-tooltip />
      <el-table-column prop="component" label="组件" min-width="120" show-overflow-tooltip />
      <el-table-column prop="permCode" label="权限码" min-width="120" show-overflow-tooltip />
      <el-table-column prop="sort" label="排序" width="80" align="center" />
      <el-table-column label="隐藏" width="80" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.menuType !== 2" :type="row.hidden ? 'warning' : 'info'" size="small">
            {{ row.hidden ? '是' : '否' }}
          </el-tag>
          <span v-else class="text-gray-400">-</span>
        </template>
      </el-table-column>
      <el-table-column label="缓存" width="80" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.menuType === 1" :type="row.keepAlive ? 'success' : 'info'" size="small">
            {{ row.keepAlive ? '是' : '否' }}
          </el-tag>
          <span v-else class="text-gray-400">-</span>
        </template>
      </el-table-column>
      <el-table-column label="启用" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.enabled ? 'success' : 'danger'" size="small">
            {{ row.enabled ? '是' : '否' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="170">
        <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.menuType !== 2" link type="primary" @click="openCreate(row)">子项</el-button>
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button link type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑菜单' : '新增菜单'" width="600px" destroy-on-close>
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        :validate-on-rule-change="false"
        label-width="90px"
      >
        <el-form-item label="类型">
          <el-radio-group v-if="!editingId" v-model="form.menuType">
            <el-radio v-for="t in menuTypes" :key="t.value" :value="t.value">{{ t.label }}</el-radio>
          </el-radio-group>
          <el-tag v-else :type="menuTypeTag[form.menuType] || 'info'" size="small">
            {{ menuTypeLabel[form.menuType] || '-' }}
          </el-tag>
        </el-form-item>
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="侧栏与面包屑显示名称" />
        </el-form-item>
        <el-form-item v-if="form.menuType === 2" label="权限码" prop="permCode">
          <el-input v-model="form.permCode" placeholder="如 system:user:add" />
        </el-form-item>
        <el-form-item v-if="form.menuType !== 2" label="路径" prop="path">
          <el-input v-model="form.path" :placeholder="pathPlaceholder" />
        </el-form-item>
        <el-form-item v-if="form.menuType !== 2" label="路由名" prop="name">
          <el-input v-model="form.name" placeholder="Vue 路由唯一名称，如 SystemUser（用于 Tab/缓存，不是 URL）" />
        </el-form-item>
        <el-form-item v-if="form.menuType !== 2" label="组件" prop="component">
          <el-input v-model="form.component" :placeholder="componentPlaceholder" />
        </el-form-item>
        <el-form-item v-if="form.menuType !== 2" label="图标">
          <IconPicker v-model="form.icon" />
        </el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" /></el-form-item>
        <el-form-item v-if="form.menuType !== 2" label="隐藏">
          <el-switch v-model="form.hidden" />
          <span class="text-muted ml-2">隐藏后不出现在侧栏，仍可通过地址访问</span>
        </el-form-item>
        <el-form-item v-if="form.menuType === 1" label="页面缓存">
          <el-switch v-model="form.keepAlive" />
          <span class="text-muted ml-2">开启后切换 Tab 保留页面状态</span>
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
          <span class="text-muted ml-2">禁用后不可访问且不参与权限分配</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
defineOptions({ name: 'SystemMenu' })

import { computed, nextTick, onMounted, ref, watch } from 'vue'
import * as Icons from '@element-plus/icons-vue'
import { menuApi } from '@/api/system'
import { ElMessage, ElMessageBox } from 'element-plus'
import IconPicker from '@/components/IconPicker.vue'

function resolveIcon(name) {
  return Icons[name] || Icons.Menu
}

const loading = ref(false)
const tree = ref([])
const dialogVisible = ref(false)
const editingId = ref(null)
const formRef = ref()

const menuTypes = [
  { label: '目录', value: 0 },
  { label: '页面', value: 1 },
  { label: '按钮', value: 2 },
]

const menuTypeLabel = { 0: '目录', 1: '页面', 2: '按钮' }
const menuTypeTag = { 0: 'info', 1: 'primary', 2: 'warning' }

const form = ref({
  parentId: null, menuType: 1, permCode: '', path: '', name: '', component: '',
  title: '', icon: '', sort: 0, hidden: false, keepAlive: false, enabled: true,
})

const pathPlaceholder = computed(() => {
  if (form.value.menuType === 0) {
    return '浏览器地址路径，根目录用 /system；仅分组，本身不打开页面'
  }
  return form.value.parentId
    ? '相对父级路径，如 user → 实际访问 /system/user'
    : '顶级页面用绝对路径，如 /dashboard'
})

const componentPlaceholder = computed(() => {
  if (form.value.menuType === 0) {
    return '目录填 Layout，加载子路由容器'
  }
  return '对应 views 下文件，不含 .vue，如 system/user/index'
})

const required = (message) => [{ required: true, message, trigger: [] }]

const formRules = computed(() => {
  const rules = {
    title: required('请输入标题'),
  }
  if (form.value.menuType === 2) {
    rules.permCode = required('请输入权限码')
    return rules
  }
  rules.path = required('请输入路径')
  rules.name = required('请输入路由名')
  rules.component = required('请输入组件路径')
  return rules
})

function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString()
}

async function load() {
  loading.value = true
  try {
    tree.value = await menuApi.tree()
  } finally {
    loading.value = false
  }
}

function openCreate(parent) {
  editingId.value = null
  form.value = {
    parentId: parent?.id || null,
    menuType: 1,
    permCode: '',
    path: '',
    name: '',
    component: '',
    title: '',
    icon: 'Menu',
    sort: 0,
    hidden: false,
    keepAlive: false,
    enabled: true,
  }
  dialogVisible.value = true
}

function openEdit(row) {
  editingId.value = row.id
  form.value = { ...row, parentId: row.parentId }
  dialogVisible.value = true
}

watch(dialogVisible, (visible) => {
  if (visible) {
    nextTick(() => formRef.value?.clearValidate())
  }
})

async function submit() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  if (editingId.value) await menuApi.update(editingId.value, form.value)
  else await menuApi.create(form.value)
  ElMessage.success('保存成功，路由与侧栏变更需重新登录后生效')
  dialogVisible.value = false
  load()
}

async function remove(row) {
  await ElMessageBox.confirm(`确定删除「${row.title}」？`, '提示', { type: 'warning' })
  await menuApi.remove(row.id)
  ElMessage.success('已删除')
  load()
}

onMounted(load)
</script>

