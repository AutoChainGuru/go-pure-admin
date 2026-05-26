import { childRoutePath, stripLeadingSlash } from '@/utils/menu-path'

const viewModules = import.meta.glob('@/views/**/*.vue')

function resolveView(component, menuType) {
  if (menuType === 0 || component === 'Layout') {
    return () => import('@/layouts/ParentLayout.vue')
  }
  if (!component) {
    return () => import('@/views/error/404.vue')
  }
  const key = `/src/views/${component}.vue`
  return viewModules[key] || viewModules[`/src/views/${component}/index.vue`] || (() => import('@/views/error/404.vue'))
}

function menuToRoute(menu, parentMenuPath = '') {
  const routePath = parentMenuPath
    ? childRoutePath(menu.path, parentMenuPath)
    : stripLeadingSlash(menu.path)

  const route = {
    path: routePath,
    name: menu.name || `Menu${menu.id}`,
    meta: {
      title: menu.title,
      icon: menu.icon,
      hidden: menu.hidden,
      keepAlive: menu.keepAlive,
    },
    component: resolveView(menu.component, menu.menuType),
  }

  const childMenus = (menu.children || []).filter((c) => c.menuType !== 2 && c.enabled !== false)
  if (childMenus.length) {
    const parentPath = stripLeadingSlash(menu.path)
    route.children = childMenus.map((c) => menuToRoute(c, parentPath))
    const first = route.children.find((c) => !c.meta?.hidden)
    if (first) {
      route.redirect = first.name ? { name: first.name } : first.path
    }
  }

  return route
}

export function buildRoutesFromMenus(menus) {
  return menus
    .filter((m) => m.menuType !== 2 && m.enabled !== false)
    .filter((m) => stripLeadingSlash(m.path) !== 'dashboard')
    .map((m) => menuToRoute(m))
}

/** 角色分配等场景：仅保留已启用菜单 */
export function filterEnabledMenuTree(menus) {
  return menus
    .filter((m) => m.enabled !== false)
    .map((m) => ({
      ...m,
      children: m.children?.length ? filterEnabledMenuTree(m.children) : [],
    }))
}

/** 侧栏展示：已启用且非隐藏（按钮由 SidebarItem 再过滤） */
export function filterMenusForSidebar(menus) {
  return menus
    .filter((m) => m.enabled !== false && !m.hidden && m.menuType !== 2)
    .map((m) => ({
      ...m,
      children: m.children?.length ? filterMenusForSidebar(m.children) : [],
    }))
}
