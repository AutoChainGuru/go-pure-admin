import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { usePermissionStore } from '@/stores/permission'
import { useTabsStore } from '@/stores/tabs'

/** 固定路由：工作台在动态菜单注入前就必须可访问 */
const constantRoutes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { public: true },
  },
  {
    path: '/',
    name: 'Root',
    component: () => import('@/layouts/BasicLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '工作台', icon: 'Odometer' },
      },
      {
        path: 'profile',
        name: 'UserProfile',
        component: () => import('@/views/profile/index.vue'),
        meta: { title: '个人信息', hidden: true },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes: constantRoutes,
})

function addNotFoundRoute() {
  if (router.hasRoute('NotFound')) return
  router.addRoute({
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
  })
}

router.beforeEach(async (to) => {
  const userStore = useUserStore()
  const permStore = usePermissionStore()

  if (to.meta.public) {
    if (to.path === '/login' && userStore.token) {
      return '/'
    }
    return true
  }

  if (!userStore.token) {
    return `/login?redirect=${encodeURIComponent(to.fullPath)}`
  }

  if (!permStore.routesAdded) {
    try {
      const info = await userStore.fetchInfo()
      const authMenus = info.menus || []
      permStore.setMenus(authMenus)
      const dynamicRoutes = permStore.generateRoutes(authMenus)
      dynamicRoutes.forEach((r) => router.addRoute('Root', r))
      addNotFoundRoute()
      permStore.markRoutesAdded()
      useTabsStore().pruneInvalid(router)
      return { ...to, replace: true }
    } catch {
      userStore.reset()
      permStore.reset()
      return '/login'
    }
  }

  if (to.matched.length === 0) {
    return { name: 'NotFound', replace: true }
  }

  return true
})

const KEEP_ROUTE_NAMES = new Set(['Login', 'Root', 'Dashboard', 'UserProfile'])

export function resetRouter() {
  router.getRoutes().forEach((route) => {
    if (route.name && !KEEP_ROUTE_NAMES.has(route.name)) {
      router.removeRoute(route.name)
    }
  })
}

export default router
