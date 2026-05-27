<template>
  <div class="layout-shell h-screen flex">
    <aside
      class="layout-panel layout-aside flex flex-col border-r transition-all duration-200"
      :class="collapsed ? 'w-16' : 'w-56'"
    >
      <div class="layout-brand h-14 flex-center border-b px-3 gap-2">
        <div class="w-8 h-8 rounded-lg bg-indigo-500 flex-center text-white font-bold text-sm">P</div>
        <span v-show="!collapsed" class="layout-brand-title font-semibold truncate">Pure Admin</span>
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
      <header class="layout-panel layout-header h-14 shrink-0 flex items-center justify-between px-4">
        <div class="flex items-center gap-3">
          <el-button text @click="collapsed = !collapsed">
            <el-icon :size="18"><component :is="collapsed ? Expand : Fold" /></el-icon>
          </el-button>
          <LayoutBreadcrumb />
        </div>
        <div class="header-toolbar flex items-center gap-1">
          <el-tooltip content="刷新页面" placement="bottom">
            <el-button class="header-tool-btn" text circle @click="appStore.refreshPage()">
              <el-icon :size="18"><Refresh /></el-icon>
            </el-button>
          </el-tooltip>
          <el-tooltip :content="appStore.isDark ? '切换浅色' : '切换深色'" placement="bottom">
            <el-button class="header-tool-btn" text circle @click="appStore.toggleTheme()">
              <el-icon :size="18">
                <component :is="appStore.isDark ? Sunny : Moon" />
              </el-icon>
            </el-button>
          </el-tooltip>
          <UserInfoDropdown @logout="logout" />
        </div>
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
import { useAppStore } from '@/stores/app'
import PageRouterView from '@/components/PageRouterView.vue'
import SidebarItem from '@/components/SidebarItem.vue'
import LayoutTabs from '@/components/LayoutTabs.vue'
import LayoutBreadcrumb from '@/components/LayoutBreadcrumb.vue'
import UserInfoDropdown from '@/components/UserInfoDropdown.vue'
import { Fold, Expand, Refresh, Moon, Sunny } from '@element-plus/icons-vue'

const route = useRoute()
const appStore = useAppStore()
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
.layout-shell {
  background: var(--layout-shell-bg, #f9fafb);
}

.layout-panel {
  background: var(--layout-panel-bg, #fff);
  border-color: var(--layout-border, #e5e7eb);
}

.layout-brand {
  border-color: var(--layout-border-light, #f3f4f6);
}

.layout-brand-title {
  color: var(--layout-text-strong, #1f2937);
}

.layout-header {
  position: relative;
  z-index: 10;
  box-shadow: var(--layout-header-shadow, 0 1px 4px rgba(0, 0, 0, 0.06));
}

.header-toolbar :deep(.header-tool-btn) {
  color: var(--layout-text-muted, #6b7280);
}

.header-toolbar :deep(.header-tool-btn:hover) {
  color: var(--el-color-primary);
  background: var(--layout-hover-bg, #f3f4f6);
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
