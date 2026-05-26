<template>
  <template v-if="item.enabled !== false && !item.hidden && item.menuType !== 2">
    <el-sub-menu v-if="visibleChildren.length" :index="indexPath">
      <template #title>
        <el-icon><component :is="iconComponent" /></el-icon>
        <span>{{ item.title }}</span>
      </template>
      <SidebarItem
        v-for="child in visibleChildren"
        :key="child.id"
        :item="child"
        :base-path="indexPath"
      />
    </el-sub-menu>
    <el-menu-item v-else :index="indexPath">
      <el-icon><component :is="iconComponent" /></el-icon>
      <span>{{ item.title }}</span>
    </el-menu-item>
  </template>
</template>

<script setup>
import * as Icons from '@element-plus/icons-vue'
import { computed } from 'vue'
import { resolveMenuPath } from '@/utils/menu-path'

const props = defineProps({
  item: { type: Object, required: true },
  basePath: { type: String, default: '' },
})

const visibleChildren = computed(() =>
  (props.item.children || []).filter(
    (c) => c.enabled !== false && c.menuType !== 2 && !c.hidden,
  ),
)

const iconComponent = computed(() => Icons[props.item.icon] || Icons.Menu)

const indexPath = computed(() => resolvePath(props.item, props.basePath))

function resolvePath(menu, parent) {
  return resolveMenuPath(menu.path, parent)
}
</script>
