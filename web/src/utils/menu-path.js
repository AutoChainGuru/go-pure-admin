/** 去掉开头斜杠 */
export function stripLeadingSlash(p) {
  return p?.startsWith('/') ? p.slice(1) : p || ''
}

/**
 * 将菜单 path 转为相对父级的路由片段（用于 vue-router 嵌套 children）
 * 例如父级 system + 菜单 /system/dict → dict
 */
export function childRoutePath(menuPath, parentMenuPath) {
  let p = stripLeadingSlash(menuPath || '')
  const parent = stripLeadingSlash(parentMenuPath || '')
  if (!parent) return p
  if (p === parent) return ''
  if (p.startsWith(`${parent}/`)) return p.slice(parent.length + 1)
  const parts = p.split('/').filter(Boolean)
  const parentParts = parent.split('/').filter(Boolean)
  if (
    parts.length > parentParts.length &&
    parts.slice(0, parentParts.length).join('/') === parent
  ) {
    return parts.slice(parentParts.length).join('/')
  }
  return p
}

/** 侧栏 / 页签使用的完整访问路径 */
export function resolveMenuPath(menuPath, basePath = '') {
  const p = menuPath || ''
  if (p.startsWith('/')) return p.replace(/\/+/g, '/')
  const base = basePath.replace(/\/$/, '')
  const parentSeg = stripLeadingSlash(base.replace(/^\//, ''))
  const segment = parentSeg ? childRoutePath(p, parentSeg) : stripLeadingSlash(p)
  return `${base}/${segment}`.replace(/\/+/g, '/')
}
