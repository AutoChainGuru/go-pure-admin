import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginApi, getAuthInfo } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(null)
  const roles = ref([])
  const permissions = ref([])

  function setToken(t) {
    token.value = t
    if (t) localStorage.setItem('token', t)
    else localStorage.removeItem('token')
  }

  async function login(form) {
    const data = await loginApi(form)
    setToken(data.token)
  }

  async function fetchInfo() {
    const data = await getAuthInfo()
    user.value = data.user
    roles.value = data.roles || []
    permissions.value = data.permissions || []
    return data
  }

  function reset() {
    setToken('')
    user.value = null
    roles.value = []
    permissions.value = []
  }

  function hasPerm(code) {
    if (user.value?.superAdmin) return true
    return permissions.value.includes(code)
  }

  return { token, user, roles, permissions, setToken, login, fetchInfo, reset, hasPerm }
})
