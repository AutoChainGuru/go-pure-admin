<template>
  <div class="min-h-screen flex-center bg-gradient-to-br from-indigo-50 via-white to-slate-100">
    <div class="w-full max-w-md px-6">
      <div class="page-card shadow-md!">
        <div class="text-center mb-8">
          <div class="w-12 h-12 mx-auto rounded-xl bg-indigo-500 flex-center text-white text-xl font-bold mb-3">P</div>
          <h1 class="text-2xl font-semibold text-gray-800 m-0">go-pure-admin</h1>
          <p class="text-muted mt-2">简洁现代的管理后台</p>
        </div>
        <el-form size="large" @submit.prevent="onSubmit">
          <el-form-item>
            <el-input v-model="form.username" placeholder="用户名" :prefix-icon="User" />
          </el-form-item>
          <el-form-item>
            <el-input v-model="form.password" type="password" placeholder="密码" show-password :prefix-icon="Lock" />
          </el-form-item>
          <el-button type="primary" class="w-full" :loading="loading" native-type="submit">登 录</el-button>
        </el-form>
        <p class="text-muted text-center mt-4 mb-0">默认账号 admin / admin123</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { usePermissionStore } from '@/stores/permission'
import { resetRouter } from '@/router'
import { User, Lock } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)

const form = reactive({
  username: 'admin',
  password: 'admin123',
})

async function onSubmit() {
  loading.value = true
  try {
    await userStore.login(form)
    resetRouter()
    usePermissionStore().reset()
    const redirect = route.query.redirect || '/'
    router.replace(String(redirect))
  } finally {
    loading.value = false
  }
}
</script>
