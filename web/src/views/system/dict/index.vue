<template>
  <div class="dict-split">
    <!-- 左侧：字典类型 -->
    <div class="dict-panel page-card flex flex-col min-w-0">
      <div class="font-medium text-gray-800 mb-3">字典类型</div>
      <div class="mb-3">
        <el-button type="primary" :icon="Plus" @click="openDictCreate">新增字典</el-button>
      </div>
      <div class="dict-table-wrap flex-1 min-h-0">
        <el-table
          ref="dictTableRef"
          v-loading="dictLoading"
          :data="dictList"
          height="100%"
          highlight-current-row
          size="small"
          row-key="id"
          @row-click="selectDict"
        >
          <el-table-column prop="name" label="名称" min-width="100" show-overflow-tooltip />
          <el-table-column prop="type" label="类型" min-width="90" show-overflow-tooltip />
          <el-table-column label="操作" width="88" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" size="small" @click="openDictEdit(row, $event)">编辑</el-button>
              <el-button link type="danger" size="small" @click="removeDict(row, $event)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div class="flex justify-end mt-3 shrink-0">
        <AppPagination
          v-model:page="query.page"
          v-model:page-size="query.pageSize"
          :total="dictTotal"
          size="small"
          @change="loadDicts"
        />
      </div>
    </div>

    <!-- 右侧：字典项 -->
    <div class="detail-panel page-card flex flex-col min-w-0">
      <div class="flex items-center justify-between gap-3 mb-3 shrink-0">
        <div class="min-w-0">
          <div class="font-medium text-gray-800">字典项</div>
          <p v-if="selectedDict" class="text-muted m-0 mt-1 truncate">
            {{ selectedDict.name }}
            <span class="text-gray-400 font-mono text-xs ml-1">{{ selectedDict.type }}</span>
          </p>
          <p v-else class="text-muted m-0 mt-1">请在左侧选择字典</p>
        </div>
        <el-button type="primary" :icon="Plus" :disabled="!selectedDict" @click="openDetailCreate">
          新增字典项
        </el-button>
      </div>
      <div class="dict-table-wrap flex-1 min-h-0">
        <el-empty v-if="!selectedDict" description="选择左侧字典后管理字典项" class="h-full flex-center" />
        <el-table v-else v-loading="detailLoading" :data="details" height="100%" stripe size="small">
          <el-table-column prop="label" label="标签" min-width="100" />
          <el-table-column prop="value" label="值" min-width="100" />
          <el-table-column prop="sort" label="排序" width="64" align="center" />
          <el-table-column label="状态" width="72" align="center">
            <template #default="{ row }">
              <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
                {{ row.enabled ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" min-width="100" show-overflow-tooltip />
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" size="small" @click="openDetailEdit(row)">编辑</el-button>
              <el-button link type="danger" size="small" @click="removeDetail(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <el-dialog v-model="dictDialogVisible" :title="dictEditingId ? '编辑字典' : '新增字典'" width="440px" destroy-on-close>
      <el-form label-width="80px">
        <el-form-item label="名称" required><el-input v-model="dictForm.name" /></el-form-item>
        <el-form-item label="类型" required>
          <el-input v-model="dictForm.type" :disabled="!!dictEditingId" placeholder="唯一标识，如 sys_status" />
        </el-form-item>
        <el-form-item label="描述"><el-input v-model="dictForm.description" type="textarea" :rows="2" /></el-form-item>
        <el-form-item label="状态"><el-switch v-model="dictForm.enabled" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dictDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitDict">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="detailDialogVisible"
      :title="detailEditingId ? '编辑字典项' : '新增字典项'"
      width="440px"
      destroy-on-close
    >
      <el-form label-width="80px">
        <el-form-item label="标签" required><el-input v-model="detailForm.label" /></el-form-item>
        <el-form-item label="值" required><el-input v-model="detailForm.value" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="detailForm.sort" :min="0" class="w-full!" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="detailForm.remark" type="textarea" :rows="2" /></el-form-item>
        <el-form-item label="状态"><el-switch v-model="detailForm.enabled" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="detailDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitDetail">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { nextTick, onMounted, reactive, ref, watch } from 'vue'
import AppPagination from '@/components/AppPagination.vue'
import { dictApi } from '@/api/system'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

const dictTableRef = ref()
const dictLoading = ref(false)
const detailLoading = ref(false)
const dictList = ref([])
const dictTotal = ref(0)
const details = ref([])
const selectedDict = ref(null)

const dictDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const dictEditingId = ref(null)
const detailEditingId = ref(null)

const query = reactive({ page: 1, pageSize: 10 })
const dictForm = reactive({ name: '', type: '', description: '', enabled: true })
const detailForm = reactive({
  dictId: null,
  label: '',
  value: '',
  sort: 0,
  enabled: true,
  remark: '',
})

async function loadDicts() {
  dictLoading.value = true
  try {
    const data = await dictApi.list(query)
    dictList.value = data.list
    dictTotal.value = data.total
    if (selectedDict.value) {
      const still = data.list.find((d) => d.id === selectedDict.value.id)
      selectedDict.value = still || data.list[0] || null
    } else if (data.list.length) {
      selectedDict.value = data.list[0]
    }
    await nextTickSetCurrentRow()
  } finally {
    dictLoading.value = false
  }
}

async function nextTickSetCurrentRow() {
  await nextTick()
  if (selectedDict.value) {
    dictTableRef.value?.setCurrentRow(selectedDict.value)
  }
}

async function loadDetails() {
  if (!selectedDict.value?.id) {
    details.value = []
    return
  }
  detailLoading.value = true
  try {
    details.value = await dictApi.listDetails(selectedDict.value.id)
  } finally {
    detailLoading.value = false
  }
}

watch(selectedDict, () => {
  loadDetails()
})

function selectDict(row) {
  selectedDict.value = row
}

function openDictCreate() {
  dictEditingId.value = null
  Object.assign(dictForm, { name: '', type: '', description: '', enabled: true })
  dictDialogVisible.value = true
}

function openDictEdit(row, e) {
  e?.stopPropagation?.()
  dictEditingId.value = row.id
  Object.assign(dictForm, {
    name: row.name,
    type: row.type,
    description: row.description,
    enabled: row.enabled,
  })
  dictDialogVisible.value = true
}

async function submitDict() {
  if (!dictForm.name?.trim() || !dictForm.type?.trim()) {
    ElMessage.warning('请填写名称和类型')
    return
  }
  if (dictEditingId.value) {
    await dictApi.update(dictEditingId.value, { ...dictForm })
  } else {
    await dictApi.create({ ...dictForm })
  }
  ElMessage.success('保存成功')
  dictDialogVisible.value = false
  await loadDicts()
}

async function removeDict(row, e) {
  e?.stopPropagation?.()
  await ElMessageBox.confirm(`确定删除字典「${row.name}」？其下字典项将一并删除。`, '提示', { type: 'warning' })
  await dictApi.remove(row.id)
  ElMessage.success('已删除')
  if (selectedDict.value?.id === row.id) {
    selectedDict.value = null
  }
  loadDicts()
}

function openDetailCreate() {
  if (!selectedDict.value) {
    ElMessage.warning('请先选择左侧字典')
    return
  }
  detailEditingId.value = null
  Object.assign(detailForm, {
    dictId: selectedDict.value.id,
    label: '',
    value: '',
    sort: details.value.length,
    enabled: true,
    remark: '',
  })
  detailDialogVisible.value = true
}

function openDetailEdit(row) {
  detailEditingId.value = row.id
  Object.assign(detailForm, {
    dictId: row.dictId,
    label: row.label,
    value: row.value,
    sort: row.sort,
    enabled: row.enabled,
    remark: row.remark || '',
  })
  detailDialogVisible.value = true
}

async function submitDetail() {
  if (!detailForm.label?.trim() || !detailForm.value?.trim()) {
    ElMessage.warning('请填写标签和值')
    return
  }
  if (detailEditingId.value) {
    await dictApi.updateDetail(detailEditingId.value, { ...detailForm })
  } else {
    await dictApi.createDetail({ ...detailForm })
  }
  ElMessage.success('保存成功')
  detailDialogVisible.value = false
  loadDetails()
}

async function removeDetail(row) {
  await ElMessageBox.confirm(`确定删除字典项「${row.label}」？`, '提示', { type: 'warning' })
  await dictApi.removeDetail(row.id)
  ElMessage.success('已删除')
  loadDetails()
}

onMounted(async () => {
  await loadDicts()
})
</script>

<style scoped>
.dict-split {
  display: flex;
  gap: var(--layout-content-gap);
  width: 100%;
  min-height: 32rem;
  height: calc(100vh - 10.5rem);
}

.dict-panel {
  flex: 3;
  width: 380px;
  max-width: 40%;
  min-height: 0;
}

.detail-panel {
  flex: 7;
  min-width: 0;
  min-height: 0;
}

.dict-table-wrap {
  min-height: 12rem;
}

.dict-split :deep(.el-table__body tr) {
  cursor: pointer;
}
</style>
