<template>
  <div class="page-layout">
    <div class="page-card">
      <div class="flex flex-wrap gap-3 items-center">
        <el-input v-model="query.keyword" placeholder="名称/编码" clearable class="w-48!" @keyup.enter="load" />
        <el-button type="primary" @click="load">查询</el-button>
        <el-button @click="resetQuery">重置</el-button>
      </div>
    </div>

    <div class="page-card">
      <div class="mb-3">
        <el-button v-auth="'system:role:add'" type="primary" @click="openCreate">新增角色</el-button>
      </div>

      <el-table v-loading="loading" :data="list" stripe>
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="code" label="编码" />
        <el-table-column label="数据范围" width="130">
          <template #default="{ row }">
            {{ dataScopeOptions.find((o) => o.value === row.dataScope)?.label }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">{{ row.enabled ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button v-auth="'system:role:edit'" link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button v-auth="'system:role:edit'" link type="primary" @click="openMenus(row)">菜单</el-button>
            <el-button v-auth="'system:role:delete'" link type="danger" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="flex justify-end mt-4">
        <AppPagination v-model:page="query.page" v-model:page-size="query.pageSize" :total="total" @change="load" />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑角色' : '新增角色'" width="480px">
      <el-form label-width="90px">
        <el-form-item label="名称" required><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="编码" required><el-input v-model="form.code" :disabled="!!editingId" /></el-form-item>
        <el-form-item label="数据范围">
          <el-select v-model="form.dataScope" class="w-full">
            <el-option v-for="o in dataScopeOptions" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" /></el-form-item>
        <el-form-item label="状态"><el-switch v-model="form.enabled" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="menuDialogVisible" title="分配菜单与按钮权限" width="520px" destroy-on-close>
      <p class="text-muted text-sm mt-0 mb-3">
        树中包含目录、页面及页面下的按钮；父子节点联动，子节点未全选时父节点为半选状态。
      </p>
      <el-scrollbar max-height="420px">
        <el-tree
          ref="menuTreeRef"
          :data="menuTree"
          node-key="id"
          show-checkbox
          default-expand-all
          :props="{ label: 'title', children: 'children' }"
        >
          <template #default="{ data }">
            <span class="inline-flex items-center gap-2 flex-wrap">
              <el-tag :type="menuTypeTag[data.menuType] || 'info'" size="small">
                {{ menuTypeLabel[data.menuType] || '菜单' }}
              </el-tag>
              <span>{{ data.title }}</span>
              <span v-if="data.menuType === 2 && data.permCode" class="text-xs text-gray-400 font-mono">
                {{ data.permCode }}
              </span>
            </span>
          </template>
        </el-tree>
      </el-scrollbar>
      <template #footer>
        <el-button @click="menuDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveMenus">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { nextTick, onMounted, reactive, ref, watch } from 'vue'
import AppPagination from '@/components/AppPagination.vue'
import { roleApi, menuApi } from '@/api/system'
import { filterEnabledMenuTree } from '@/utils/route'
import { ElMessage, ElMessageBox } from 'element-plus'

const dataScopeOptions = [
  { label: '全部', value: 1 },
  { label: '自定义部门', value: 2 },
  { label: '本部门', value: 3 },
  { label: '本部门及下级', value: 4 },
  { label: '仅本人', value: 5 },
]

const menuTypeLabel = { 0: '目录', 1: '页面', 2: '按钮' }
const menuTypeTag = { 0: 'info', 1: 'primary', 2: 'warning' }

const loading = ref(false)
const list = ref([])
const total = ref(0)
const menuTree = ref([])
const dialogVisible = ref(false)
const menuDialogVisible = ref(false)
const editingId = ref(null)
const checkedMenuIds = ref([])
const menuTreeRef = ref()

const query = reactive({ page: 1, pageSize: 10, keyword: '' })
const form = reactive({
  name: '', code: '', sort: 0, enabled: true, dataScope: 5, description: '', deptIds: [],
})

async function load() {
  loading.value = true
  try {
    const data = await roleApi.list(query)
    list.value = data.list
    total.value = data.total
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.keyword = ''
  query.page = 1
  load()
}

function openCreate() {
  editingId.value = null
  Object.assign(form, { name: '', code: '', sort: 0, enabled: true, dataScope: 5, description: '', deptIds: [] })
  dialogVisible.value = true
}

function openEdit(row) {
  editingId.value = row.id
  Object.assign(form, {
    name: row.name, code: row.code, sort: row.sort, enabled: row.enabled,
    dataScope: row.dataScope, description: row.description, deptIds: row.deptIds || [],
  })
  dialogVisible.value = true
}

async function submit() {
  if (editingId.value) await roleApi.update(editingId.value, { ...form })
  else await roleApi.create({ ...form })
  ElMessage.success('保存成功')
  dialogVisible.value = false
  load()
}

/** 收集勾选与半选节点（半选表示子节点未全选） */
function collectCheckedIds() {
  const tree = menuTreeRef.value
  if (!tree) return [...checkedMenuIds.value]
  return [...new Set([...tree.getCheckedKeys(), ...tree.getHalfCheckedKeys()])]
}

/** 回显时只设置叶子节点，由树组件计算父级勾选/半选 */
function toTreeLeafKeys(menuIds, nodes) {
  const set = new Set(menuIds)
  const keys = []
  const walk = (list) => {
    for (const node of list || []) {
      const kids = node.children || []
      if (!kids.length) {
        if (set.has(node.id)) keys.push(node.id)
        continue
      }
      const pageWithButtons = node.menuType === 1 && kids.every((c) => c.menuType === 2)
      if (pageWithButtons) {
        kids.forEach((c) => {
          if (set.has(c.id)) keys.push(c.id)
        })
      } else {
        walk(kids)
      }
    }
  }
  walk(nodes)
  return keys
}

async function openMenus(row) {
  editingId.value = row.id
  menuTree.value = filterEnabledMenuTree(await menuApi.tree())
  const res = await roleApi.getMenus(row.id)
  checkedMenuIds.value = res.menuIds || []
  menuDialogVisible.value = true
}

watch(menuDialogVisible, async (visible) => {
  if (!visible) return
  await nextTick()
  const leafKeys = toTreeLeafKeys(checkedMenuIds.value, menuTree.value)
  menuTreeRef.value?.setCheckedKeys(leafKeys, true)
})

async function saveMenus() {
  const ids = collectCheckedIds()
  await roleApi.assignMenus(editingId.value, ids)
  ElMessage.success('菜单与按钮权限已更新')
  menuDialogVisible.value = false
}

async function remove(row) {
  await ElMessageBox.confirm(`确定删除角色「${row.name}」？`, '提示', { type: 'warning' })
  await roleApi.remove(row.id)
  ElMessage.success('已删除')
  load()
}

onMounted(load)
</script>
