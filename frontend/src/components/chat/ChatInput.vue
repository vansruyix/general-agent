<template>
  <div class="chat-input">
    <div class="input-wrapper">
      <el-input
        v-model="input"
        type="textarea"
        :rows="2"
        :autosize="{ minRows: 2, maxRows: 5 }"
        placeholder="输入消息..."
        :disabled="isStreaming"
        resize="none"
        @keydown="handleKeydown"
      />
      <div class="input-actions">
        <span class="shortcut-hint">Enter 发送 · Shift+Enter 换行</span>
        <el-button
          v-if="!isStreaming"
          type="primary"
          size="small"
          :disabled="!input.trim()"
          :icon="Promotion"
          class="send-btn"
          @click="send"
        >发送</el-button>
        <el-button
          v-else
          type="danger"
          size="small"
          :icon="Close"
          class="stop-btn"
          @click="$emit('stop')"
        >停止</el-button>
      </div>
    </div>
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
  max-width: 768px;
  margin: 0 auto;
  padding: 16px 16px 20px;
}
.input-wrapper {
  position: relative;
}
.chat-input :deep(.el-textarea__inner) {
  resize: none;
  border-radius: 16px;
  border-color: var(--ga-border);
  font-size: 14px;
  line-height: 1.5;
  padding: 12px 16px 44px;
  min-height: 68px !important;
  height: auto !important;
  box-shadow: 0 0 0 1px rgba(0,0,0,.04), 0 2px 8px rgba(0,0,0,.06);
  transition: border-color .15s, box-shadow .15s;
}
.chat-input :deep(.el-textarea__inner:focus) {
  border-color: var(--el-color-primary);
  box-shadow: 0 0 0 3px var(--el-color-primary-light-7);
  outline: none;
}
.input-actions {
  position: absolute;
  right: 12px;
  bottom: 8px;
  display: flex;
  align-items: center;
  gap: 10px;
}
.shortcut-hint {
  font-size: 11px;
  color: var(--ga-text-muted);
  user-select: none;
  white-space: nowrap;
}
.send-btn, .stop-btn {
  border-radius: 8px;
  font-weight: 500;
  font-size: 12px;
  padding: 5px 12px;
}
</style>