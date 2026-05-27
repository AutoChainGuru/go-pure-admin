<template>
  <el-dropdown trigger="click" placement="bottom-end" @command="onCommand">
    <div class="user-trigger">
      <el-avatar :size="32" :src="userStore.user?.avatar || undefined" class="user-avatar">
        {{ avatarText }}
      </el-avatar>
      <span class="user-nickname">{{ displayName }}</span>
      <el-icon class="user-arrow"><ArrowDown /></el-icon>
    </div>
    <template #dropdown>
      <div class="user-dropdown-panel">
        <div class="user-dropdown-header">
          <el-avatar :size="48" :src="userStore.user?.avatar || undefined">
            {{ avatarText }}
          </el-avatar>
          <div class="user-dropdown-meta">
            <div class="meta-row">
              <span class="meta-label">用户名</span>
              <span class="meta-value">{{ userStore.user?.username || '-' }}</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">角色</span>
              <span class="meta-value">{{ roleText }}</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">部门</span>
              <span class="meta-value">{{ deptText }}</span>
            </div>
          </div>
        </div>
        <el-dropdown-menu>
          <el-dropdown-item command="profile">
            <el-icon><User /></el-icon>
            <span>个人信息</span>
          </el-dropdown-item>
          <el-dropdown-item divided command="logout">
            <el-icon class="text-red-500"><SwitchButton /></el-icon>
            <span class="text-red-500">退出登录</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </div>
    </template>
  </el-dropdown>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowDown, SwitchButton, User } from '@element-plus/icons-vue'
import { ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'

const emit = defineEmits(['logout'])

const router = useRouter()
const userStore = useUserStore()

const displayName = computed(
  () => userStore.user?.nickname || userStore.user?.username || '用户',
)

const avatarText = computed(() => {
  const name = displayName.value
  return name ? name.charAt(0).toUpperCase() : 'U'
})

const roleText = computed(() => {
  if (userStore.user?.superAdmin) return '超级管理员'
  const names = (userStore.roles || []).map((r) => r.name).filter(Boolean)
  return names.length ? names.join('、') : '-'
})

const deptText = computed(() => userStore.user?.dept?.name || '-')

function onCommand(cmd) {
  if (cmd === 'profile') {
    router.push('/profile')
    return
  }
  if (cmd === 'logout') {
    ElMessageBox.confirm('确定退出登录？', '提示', { type: 'warning' })
      .then(() => emit('logout'))
      .catch(() => {})
  }
}
</script>

<style scoped>
.user-trigger {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 4px 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s;
}

.user-trigger:hover {
  background: var(--layout-hover-bg, #f5f7fa);
}

.user-avatar {
  flex-shrink: 0;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
  font-weight: 600;
}

.user-nickname {
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
  color: var(--layout-text-strong, #303133);
}

.user-arrow {
  font-size: 12px;
  color: var(--layout-text-muted, #909399);
}

.user-dropdown-panel {
  min-width: 260px;
}

.user-dropdown-header {
  display: flex;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--layout-tabs-border, #ebeef5);
}

.user-dropdown-meta {
  flex: 1;
  min-width: 0;
}

.meta-row {
  display: flex;
  gap: 8px;
  font-size: 13px;
  line-height: 1.6;
}

.meta-label {
  flex-shrink: 0;
  color: #909399;
  width: 42px;
}

.meta-value {
  color: var(--layout-text-strong, #303133);
  word-break: break-all;
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
