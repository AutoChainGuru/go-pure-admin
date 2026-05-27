import { defineStore } from 'pinia'
import { ref } from 'vue'

const THEME_STORAGE_KEY = 'pure-admin-theme'

function readStoredDark() {
  return localStorage.getItem(THEME_STORAGE_KEY) === 'dark'
}

export function applyDocumentTheme(isDark) {
  document.documentElement.classList.toggle('dark', isDark)
}

applyDocumentTheme(readStoredDark())

export const useAppStore = defineStore('app', () => {
  const isDark = ref(readStoredDark())
  const refreshTick = ref(0)

  function toggleTheme() {
    isDark.value = !isDark.value
    localStorage.setItem(THEME_STORAGE_KEY, isDark.value ? 'dark' : 'light')
    applyDocumentTheme(isDark.value)
  }

  function refreshPage() {
    refreshTick.value += 1
  }

  return { isDark, refreshTick, toggleTheme, refreshPage }
})
