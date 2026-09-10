<template>
  <div class="message-bubble" :class="message.role">
    <el-avatar
      :size="32"
      :icon="message.role === 'user' ? UserFilled : Service"
      class="avatar"
    />
    <div class="body">
      <ThinkingBlock v-if="message.thinking" :content="message.thinking" />
      <div v-if="message.toolCalls?.length" class="tool-calls">
        <el-collapse>
          <el-collapse-item
            v-for="tc in message.toolCalls"
            :key="tc.callId || tc.toolName"
            :title="`🔧 ${tc.toolName || '工具调用'}`"
          >
            <div class="tool-section">
              <strong>参数:</strong>
              <pre>{{ tc.params }}</pre>
            </div>
            <div v-if="tc.result" class="tool-section">
              <strong>结果:</strong>
              <pre>{{ tc.result }}</pre>
            </div>
          </el-collapse-item>
        </el-collapse>
      </div>
      <div
        v-if="message.content"
        class="content"
        v-html="renderedContent"
      ></div>
      <div
        v-if="message.status === 'streaming' && !message.content"
        class="typing-indicator"
      >
        <span class="dot"></span><span class="dot"></span><span class="dot"></span>
      </div>
      <span v-if="message.status === 'error'" class="error-mark">⚠️</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { UserFilled, Service } from '@element-plus/icons-vue'
import type { ChatMessage } from '../../types/chat'
import { renderMarkdown } from '../../utils/markdown'
import ThinkingBlock from './ThinkingBlock.vue'

const props = defineProps<{
  message: ChatMessage
}>()

const renderedContent = computed(() => renderMarkdown(props.message.content))
</script>

<style scoped>
.message-bubble {
  display: flex;
  gap: 10px;
  padding: 8px 16px;
  max-width: 88%;
}
.message-bubble.user {
  flex-direction: row-reverse;
  align-self: flex-end;
  margin-left: auto;
}
.message-bubble.assistant {
  align-self: flex-start;
}
.avatar {
  flex-shrink: 0;
  margin-top: 2px;
}
.body {
  flex: 1;
  min-width: 0;
}
.content {
  font-size: 14px;
  line-height: 1.75;
  word-break: break-word;
  color: var(--ga-text-primary);
}
.user .body {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}
.user .content {
  background: var(--ga-bg-user-bubble);
  padding: 10px 14px;
  border-radius: var(--ga-radius-lg) 4px var(--ga-radius-lg) var(--ga-radius-lg);
  max-width: max-content;
}
.content :deep(pre) {
  background: var(--ga-bg-code);
  padding: 12px 14px;
  border-radius: var(--ga-radius-sm);
  overflow-x: auto;
  border: 1px solid var(--ga-border);
  border-left: 3px solid var(--el-color-primary);
  margin: 8px 0;
  line-height: 1.55;
}
.content :deep(code) {
  font-family: var(--ga-font-mono);
  font-size: 13px;
  color: #1a1a1a;
}
.content :deep(p) {
  margin: 0;
}
.content :deep(p + p) {
  margin-top: 8px;
}
.content :deep(ul), .content :deep(ol) {
  padding-left: 20px;
  margin: 4px 0;
}
.content :deep(li) {
  margin: 2px 0;
}
.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 6px 0;
}
.dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--ga-text-muted);
  animation: bounce 1.4s infinite ease-in-out both;
}
.dot:nth-child(1) { animation-delay: -0.32s; }
.dot:nth-child(2) { animation-delay: -0.16s; }
@keyframes bounce {
  0%, 80%, 100% { transform: scale(0.4); }
  40% { transform: scale(1); }
}
.tool-calls {
  margin-bottom: 6px;
  font-size: 13px;
}
.tool-calls :deep(.el-collapse) {
  border: none;
  --el-collapse-header-height: 32px;
}
.tool-calls :deep(.el-collapse-item) {
  background: var(--ga-bg-code);
  border: 1px solid var(--ga-border);
  border-radius: var(--ga-radius-sm);
  margin-bottom: 4px;
  overflow: hidden;
}
.tool-calls :deep(.el-collapse-item__header) {
  padding: 0 10px;
  font-size: 13px;
  font-weight: 500;
  color: var(--ga-text-secondary);
  border: none;
  background: transparent;
}
.tool-calls :deep(.el-collapse-item__wrap) {
  border: none;
  background: transparent;
}
.tool-calls :deep(.el-collapse-item__content) {
  padding: 0 10px 8px;
}
.tool-section {
  margin-bottom: 4px;
}
.tool-section strong {
  font-weight: 500;
  color: var(--ga-text-secondary);
  font-size: 12px;
}
.tool-section pre {
  background: #fff;
  padding: 6px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-family: var(--ga-font-mono);
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-all;
  margin-top: 2px;
  border: 1px solid var(--ga-border);
}
.error-mark { color: var(--el-color-danger); }
</style>