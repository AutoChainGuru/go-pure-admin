<template>
  <div class="page-layout">
    <div class="page-card">
      <div class="flex flex-wrap gap-3 items-center">
        <el-input v-model="query.username" placeholder="操作人" clearable class="w-36!" @keyup.enter="load" />
        <el-input v-model="query.module" placeholder="模块" clearable class="w-36!" @keyup.enter="load" />
        <el-select v-model="query.status" placeholder="状态" clearable class="w-28!">
          <el-option label="成功" :value="1" />
          <el-option label="失败" :value="0" />
        </el-select>
        <el-date-picker
          v-model="dateRange"
          type="datetimerange"
          range-separator="-"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          value-format="YYYY-MM-DD HH:mm:ss"
          :default-time="defaultTimeRange"
          class="filter-daterange"
        />
        <el-button type="primary" @click="load">查询</el-button>
        <el-button @click="reset">重置</el-button>
      </div>
    </div>

    <div class="page-card">
      <el-table v-loading="loading" :data="list" stripe>
      <el-table-column prop="module" label="模块" width="110" />
      <el-table-column prop="action" label="操作" width="120" />
      <el-table-column prop="username" label="操作人" width="100" />
      <el-table-column prop="method" label="方法" width="90" />
      <el-table-column prop="path" label="路径" min-width="200" show-overflow-tooltip />
      <el-table-column label="状态" width="72">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '成功' : '失败' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="durationMs" label="耗时(ms)" width="96" align="right" />
      <el-table-column prop="ip" label="IP" width="120" />
      <el-table-column label="时间" width="170">
        <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column label="详情" width="72" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="openDetail(row)">查看</el-button>
        </template>
      </el-table-column>
    </el-table>

      <div class="flex justify-end mt-4">
        <AppPagination v-model:page="query.page" v-model:page-size="query.pageSize" :total="total" @change="load" />
      </div>
    </div>

    <el-dialog v-model="detailVisible" title="操作日志详情" width="720px" destroy-on-close>
      <el-descriptions v-if="detail" :column="2" border size="small">
        <el-descriptions-item label="模块">{{ detail.module }}</el-descriptions-item>
        <el-descriptions-item label="操作">{{ detail.action }}</el-descriptions-item>
        <el-descriptions-item label="操作人">{{ detail.username || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="detail.status === 1 ? 'success' : 'danger'" size="small">
            {{ detail.status === 1 ? '成功' : '失败' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="请求" :span="2">{{ detail.method }} {{ detail.path }}</el-descriptions-item>
        <el-descriptions-item label="IP">{{ detail.ip || '-' }}</el-descriptions-item>
        <el-descriptions-item label="耗时">{{ detail.durationMs }} ms</el-descriptions-item>
        <el-descriptions-item v-if="detail.errorMsg" label="错误" :span="2">{{ detail.errorMsg }}</el-descriptions-item>
        <el-descriptions-item label="时间" :span="2">{{ formatTime(detail.createdAt) }}</el-descriptions-item>
      </el-descriptions>
      <div v-if="detail?.requestBody" class="mt-4">
        <div class="text-sm text-gray-500 mb-1">请求参数</div>
        <pre class="log-pre">{{ detail.requestBody }}</pre>
      </div>
      <div v-if="detail?.userAgent" class="mt-3">
        <div class="text-sm text-gray-500 mb-1">User-Agent</div>
        <pre class="log-pre text-xs">{{ detail.userAgent }}</pre>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import AppPagination from '@/components/AppPagination.vue'
import { operLogApi } from '@/api/system'

const loading = ref(false)
const list = ref([])
const total = ref(0)
const dateRange = ref(null)
const defaultTimeRange = [
  new Date(2000, 0, 1, 0, 0, 0),
  new Date(2000, 0, 1, 23, 59, 59),
]
const detailVisible = ref(false)
const detail = ref(null)

const query = reactive({
  page: 1,
  pageSize: 10,
  username: '',
  module: '',
  status: undefined,
})

function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString()
}

function buildParams() {
  const params = { ...query }
  if (dateRange.value?.length === 2) {
    params.startTime = dateRange.value[0]
    params.endTime = dateRange.value[1]
  }
  if (params.status === undefined || params.status === '') {
    delete params.status
  }
  return params
}

async function load() {
  loading.value = true
  try {
    const data = await operLogApi.list(buildParams())
    list.value = data.list
    total.value = data.total
  } finally {
    loading.value = false
  }
}

async function openDetail(row) {
  detail.value = await operLogApi.get(row.id)
  detailVisible.value = true
}

function reset() {
  query.username = ''
  query.module = ''
  query.status = undefined
  dateRange.value = null
  query.page = 1
  load()
}

onMounted(load)
</script>

<style scoped>
.log-pre {
  margin: 0;
  padding: 10px;
  max-height: 240px;
  overflow: auto;
  font-size: 12px;
  line-height: 1.5;
  background: #f5f7fa;
  border-radius: 6px;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
