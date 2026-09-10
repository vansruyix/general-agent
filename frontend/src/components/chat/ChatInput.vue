<template>
  <div class="chat-input">
    <el-input
      v-model="input"
      type="textarea"
      :rows="1"
      :autosize="{ minRows: 1, maxRows: 4 }"
      placeholder="输入消息... (Enter 发送, Shift+Enter 换行)"
      :disabled="isStreaming"
      resize="none"
      @keydown="handleKeydown"
    />
    <el-button
      v-if="!isStreaming"
      type="primary"
      :disabled="!input.trim()"
      :icon="Promotion"
      @click="send"
    >发送</el-button>
    <el-button
      v-else
      type="danger"
      :icon="Close"
      @click="$emit('stop')"
    >停止</el-button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Promotion, Close } from '@element-plus/icons-vue'

defineProps<{
  isStreaming: boolean
}>()

const emit = defineEmits<{
  send: [text: string]
  stop: []
}>()

const input = ref('')

/** 监听键盘事件，Enter 发送（Shift+Enter 换行，输入法组合中不触发） */
function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey && !e.isComposing) {
    e.preventDefault()
    send()
  }
}

function send() {
  const text = input.value.trim()
  if (!text) return
  emit('send', text)
  input.value = ''
}
</script>

<style scoped>
.chat-input {
  display: flex;
  gap: 10px;
  align-items: flex-end;
  padding: 12px 16px;
  border-top: 1px solid var(--el-border-color-light);
  background: var(--el-bg-color);
}
.chat-input :deep(.el-textarea__inner) {
  resize: none;
}
</style>