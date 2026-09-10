<template>
  <div class="sidebar">
    <div class="sidebar-header">
      <span class="logo">General Agent</span>
    </div>
    <el-button
      type="primary"
      class="new-session-btn"
      :icon="Plus"
      @click="$emit('create')"
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
  create: []
  select: [id: string]
  delete: [id: string]
}>()
</script>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--ga-bg-sidebar);
  border-right: 1px solid var(--ga-border);
}
.sidebar-header {
  padding: 20px 16px 12px;
  font-size: 15px;
  font-weight: 700;
  color: var(--ga-text-primary);
  letter-spacing: -0.01em;
}
.new-session-btn {
  margin: 8px 12px 4px;
  border-radius: var(--ga-radius);
  font-weight: 500;
}
</style>