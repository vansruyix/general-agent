<template>
  <div class="message-list" ref="listRef">
    <div v-if="messages.length === 0" class="empty-state">
      <p class="empty-title">General Agent</p>
      <p class="empty-desc">输入消息开始对话</p>
    </div>
    <TransitionGroup name="msg">
      <MessageBubble
        v-for="msg in messages"
        :key="msg.id"
        :message="msg"
      />
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import type { ChatMessage } from '../../types/chat'
import MessageBubble from './MessageBubble.vue'

const props = defineProps<{
  messages: ChatMessage[]
}>()

const listRef = ref<HTMLElement | null>(null)

function scrollToBottom() {
  nextTick(() => {
    if (listRef.value) {
      listRef.value.scrollTop = listRef.value.scrollHeight
    }
  })
}

// 消息数量变化或最后一条消息内容变化时自动滚动到底部
watch(() => props.messages.length, scrollToBottom)
watch(
  () => {
    const last = props.messages[props.messages.length - 1]
    return last?.content
  },
  scrollToBottom
)
</script>

<style scoped>
.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px 0;
  display: flex;
  flex-direction: column;
}
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  user-select: none;
}
.empty-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--ga-text-primary);
  letter-spacing: -0.03em;
}
.empty-desc {
  font-size: 14px;
  color: var(--ga-text-muted);
}
.msg-enter-active {
  transition: all 0.15s ease;
}
.msg-enter-from {
  opacity: 0;
  transform: translateY(6px);
}

/* 自定义滚动条 — 细竖线，无箭头 */
.message-list::-webkit-scrollbar {
  width: 5px;
}
.message-list::-webkit-scrollbar-track {
  background: transparent;
}
.message-list::-webkit-scrollbar-thumb {
  background: #d1d3d6;
  border-radius: 3px;
}
.message-list::-webkit-scrollbar-thumb:hover {
  background: #b0b3b8;
}
.message-list::-webkit-scrollbar-button {
  display: none;
}
/* Firefox */
.message-list {
  scrollbar-width: thin;
  scrollbar-color: #d1d3d6 transparent;
}
</style>