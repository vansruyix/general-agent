<template>
  <el-container class="app-layout">
    <el-aside width="260px" class="sidebar-panel">
      <Sidebar
        :sessions="sessionStore.sessions"
        :currentSessionId="sessionStore.currentSessionId"
        @create="handleNew"
        @select="handleSelect"
        @delete="handleDelete"
      />
    </el-aside>
    <el-main class="main-panel">
      <ChatArea />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import Sidebar from './Sidebar.vue'
import ChatArea from './ChatArea.vue'
import { useSessionStore } from '../../stores/session'
import { useChatStore } from '../../stores/chat'

const sessionStore = useSessionStore()
const chatStore = useChatStore()

/** 应用启动时如果没有会话，自动创建一个 */
onMounted(() => {
  if (sessionStore.sessions.length === 0) {
    sessionStore.createSession()
  }
})

function handleNew() {
  sessionStore.createSession()
  chatStore.clearMessages()
}

function handleSelect(id: string) {
  sessionStore.switchSession(id)
  chatStore.loadMessages()
}

function handleDelete(id: string) {
  sessionStore.deleteSession(id)
  if (sessionStore.currentSessionId) {
    chatStore.loadMessages()
  } else {
    chatStore.clearMessages()
  }
}
</script>

<style scoped>
.app-layout {
  height: 100vh;
  overflow: hidden;
}
.sidebar-panel {
  height: 100vh;
  overflow: hidden;
}
.main-panel {
  padding: 0;
  height: 100vh;
  overflow: hidden;
}
</style>