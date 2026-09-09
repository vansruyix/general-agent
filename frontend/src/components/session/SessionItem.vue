<template>
  <div
    class="session-item"
    :class="{ active: isActive }"
    @click="$emit('select')"
  >
    <div class="session-info">
      <span class="title">{{ session.title }}</span>
      <span class="time">{{ formatTime(session.updatedAt) }}</span>
    </div>
    <el-button
      class="delete-btn"
      size="small"
      circle
      :icon="Delete"
      @click.stop="$emit('delete')"
    />
  </div>
</template>

<script setup lang="ts">
import { Delete } from '@element-plus/icons-vue'
import type { Session } from '../../types/session'

defineProps<{
  session: Session
  isActive: boolean
}>()

defineEmits<{
  select: []
  delete: []
}>()

function formatTime(ts: number): string {
  const d = new Date(ts)
  const now = new Date()
  if (d.toDateString() === now.toDateString()) {
    return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  return d.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}
</script>

<style scoped>
.session-item {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  cursor: pointer;
  border-radius: 6px;
  margin: 2px 8px;
  transition: background 0.15s;
}
.session-item:hover { background: var(--el-fill-color-light); }
.session-item.active { background: var(--el-color-primary-light-9); }
.session-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.title {
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.time {
  font-size: 11px;
  color: var(--el-text-color-secondary);
}
.delete-btn {
  opacity: 0;
  transition: opacity 0.15s;
}
.session-item:hover .delete-btn { opacity: 1; }
</style>