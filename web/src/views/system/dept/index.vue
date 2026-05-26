<template>
  <div class="page-card">
    <div class="mb-4"><el-button type="primary" @click="openCreate()">新增部门</el-button></div>
    <el-table v-loading="loading" :data="tree" row-key="id" default-expand-all>
      <el-table-column prop="name" label="部门名称" min-width="200" />
      <el-table-column prop="sort" label="排序" width="80" />
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="row.enabled ? 'success' : 'info'" size="small">{{ row.enabled ? '启用' : '禁用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button link type="primary" @click="openCreate(row)">子部门</el-button>
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button link type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑部门' : '新增部门'" width="400px">
      <el-form label-width="80px">
        <el-form-item label="名称" required><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" /></el-form-item>
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
defineOptions({ name: 'SystemDept' })

import { onMounted, ref } from 'vue'
import { deptApi } from '@/api/system'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const tree = ref([])
const dialogVisible = ref(false)
const editingId = ref(null)
const form = ref({ parentId: null, name: '', sort: 0, enabled: true })

async function load() {
  loading.value = true
  try {
    tree.value = await deptApi.tree()
  } finally {
    loading.value = false
  }
}

function openCreate(parent) {
  editingId.value = null
  form.value = { parentId: parent?.id || null, name: '', sort: 0, enabled: true }
  dialogVisible.value = true
}

function openEdit(row) {
  editingId.value = row.id
  form.value = { parentId: row.parentId, name: row.name, sort: row.sort, enabled: row.enabled }
  dialogVisible.value = true
}

async function submit() {
  if (editingId.value) await deptApi.update(editingId.value, form.value)
  else await deptApi.create(form.value)
  ElMessage.success('保存成功')
  dialogVisible.value = false
  load()
}

async function remove(row) {
  await ElMessageBox.confirm(`确定删除部门「${row.name}」？`, '提示', { type: 'warning' })
  await deptApi.remove(row.id)
  ElMessage.success('已删除')
  load()
}

onMounted(load)
</script>
