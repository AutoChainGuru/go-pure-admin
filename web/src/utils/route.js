import { childRoutePath, stripLeadingSlash } from '@/utils/menu-path'

const viewModules = import.meta.glob('@/views/**/*.vue')

function toAbsolutePath(relativePath) {
  const p = stripLeadingSlash(relativePath || '')
  return p ? `/${p}` : '/'
}

function joinMenuRoutePath(parentMenuPath, menuPath) {
  const parent = stripLeadingSlash(parentMenuPath || '')
  const segment = parent ? childRoutePath(menuPath, parent) : stripLeadingSlash(menuPath)
  if (!parent) return segment
  if (!segment) return parent
  return `${parent}/${segment}`
}

function resolveView(component, routeName) {
  if (!component) {
    return () => import('@/views/error/404.vue')
  }
  const key = `/src/views/${component}.vue`
  const loader =
    viewModules[key] ||
    viewModules[`/src/views/${component}/index.vue`] ||
    (() => import('@/views/error/404.vue'))

  return async () => {
    const mod = await loader()
    if (routeName && mod.default) {
      mod.default.name = routeName
    }
    return mod
  }
}

/** 目录：menu_type=0 或 component=Layout（兼容误设为页面的目录） */
function isDirectoryMenu(menu) {
  return Number(menu.menuType) === 0 || menu.component === 'Layout'
}

function menuToRoute(menu, relativePath, parentTitle = '') {
  const routeName = menu.name || `Menu${menu.id}`
  return {
    path: toAbsolutePath(relativePath),
    name: routeName,
    meta: {
      title: menu.title,
      parentTitle: parentTitle || undefined,
      icon: menu.icon,
      hidden: menu.hidden,
      keepAlive: !!menu.keepAlive,
    },
    component: resolveView(menu.component, routeName),
  }
}

/**
 * 目录仅作侧栏分组；页面以 /system/user 等绝对路径注册为 Root 子路由，
 * 全部在 BasicLayout 同一 router-view 渲染，单一 keep-alive 即可模块内/跨模块缓存。
 */
function flattenRoutesFromMenus(menus, parentMenuPath = '', parentTitle = '') {
  const routes = []
  for (const m of menus || []) {
    if (m.menuType === 2 || m.enabled === false) continue

    const fullPath = joinMenuRoutePath(parentMenuPath, m.path)

    if (isDirectoryMenu(m)) {
      routes.push(...flattenRoutesFromMenus(m.children || [], fullPath, m.title))

      const first = (m.children || []).find(
        (c) => Number(c.menuType) === 1 && c.enabled !== false && !c.hidden,
      )
      if (first) {
        routes.push({
          path: toAbsolutePath(fullPath),
          redirect: toAbsolutePath(joinMenuRoutePath(fullPath, first.path)),
        })
      }
      continue
    }

    if (Number(m.menuType) === 1 && !isDirectoryMenu(m)) {
      routes.push(menuToRoute(m, fullPath, parentTitle))
    }
  }
  return routes
}

/** 需缓存的组件名（与 defineOptions.name / 路由 name 一致） */
export function collectKeepAliveNames(menus) {
  const names = []
  function walk(list) {
    for (const m of list || []) {
      if (Number(m.menuType) === 1 && m.keepAlive && m.name) {
        names.push(m.name)
      }
      walk(m.children)
    }
  }
  walk(menus)
  return names
}

export function collectKeepAliveNamesFromRouter(router) {
  return [
    ...new Set(
      router
        .getRoutes()
        .filter((r) => r.meta?.keepAlive && r.name)
        .map((r) => String(r.name)),
    ),
  ]
}

export function buildRoutesFromMenus(menus) {
  return flattenRoutesFromMenus(menus).filter((r) => {
    if (r.redirect) return true
    return stripLeadingSlash(r.path) !== 'dashboard'
  })
}

export function filterEnabledMenuTree(menus) {
  return menus
    .filter((m) => m.enabled !== false)
    .map((m) => ({
      ...m,
      children: m.children?.length ? filterEnabledMenuTree(m.children) : [],
    }))
}

export function filterMenusForSidebar(menus) {
  return menus
    .filter((m) => m.enabled !== false && !m.hidden && m.menuType !== 2)
    .map((m) => ({
      ...m,
      children: m.children?.length ? filterMenusForSidebar(m.children) : [],
    }))
}
