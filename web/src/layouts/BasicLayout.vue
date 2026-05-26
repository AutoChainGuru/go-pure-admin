<template>
  <div class="h-screen flex bg-gray-50">
    <aside
      class="flex flex-col bg-white border-r border-gray-200 transition-all duration-200"
      :class="collapsed ? 'w-16' : 'w-56'"
    >
      <div class="h-14 flex-center border-b border-gray-100 px-3 gap-2">
        <div class="w-8 h-8 rounded-lg bg-indigo-500 flex-center text-white font-bold text-sm">P</div>
        <span v-show="!collapsed" class="font-semibold text-gray-800 truncate">Pure Admin</span>
      </div>
      <el-scrollbar class="flex-1">
        <el-menu
          :default-active="activeMenu"
          :collapse="collapsed"
          router
          class="layout-sidebar-menu border-none!"
        >
          <SidebarItem v-for="item in permStore.menus" :key="item.id" :item="item" />
        </el-menu>
      </el-scrollbar>
    </aside>

    <div class="flex-1 flex flex-col min-w-0">
      <header class="layout-header h-14 shrink-0 bg-white flex items-center justify-between px-4">
        <div class="flex items-center gap-3">
          <el-button text @click="collapsed = !collapsed">
            <el-icon :size="18"><component :is="collapsed ? Expand : Fold" /></el-icon>
          </el-button>
          <LayoutBreadcrumb />
        </div>
        <UserInfoDropdown @logout="logout" />
      </header>
      <LayoutTabs />
      <main class="layout-main flex-1 overflow-auto">
        <PageRouterView />
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { resetRouter } from '@/router'
import { useUserStore } from '@/stores/user'
import { usePermissionStore } from '@/stores/permission'
import { useTabsStore } from '@/stores/tabs'
import PageRouterView from '@/components/PageRouterView.vue'
import SidebarItem from '@/components/SidebarItem.vue'
import LayoutTabs from '@/components/LayoutTabs.vue'
import LayoutBreadcrumb from '@/components/LayoutBreadcrumb.vue'
import UserInfoDropdown from '@/components/UserInfoDropdown.vue'
import { Fold, Expand } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const permStore = usePermissionStore()
const tabsStore = useTabsStore()
const collapsed = ref(false)
const activeMenu = computed(() => route.path)

watch(
  () => route.fullPath,
  () => {
    tabsStore.addTab(route)
  },
  { immediate: true },
)

async function logout() {
  userStore.reset()
  permStore.reset()
  tabsStore.reset()
  resetRouter()
  router.replace('/login')
}
</script>

<style scoped>
.layout-header {
  position: relative;
  z-index: 10;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
}

.layout-sidebar-menu {
  --sidebar-menu-active-bg: #eef2ff;
  --sidebar-menu-active-color: var(--el-color-primary);
  padding: 8px;
  box-sizing: border-box;
}

.layout-sidebar-menu:deep(.el-menu-item),
.layout-sidebar-menu:deep(.el-sub-menu__title) {
  height: 44px;
  line-height: 44px;
  margin-bottom: 4px;
  border-radius: 8px;
  color: #4b5563;
}

.layout-sidebar-menu:deep(.el-menu-item:hover),
.layout-sidebar-menu:deep(.el-sub-menu__title:hover) {
  background-color: #f3f4f6;
  color: #374151;
}

.layout-sidebar-menu:deep(.el-menu-item.is-active) {
  background-color: var(--sidebar-menu-active-bg);
  color: var(--sidebar-menu-active-color);
  font-weight: 500;
}

.layout-sidebar-menu:deep(.el-menu-item.is-active .el-icon) {
  color: var(--sidebar-menu-active-color);
}

.layout-sidebar-menu:deep(.el-menu-item.is-active:hover) {
  background-color: var(--sidebar-menu-active-bg);
}

.layout-sidebar-menu:deep(.el-sub-menu.is-active > .el-sub-menu__title) {
  color: var(--sidebar-menu-active-color);
  font-weight: 500;
}

.layout-sidebar-menu:deep(.el-sub-menu.is-active > .el-sub-menu__title .el-icon) {
  color: var(--sidebar-menu-active-color);
}

.layout-sidebar-menu:deep(.el-sub-menu .el-menu) {
  padding: 4px 0 4px 8px;
  background-color: transparent;
}

.layout-sidebar-menu:deep(.el-sub-menu .el-menu-item) {
  min-width: auto;
  height: 40px;
  line-height: 40px;
}
</style>
