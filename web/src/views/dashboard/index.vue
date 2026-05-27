<template>
  <div class="page-card">
    <h2 class="text-xl font-semibold text-gray-800 mt-0">欢迎回来</h2>
    <p class="text-muted">
      {{ userStore.user?.nickname || userStore.user?.username }}，祝你今天工作顺利。
    </p>
  </div>
</template>

<script setup>
defineOptions({ name: 'Dashboard' })

import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
</script>
