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
        <el-menu :default-active="activeMenu" :collapse="collapsed" router class="border-none!">
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
        <router-view v-slot="{ Component, route: viewRoute }">
          <transition name="page-slide" mode="out-in">
            <div
              :key="viewRoute.meta?.keepAlive ? viewRoute.name : viewRoute.fullPath"
              class="page-view min-h-0"
            >
              <keep-alive v-if="viewRoute.meta?.keepAlive">
                <component :is="Component" />
              </keep-alive>
              <component v-else :is="Component" />
            </div>
          </transition>
        </router-view>
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

.page-slide-enter-active,
.page-slide-leave-active {
  transition:
    opacity 0.22s ease,
    transform 0.22s ease;
}

.page-slide-enter-from {
  opacity: 0;
  transform: translateX(10px);
}

.page-slide-leave-to {
  opacity: 0;
  transform: translateX(-10px);
}
</style>
