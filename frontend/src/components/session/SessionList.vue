<template>
  <div class="session-list">
    <SessionItem
      v-for="session in sessions"
      :key="session.id"
      :session="session"
      :isActive="session.id === currentSessionId"
      @select="$emit('select', session.id)"
      @delete="$emit('delete', session.id)"
    />
    <el-empty
      v-if="sessions.length === 0"
      description="暂无会话"
      :image-size="60"
    />
  </div>
</template>

<script setup lang="ts">
import type { Session } from '../../types/session'
import SessionItem from './SessionItem.vue'

defineProps<{
  sessions: Session[]
  currentSessionId: string | null
}>()

defineEmits<{
  select: [id: string]
  delete: [id: string]
}>()
</script>

<style scoped>
.session-list {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;
}
</style>