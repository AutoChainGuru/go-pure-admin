<template>
  <el-pagination
    :current-page="page"
    :page-size="pageSize"
    :page-sizes="pageSizes"
    :total="total"
    :size="size"
    layout="total, sizes, prev, pager, next, jumper"
    background
    @size-change="onSizeChange"
    @current-change="onCurrentChange"
  />
</template>

<script setup>
const page = defineModel('page', { type: Number, default: 1 })
const pageSize = defineModel('pageSize', { type: Number, default: 10 })

defineProps({
  total: { type: Number, default: 0 },
  size: { type: String, default: 'default' },
})

const emit = defineEmits(['change'])

const pageSizes = [10, 20, 50, 100]

function onSizeChange(size) {
  pageSize.value = size
  page.value = 1
  emit('change')
}

function onCurrentChange(p) {
  page.value = p
  emit('change')
}
</script>
