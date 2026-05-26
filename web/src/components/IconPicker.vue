<template>
  <div class="icon-picker">
    <el-popover placement="bottom-start" :width="400" trigger="click">
      <template #reference>
        <el-input
          :model-value="displayLabel"
          readonly
          placeholder="点击选择图标"
          class="icon-picker-trigger"
        >
          <template #prefix>
            <el-icon v-if="resolvedIcon" class="text-base">
              <component :is="resolvedIcon" />
            </el-icon>
          </template>
        </el-input>
      </template>
      <div class="icon-picker-panel">
        <el-input v-model="keyword" placeholder="搜索图标名称" clearable class="mb-2" />
        <el-scrollbar height="260px">
          <div class="icon-grid">
            <div
              v-for="name in filteredIcons"
              :key="name"
              class="icon-item"
              :class="{ 'is-active': modelValue === name }"
              :title="name"
              @click="select(name)"
            >
              <el-icon :size="18"><component :is="Icons[name]" /></el-icon>
              <span class="icon-name">{{ name }}</span>
            </div>
          </div>
        </el-scrollbar>
      </div>
    </el-popover>
    <el-button v-if="modelValue" link type="primary" @click="clear">清除</el-button>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import * as Icons from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue'])

const keyword = ref('')

const iconNames = Object.keys(Icons).sort()

const resolvedIcon = computed(() => {
  if (!props.modelValue) return null
  return Icons[props.modelValue] || Icons.Menu
})

const displayLabel = computed(() => props.modelValue || '')

const filteredIcons = computed(() => {
  const q = keyword.value.trim().toLowerCase()
  if (!q) return iconNames
  return iconNames.filter((n) => n.toLowerCase().includes(q))
})

function select(name) {
  emit('update:modelValue', name)
}

function clear() {
  emit('update:modelValue', '')
}
</script>

<style scoped>
.icon-picker {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.icon-picker-trigger {
  flex: 1;
}

.icon-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 6px;
}

.icon-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px 4px;
  border-radius: 6px;
  cursor: pointer;
  border: 1px solid transparent;
  transition: background 0.15s, border-color 0.15s;
}

.icon-item:hover {
  background: #f5f7fa;
}

.icon-item.is-active {
  background: #ecf5ff;
  border-color: #b3d8ff;
  color: #409eff;
}

.icon-name {
  font-size: 10px;
  line-height: 1.2;
  text-align: center;
  word-break: break-all;
  max-width: 100%;
  color: #909399;
}

.icon-item.is-active .icon-name {
  color: #409eff;
}
</style>
