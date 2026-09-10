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
  padding: 8px 12px;
  cursor: pointer;
  border-radius: var(--ga-radius);
  margin: 1px 8px;
  transition: background 0.12s;
  position: relative;
}
.session-item:hover { background: #ececf0; }
.session-item.active { background: #e8e8ed; }
.session-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 6px;
  bottom: 6px;
  width: 3px;
  border-radius: 2px;
  background: var(--el-color-primary);
}
.session-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 1px;
}
.title {
  font-size: 13px;
  line-height: 1.4;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--ga-text-primary);
}
.session-item.active .title { font-weight: 500; }
.time {
  font-size: 11px;
  color: var(--ga-text-muted);
}
.delete-btn {
  opacity: 0;
  transition: opacity 0.12s;
  margin-left: 4px;
  flex-shrink: 0;
}
.session-item:hover .delete-btn { opacity: 1; }
</style>