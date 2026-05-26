<template>
  <router-view v-slot="{ Component, route }">
    <keep-alive :include="keepAliveNames" :max="30">
      <component :is="Component" :key="pageKey(route)" />
    </keep-alive>
  </router-view>
</template>

<script setup>
import { storeToRefs } from 'pinia'
import { usePermissionStore } from '@/stores/permission'

const { keepAliveNames } = storeToRefs(usePermissionStore())

function pageKey(route) {
  if (route.name && keepAliveNames.value.includes(route.name)) {
    return route.name
  }
  return route.fullPath
}
</script>
