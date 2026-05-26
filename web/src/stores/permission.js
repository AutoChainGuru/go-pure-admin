import { defineStore } from 'pinia'
import { ref } from 'vue'
import { buildRoutesFromMenus, filterMenusForSidebar } from '@/utils/route'

export const usePermissionStore = defineStore('permission', () => {
  const menus = ref([])
  const routesAdded = ref(false)

  function setMenus(list) {
    menus.value = filterMenusForSidebar(list || [])
  }

  function generateRoutes(menuTree) {
    return buildRoutesFromMenus(menuTree)
  }

  function markRoutesAdded() {
    routesAdded.value = true
  }

  function reset() {
    menus.value = []
    routesAdded.value = false
  }

  return { menus, routesAdded, setMenus, generateRoutes, markRoutesAdded, reset }
})
