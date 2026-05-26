<template>
  <div class="page-card">
    <h2 class="text-xl font-semibold text-gray-800 mt-0">欢迎回来</h2>
    <p class="text-muted">
      {{ userStore.user?.nickname || userStore.user?.username }}，祝你今天工作顺利。
    </p>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
      <div class="p-4 rounded-lg bg-indigo-50 border border-indigo-100">
        <div class="text-indigo-600 text-sm">角色</div>
        <div class="text-2xl font-bold text-indigo-700 mt-1">{{ userStore.roles.length }}</div>
      </div>
      <div class="p-4 rounded-lg bg-emerald-50 border border-emerald-100">
        <div class="text-emerald-600 text-sm">按钮权限</div>
        <div class="text-2xl font-bold text-emerald-700 mt-1">{{ userStore.permissions.length }}</div>
      </div>
      <div class="p-4 rounded-lg bg-amber-50 border border-amber-100">
        <div class="text-amber-600 text-sm">账号类型</div>
        <div class="text-lg font-semibold text-amber-800 mt-2">
          {{ userStore.user?.superAdmin ? '超级管理员' : '普通用户' }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineOptions({ name: 'Dashboard' })

import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
</script>
