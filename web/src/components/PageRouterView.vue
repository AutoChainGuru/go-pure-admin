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
import { useAppStore } from '@/stores/app'

const { keepAliveNames } = storeToRefs(usePermissionStore())
const { refreshTick } = storeToRefs(useAppStore())

function pageKey(route) {
  let base
  if (route.name && keepAliveNames.value.includes(route.name)) {
    base = route.name
  } else {
    base = route.fullPath
  }
  return `${base}__${refreshTick.value}`
}
</script>
