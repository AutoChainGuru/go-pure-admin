<template>
  <div class="page-layout">
    <div class="page-card">
      <div class="flex flex-wrap gap-3 items-center">
        <el-input v-model="query.username" placeholder="用户名" clearable class="w-40!" @keyup.enter="load" />
        <el-select v-model="query.status" placeholder="状态" clearable class="w-28!">
          <el-option label="成功" :value="1" />
          <el-option label="失败" :value="0" />
        </el-select>
        <el-input v-model="query.ip" placeholder="IP" clearable class="w-40!" @keyup.enter="load" />
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
      <el-table-column prop="username" label="用户名" width="120" />
      <el-table-column prop="nickname" label="昵称" width="120" />
      <el-table-column label="状态" width="88">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '成功' : '失败' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="message" label="说明" min-width="160" show-overflow-tooltip />
      <el-table-column prop="ip" label="IP" width="130" />
      <el-table-column prop="userAgent" label="浏览器" min-width="200" show-overflow-tooltip />
      <el-table-column label="登录时间" width="170">
        <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
      </el-table-column>
    </el-table>

      <div class="flex justify-end mt-4">
        <AppPagination v-model:page="query.page" v-model:page-size="query.pageSize" :total="total" @change="load" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import AppPagination from '@/components/AppPagination.vue'
import { loginLogApi } from '@/api/system'

const loading = ref(false)
const list = ref([])
const total = ref(0)
const dateRange = ref(null)
const defaultTimeRange = [
  new Date(2000, 0, 1, 0, 0, 0),
  new Date(2000, 0, 1, 23, 59, 59),
]

const query = reactive({
  page: 1,
  pageSize: 10,
  username: '',
  status: undefined,
  ip: '',
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
    const data = await loginLogApi.list(buildParams())
    list.value = data.list
    total.value = data.total
  } finally {
    loading.value = false
  }
}

function reset() {
  query.username = ''
  query.status = undefined
  query.ip = ''
  dateRange.value = null
  query.page = 1
  load()
}

onMounted(load)
</script>
