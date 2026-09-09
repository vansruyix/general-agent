<template>
  <div class="sidebar">
    <div class="sidebar-header">
      <span class="logo">General Agent</span>
    </div>
    <el-button
      type="primary"
      class="new-session-btn"
      :icon="Plus"
      @click="$emit('new')"
    >新会话</el-button>
    <SessionList
      :sessions="sessions"
      :currentSessionId="currentSessionId"
      @select="$emit('select', $event)"
      @delete="$emit('delete', $event)"
    />
  </div>
</template>

<script setup lang="ts">
import { Plus } from '@element-plus/icons-vue'
import type { Session } from '../../types/session'
import SessionList from '../session/SessionList.vue'

defineProps<{
  sessions: Session[]
  currentSessionId: string | null
}>()

defineEmits<{
  new: []
  select: [id: string]
  delete: [id: string]
}>()
</script>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--el-bg-color-page);
  border-right: 1px solid var(--el-border-color-light);
}
.sidebar-header {
  padding: 16px;
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  border-bottom: 1px solid var(--el-border-color-lighter);
}
.new-session-btn {
  margin: 12px;
}
</style>