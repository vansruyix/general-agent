<template>
  <el-container class="chat-area">
    <el-main class="message-area">
      <MessageList :messages="chatStore.messages" />
    </el-main>
    <el-footer class="input-area" height="auto">
      <ChatInput
        :isStreaming="chatStore.isStreaming"
        @send="chatStore.sendMessage"
        @stop="chatStore.stopStreaming"
      />
    </el-footer>
  </el-container>
</template>

<script setup lang="ts">
import MessageList from '../chat/MessageList.vue'
import ChatInput from '../chat/ChatInput.vue'
import { useChatStore } from '../../stores/chat'

const chatStore = useChatStore()
</script>

<style scoped>
.chat-area {
  height: 100%;
}
.message-area {
  padding: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.input-area {
  padding: 0;
  background: var(--ga-bg-chat);
}
</style>