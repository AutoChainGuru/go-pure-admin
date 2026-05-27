<template>
  <div class="layout-tabs shrink-0">
    <el-scrollbar>
      <transition-group
        name="tab-list"
        tag="div"
        class="tabs-track flex items-end px-2 gap-1 min-w-max"
      >
        <div
          v-for="tab in tabsStore.tabs"
          :key="tab.path"
          class="tab-item"
          :class="{ 'is-active': isActive(tab) }"
          @click="go(tab)"
          @contextmenu.prevent="openContextMenu($event, tab)"
        >
          <span class="tab-title">{{ tab.title }}</span>
          <el-icon
            v-if="!tab.affix"
            class="tab-close"
            @click.stop="close(tab)"
          >
            <Close />
          </el-icon>
        </div>
      </transition-group>
    </el-scrollbar>

    <ul
      v-show="contextMenu.visible"
      class="tab-context-menu"
      :style="{ left: `${contextMenu.x}px`, top: `${contextMenu.y}px` }"
      @click.stop
    >
      <li :class="{ disabled: !canCloseLeft }" @click="runClose('left')">关闭左侧</li>
      <li :class="{ disabled: !canCloseRight }" @click="runClose('right')">关闭右侧</li>
      <li :class="{ disabled: !canCloseOthers }" @click="runClose('others')">关闭其他</li>
      <li :class="{ disabled: !canCloseAll }" @click="runClose('all')">关闭全部</li>
    </ul>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Close } from '@element-plus/icons-vue'
import { useTabsStore } from '@/stores/tabs'

const route = useRoute()
const router = useRouter()
const tabsStore = useTabsStore()

const contextMenu = reactive({
  visible: false,
  x: 0,
  y: 0,
  tab: null,
})

const currentPath = computed(() => route.path)

const contextIndex = computed(() =>
  tabsStore.tabs.findIndex((t) => t.path === contextMenu.tab?.path),
)

const canCloseLeft = computed(() => {
  if (contextIndex.value <= 0) return false
  return tabsStore.tabs.slice(0, contextIndex.value).some((t) => !t.affix)
})

const canCloseRight = computed(() => {
  if (contextIndex.value < 0) return false
  return tabsStore.tabs.slice(contextIndex.value + 1).some((t) => !t.affix)
})

const canCloseOthers = computed(() => {
  if (!contextMenu.tab) return false
  return tabsStore.tabs.some((t) => t.path !== contextMenu.tab.path && !t.affix)
})

const canCloseAll = computed(() => tabsStore.tabs.some((t) => !t.affix))

function isActive(tab) {
  return tab.path === route.path || tab.path === route.path.replace(/\/$/, '')
}

function go(tab) {
  if (isActive(tab)) return
  router.push(tab.path)
}

function close(tab) {
  tabsStore.closeTab(tab.path, router, currentPath.value)
}

function openContextMenu(e, tab) {
  contextMenu.tab = tab
  contextMenu.x = e.clientX
  contextMenu.y = e.clientY
  contextMenu.visible = true
}

function hideContextMenu() {
  contextMenu.visible = false
  contextMenu.tab = null
}

function runClose(type) {
  const tab = contextMenu.tab
  if (!tab) return
  const path = tab.path
  if (type === 'left' && canCloseLeft.value) {
    tabsStore.closeLeft(path, router, currentPath.value)
  } else if (type === 'right' && canCloseRight.value) {
    tabsStore.closeRight(path, router, currentPath.value)
  } else if (type === 'others' && canCloseOthers.value) {
    tabsStore.closeOthers(path, router, currentPath.value)
  } else if (type === 'all' && canCloseAll.value) {
    tabsStore.closeAll(router)
  }
  hideContextMenu()
}

onMounted(() => {
  document.addEventListener('click', hideContextMenu)
})

onUnmounted(() => {
  document.removeEventListener('click', hideContextMenu)
})
</script>

<style scoped>
.layout-tabs {
  position: relative;
  z-index: 9;
  height: 42px;
  box-sizing: border-box;
  background: var(--layout-tabs-bg, #fff);
  border-bottom: 1px solid var(--layout-tabs-border, #ebeef5);
}

.layout-tabs :deep(.el-scrollbar) {
  height: 42px;
}

.layout-tabs :deep(.el-scrollbar__wrap) {
  height: 42px !important;
}

.tabs-track {
  position: relative;
  height: 42px;
  box-sizing: border-box;
  align-items: flex-end;
}

.tab-list-move {
  transition: transform 0.22s ease;
}

.tab-list-enter-active {
  transition: opacity 0.22s ease, transform 0.22s ease;
}

.tab-list-leave-active {
  position: absolute;
  bottom: 0;
  transition: opacity 0.18s ease, transform 0.18s ease;
}

.tab-list-enter-from {
  opacity: 0;
  transform: translateY(-6px) scale(0.96);
}

.tab-list-leave-to {
  opacity: 0;
  transform: scale(0.92);
}

.tab-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  height: 36px;
  max-width: 180px;
  padding: 0 8px;
  font-size: 13px;
  line-height: 1;
  color: var(--layout-tab-text, #606266);
  border: 1px solid transparent;
  border-radius: 4px 4px 0 0;
  cursor: pointer;
  user-select: none;
  transition:
    color 0.2s ease,
    background 0.2s ease,
    border-color 0.2s ease,
    box-shadow 0.2s ease,
    transform 0.15s ease;
}

.tab-item:active:not(.is-active) {
  transform: scale(0.97);
}

.tab-item:hover {
  color: var(--el-color-primary);
  background: var(--layout-tab-hover-bg, #f5f7fa);
}

.tab-item.is-active {
  color: var(--el-color-primary);
  background: var(--layout-tab-active-bg, #ecf5ff);
  border-color: var(--layout-tab-active-border, #d9ecff);
  border-bottom-color: var(--layout-tabs-bg, #fff);
  box-shadow: 0 1px 4px rgb(64 158 255 / 12%);
}

.tab-title {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tab-close {
  flex-shrink: 0;
  font-size: 11px;
  border-radius: 50%;
  padding: 1px;
}

.tab-close:hover {
  color: #fff;
  background: #b1b3b8;
}

.tab-context-menu {
  position: fixed;
  z-index: 3000;
  margin: 0;
  padding: 4px 0;
  min-width: 128px;
  list-style: none;
  background: var(--layout-panel-bg, #fff);
  border: 1px solid var(--layout-border, #e4e7ed);
  border-radius: 4px;
  box-shadow: 0 2px 12px rgb(0 0 0 / 10%);
}

.tab-context-menu li {
  padding: 8px 16px;
  font-size: 13px;
  color: #606266;
  cursor: pointer;
}

.tab-context-menu li:hover:not(.disabled) {
  color: #409eff;
  background: #ecf5ff;
}

.tab-context-menu li.disabled {
  color: #c0c4cc;
  cursor: not-allowed;
}
</style>
