import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const DASHBOARD_PATH = '/dashboard'
const TABS_STORAGE_KEY = 'go-pure-admin-tabs'

function createDashboardTab() {
  return {
    path: DASHBOARD_PATH,
    title: '工作台',
    name: 'Dashboard',
    affix: true,
  }
}

function normalizePath(path) {
  if (!path) return path
  if (path.length > 1 && path.endsWith('/')) return path.slice(0, -1)
  return path
}

function serializeTabs(list) {
  return list.map(({ path, fullPath, title, name, affix }) => ({
    path,
    fullPath,
    title,
    name,
    affix: !!affix,
  }))
}

function loadTabsFromStorage() {
  try {
    const raw = localStorage.getItem(TABS_STORAGE_KEY)
    if (!raw) return null
    const parsed = JSON.parse(raw)
    if (!Array.isArray(parsed) || parsed.length === 0) return null
    const list = parsed
      .filter((t) => t && typeof t.path === 'string')
      .map((t) => {
        const path = normalizePath(t.path)
        return {
          path,
          fullPath: t.fullPath || path,
          title: t.title || path,
          name: t.name,
          affix: path === DASHBOARD_PATH || !!t.affix,
        }
      })
    if (!list.some((t) => t.path === DASHBOARD_PATH)) {
      list.unshift(createDashboardTab())
    }
    return list.length ? list : null
  } catch {
    return null
  }
}

function saveTabsToStorage(list) {
  try {
    localStorage.setItem(TABS_STORAGE_KEY, JSON.stringify(serializeTabs(list)))
  } catch {
    // ignore quota / private mode errors
  }
}

function createInitialTabs() {
  if (!localStorage.getItem('token')) {
    return [createDashboardTab()]
  }
  return loadTabsFromStorage() || [createDashboardTab()]
}

export function canAddTab(route) {
  if (route.meta?.public) return false
  if (!route.name || ['Root', 'NotFound', 'Login'].includes(route.name)) return false
  const last = route.matched[route.matched.length - 1]
  return !!last?.meta?.title
}

export const useTabsStore = defineStore('tabs', () => {
  const tabs = ref(createInitialTabs())

  watch(
    tabs,
    (list) => {
      if (localStorage.getItem('token')) {
        saveTabsToStorage(list)
      }
    },
    { deep: true },
  )

  function ensureAffix() {
    if (!tabs.value.some((t) => t.path === DASHBOARD_PATH)) {
      tabs.value.unshift(createDashboardTab())
    }
  }

  function addTab(route) {
    if (!canAddTab(route)) return
    ensureAffix()
    const path = normalizePath(route.path)
    if (tabs.value.some((t) => t.path === path)) return
    tabs.value.push({
      path,
      fullPath: route.fullPath,
      title: route.meta?.title || String(route.name),
      name: route.name,
      affix: path === DASHBOARD_PATH,
    })
  }

  function removePaths(paths, router, currentPath) {
    const removeSet = new Set(paths.filter((p) => p !== DASHBOARD_PATH))
    if (!removeSet.size) return

    const active = normalizePath(currentPath)
    const wasActive = removeSet.has(active)
    let nextPath = DASHBOARD_PATH
    if (wasActive) {
      const idx = tabs.value.findIndex((t) => t.path === active)
      const right = tabs.value[idx + 1]
      const left = tabs.value[idx - 1]
      nextPath = (right || left || tabs.value.find((t) => t.affix))?.path || DASHBOARD_PATH
    }

    tabs.value = tabs.value.filter((t) => !removeSet.has(t.path))
    ensureAffix()

    if (wasActive) {
      router.push(nextPath)
    }
  }

  function closeTab(path, router, currentPath) {
    const p = normalizePath(path)
    const tab = tabs.value.find((t) => t.path === p)
    if (!tab || tab.affix) return
    removePaths([p], router, currentPath)
  }

  function closeLeft(targetPath, router, currentPath) {
    const idx = tabs.value.findIndex((t) => t.path === normalizePath(targetPath))
    if (idx <= 0) return
    const paths = tabs.value.slice(0, idx).filter((t) => !t.affix).map((t) => t.path)
    removePaths(paths, router, currentPath)
  }

  function closeRight(targetPath, router, currentPath) {
    const idx = tabs.value.findIndex((t) => t.path === normalizePath(targetPath))
    if (idx < 0) return
    const paths = tabs.value.slice(idx + 1).filter((t) => !t.affix).map((t) => t.path)
    removePaths(paths, router, currentPath)
  }

  function closeOthers(targetPath, router, currentPath) {
    const keep = normalizePath(targetPath)
    const paths = tabs.value.filter((t) => t.path !== keep && !t.affix).map((t) => t.path)
    removePaths(paths, router, currentPath)
  }

  function closeAll(router) {
    tabs.value = tabs.value.filter((t) => t.affix)
    ensureAffix()
    router.push(DASHBOARD_PATH)
  }

  function reset() {
    tabs.value = [createDashboardTab()]
    localStorage.removeItem(TABS_STORAGE_KEY)
  }

  /** 动态路由注入后移除已失效页签，避免 router 解析 /system/dict 等报 No match */
  function pruneInvalid(router) {
    tabs.value = tabs.value.filter((t) => {
      if (t.affix) return true
      const resolved = router.resolve(t.path)
      return resolved.matched.length > 0 && resolved.name !== 'NotFound'
    })
    ensureAffix()
  }

  return {
    tabs,
    addTab,
    closeTab,
    closeLeft,
    closeRight,
    closeOthers,
    closeAll,
    reset,
    pruneInvalid,
  }
})
