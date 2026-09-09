<template>
  <div class="message-list" ref="listRef">
    <TransitionGroup name="msg">
      <MessageBubble
        v-for="msg in messages"
        :key="msg.id"
        :message="msg"
      />
    </TransitionGroup>
    <el-empty
      v-if="messages.length === 0"
      description="开始新的对话吧"
      :image-size="80"
    />
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
.msg-enter-active {
  transition: all 0.2s ease;
}
.msg-enter-from {
  opacity: 0;
  transform: translateY(10px);
}
</style>