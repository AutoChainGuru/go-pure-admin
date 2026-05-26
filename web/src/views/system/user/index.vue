<template>
  <div class="page-layout">
    <div class="page-card">
      <div class="flex flex-wrap gap-3 items-center">
        <el-input v-model="query.username" placeholder="用户名" clearable class="w-48!" @keyup.enter="load" />
        <el-button type="primary" @click="load">查询</el-button>
        <el-button @click="resetQuery">重置</el-button>
      </div>
    </div>

    <div class="page-card">
      <div class="mb-3">
        <el-button v-auth="'system:user:add'" type="primary" @click="openCreate">新增用户</el-button>
      </div>

      <el-table v-loading="loading" :data="list" stripe>
      <el-table-column prop="username" label="用户名" min-width="120" />
      <el-table-column prop="nickname" label="昵称" min-width="120" />
      <el-table-column prop="dept.name" label="部门" min-width="120" />
      <el-table-column label="角色" min-width="160">
        <template #default="{ row }">
          <el-tag v-for="r in row.roles || []" :key="r.id" size="small" class="mr-1">{{ r.name }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="row.enabled ? 'success' : 'info'" size="small">{{ row.enabled ? '启用' : '禁用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button v-auth="'system:user:edit'" link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button v-auth="'system:user:delete'" link type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

      <div class="flex justify-end mt-4">
        <AppPagination v-model:page="query.page" v-model:page-size="query.pageSize" :total="total" @change="load" />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑用户' : '新增用户'" width="480px" destroy-on-close>
      <el-form label-width="80px">
        <el-form-item label="用户名" required>
          <el-input v-model="form.username" :disabled="!!editingId" />
        </el-form-item>
        <el-form-item :label="editingId ? '新密码' : '密码'" :required="!editingId">
          <el-input v-model="form.password" type="password" show-password placeholder="留空则不修改" />
        </el-form-item>
        <el-form-item label="昵称"><el-input v-model="form.nickname" /></el-form-item>
        <el-form-item label="部门">
          <el-tree-select
            v-model="form.deptId"
            :data="deptTree"
            node-key="id"
            :props="{ label: 'name', children: 'children' }"
            check-strictly
            clearable
            class="w-full"
          />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.roleIds" multiple class="w-full">
            <el-option v-for="r in roles" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态"><el-switch v-model="form.enabled" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
defineOptions({ name: 'SystemUser' })

import { onMounted, reactive, ref } from 'vue'
import AppPagination from '@/components/AppPagination.vue'
import { userApi, roleApi, deptApi } from '@/api/system'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const list = ref([])
const total = ref(0)
const roles = ref([])
const deptTree = ref([])
const dialogVisible = ref(false)
const editingId = ref(null)

const query = reactive({ page: 1, pageSize: 10, username: '' })
const form = reactive({
  username: '',
  password: '',
  nickname: '',
  phone: '',
  email: '',
  enabled: true,
  deptId: null,
  roleIds: [],
})

async function load() {
  loading.value = true
  try {
    const data = await userApi.list(query)
    list.value = data.list
    total.value = data.total
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.username = ''
  query.page = 1
  load()
}

let optionsLoaded = false

/** 仅打开表单时需要角色/部门下拉，避免列表页加载未授权接口 */
async function ensureFormOptions() {
  if (optionsLoaded) return
  const [roleList, depts] = await Promise.all([roleApi.all(), deptApi.tree()])
  roles.value = roleList
  deptTree.value = depts
  optionsLoaded = true
}

async function openCreate() {
  await ensureFormOptions()
  editingId.value = null
  Object.assign(form, {
    username: '', password: '', nickname: '', phone: '', email: '',
    enabled: true, deptId: null, roleIds: [],
  })
  dialogVisible.value = true
}

async function openEdit(row) {
  await ensureFormOptions()
  editingId.value = row.id
  Object.assign(form, {
    username: row.username,
    password: '',
    nickname: row.nickname,
    phone: row.phone,
    email: row.email,
    enabled: row.enabled,
    deptId: row.deptId,
    roleIds: (row.roles || []).map((r) => r.id),
  })
  dialogVisible.value = true
}

async function submit() {
  if (!editingId.value && !form.password) {
    ElMessage.warning('请填写密码')
    return
  }
  if (editingId.value) {
    await userApi.update(editingId.value, { ...form })
  } else {
    await userApi.create({ ...form })
  }
  ElMessage.success('保存成功')
  dialogVisible.value = false
  load()
}

async function remove(row) {
  await ElMessageBox.confirm(`确定删除用户「${row.username}」？`, '提示', { type: 'warning' })
  await userApi.remove(row.id)
  ElMessage.success('已删除')
  load()
}

onMounted(load)
</script>
