/**
 * 根据当前路由 matched 生成面包屑
 * @param {import('vue-router').RouteLocationNormalizedLoaded} route
 */
export function getBreadcrumbs(route) {
  const items = []
  let base = ''

  for (const record of route.matched) {
    if (!record.meta?.title || record.name === 'Root') continue
    if (record.meta.breadcrumb === false) continue

    if (record.path.startsWith('/')) {
      base = record.path
    } else {
      base = `${base}/${record.path}`.replace(/\/+/g, '/')
    }
    items.push({ title: record.meta.title, path: base })
  }

  if (!items.length) {
    return [{ title: '工作台', path: '/dashboard' }]
  }

  return items
}
