import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  buildRoutesFromMenus,
  collectKeepAliveNames,
  collectKeepAliveNamesFromRouter,
  filterMenusForSidebar,
} from '@/utils/route'

export const usePermissionStore = defineStore('permission', () => {
  const menus = ref([])
  /** 开启页面缓存的组件名（统一名单，模块内 + 跨模块共用） */
  const keepAliveNames = ref([])
  const routesAdded = ref(false)

  function setMenus(list) {
    const tree = list || []
    menus.value = filterMenusForSidebar(tree)
    keepAliveNames.value = collectKeepAliveNames(tree)
  }

  function syncKeepAliveFromRouter(router) {
    const fromRouter = collectKeepAliveNamesFromRouter(router)
    if (fromRouter.length) {
      keepAliveNames.value = [...new Set([...keepAliveNames.value, ...fromRouter])]
    }
  }

  function generateRoutes(menuTree) {
    return buildRoutesFromMenus(menuTree)
  }

  function markRoutesAdded() {
    routesAdded.value = true
  }

  function reset() {
    menus.value = []
    keepAliveNames.value = []
    routesAdded.value = false
  }

  return {
    menus,
    keepAliveNames,
    routesAdded,
    setMenus,
    syncKeepAliveFromRouter,
    generateRoutes,
    markRoutesAdded,
    reset,
  }
})
