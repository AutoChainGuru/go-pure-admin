<template>
  <div class="page-card profile-page max-w-3xl">
    <h2 class="text-lg font-semibold text-gray-800 mt-0 mb-6">个人信息</h2>

    <el-form v-loading="loading" label-width="96px" class="max-w-lg">
      <el-form-item label="用户名">
        <el-input :model-value="profile.username" disabled />
      </el-form-item>
      <el-form-item label="昵称">
        <el-input v-model="form.nickname" placeholder="显示名称" />
      </el-form-item>
      <el-form-item label="手机号">
        <el-input v-model="form.phone" placeholder="手机号" />
      </el-form-item>
      <el-form-item label="邮箱">
        <el-input v-model="form.email" placeholder="邮箱" />
      </el-form-item>
      <el-form-item label="头像地址">
        <el-input v-model="form.avatar" placeholder="头像 URL" />
      </el-form-item>
      <el-form-item label="部门">
        <el-input :model-value="profile.dept?.name || '-'" disabled />
      </el-form-item>
      <el-form-item label="角色">
        <el-input :model-value="roleText" disabled type="textarea" :rows="2" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="saving" @click="saveProfile">保存资料</el-button>
      </el-form-item>
    </el-form>

    <el-divider />

    <h3 class="text-base font-semibold text-gray-800 mt-0 mb-4">修改密码</h3>
    <el-form label-width="96px" class="max-w-lg">
      <el-form-item label="原密码" required>
        <el-input v-model="pwdForm.oldPassword" type="password" show-password autocomplete="current-password" />
      </el-form-item>
      <el-form-item label="新密码" required>
        <el-input v-model="pwdForm.newPassword" type="password" show-password autocomplete="new-password" />
      </el-form-item>
      <el-form-item label="确认密码" required>
        <el-input v-model="pwdForm.confirmPassword" type="password" show-password autocomplete="new-password" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" plain :loading="pwdSaving" @click="savePassword">修改密码</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { getProfile, updateProfile, changePassword } from '@/api/auth'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const loading = ref(false)
const saving = ref(false)
const pwdSaving = ref(false)

const profile = ref({})
const form = reactive({
  nickname: '',
  phone: '',
  email: '',
  avatar: '',
})
const pwdForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const roleText = computed(() => {
  if (profile.value.superAdmin) return '超级管理员'
  const names = (profile.value.roles || []).map((r) => r.name).filter(Boolean)
  return names.length ? names.join('、') : '-'
})

function applyProfile(data) {
  profile.value = data
  form.nickname = data.nickname || ''
  form.phone = data.phone || ''
  form.email = data.email || ''
  form.avatar = data.avatar || ''
}

async function loadProfile() {
  loading.value = true
  try {
    const data = await getProfile()
    applyProfile(data)
  } finally {
    loading.value = false
  }
}

async function saveProfile() {
  saving.value = true
  try {
    const data = await updateProfile({ ...form })
    applyProfile(data)
    await userStore.fetchInfo()
    ElMessage.success('资料已保存')
  } finally {
    saving.value = false
  }
}

async function savePassword() {
  if (!pwdForm.oldPassword || !pwdForm.newPassword) {
    ElMessage.warning('请填写原密码和新密码')
    return
  }
  if (pwdForm.newPassword.length < 6) {
    ElMessage.warning('新密码至少 6 位')
    return
  }
  if (pwdForm.newPassword !== pwdForm.confirmPassword) {
    ElMessage.warning('两次输入的新密码不一致')
    return
  }
  pwdSaving.value = true
  try {
    await changePassword({
      oldPassword: pwdForm.oldPassword,
      newPassword: pwdForm.newPassword,
    })
    pwdForm.oldPassword = ''
    pwdForm.newPassword = ''
    pwdForm.confirmPassword = ''
    ElMessage.success('密码已修改')
  } finally {
    pwdSaving.value = false
  }
}

onMounted(loadProfile)
</script>
