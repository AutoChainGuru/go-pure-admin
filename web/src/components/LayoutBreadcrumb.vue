<template>
  <el-breadcrumb separator="/" class="layout-breadcrumb">
    <el-breadcrumb-item v-for="(item, index) in items" :key="item.path">
      <span v-if="index === items.length - 1" class="is-current">{{ item.title }}</span>
      <router-link v-else :to="item.path" class="breadcrumb-link">{{ item.title }}</router-link>
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { getBreadcrumbs } from '@/utils/breadcrumb'

const route = useRoute()
const items = computed(() => getBreadcrumbs(route))
</script>

<style scoped>
.layout-breadcrumb {
  font-size: 12px;
  line-height: 1;
}

.layout-breadcrumb :deep(.el-breadcrumb__separator) {
  color: #c0c4cc;
  font-weight: 400;
}

.layout-breadcrumb :deep(.el-breadcrumb__item:not(:last-child) .el-breadcrumb__inner),
.layout-breadcrumb :deep(.el-breadcrumb__item:not(:last-child) .el-breadcrumb__inner.is-link) {
  font-weight: 400;
}

.layout-breadcrumb :deep(.el-breadcrumb__item:last-child .el-breadcrumb__inner) {
  font-weight: 600;
}

.breadcrumb-link {
  color: #909399;
  font-weight: 400;
  text-decoration: none;
  transition: color 0.15s;
}

.breadcrumb-link:hover {
  color: #409eff;
}

.is-current {
  color: #909399;
  font-weight: 600;
}
</style>
