<template>
  <div class="message-bubble" :class="message.role">
    <el-avatar
      :size="32"
      :icon="message.role === 'user' ? UserFilled : Service"
      class="avatar"
    />
    <div class="body">
      <ThinkingBlock v-if="message.thinking" :content="message.thinking" />
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
      <div v-if="message.toolCalls?.length" class="tool-calls">
        <el-collapse>
          <el-collapse-item
            v-for="tc in message.toolCalls"
            :key="tc.callId"
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
  max-width: 85%;
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
}
.body {
  flex: 1;
  min-width: 0;
}
.content {
  font-size: 14px;
  line-height: 1.7;
  word-break: break-word;
}
.content :deep(pre) {
  background: var(--el-fill-color-light);
  padding: 12px;
  border-radius: 6px;
  overflow-x: auto;
}
.content :deep(code) {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
}
.content :deep(p) {
  margin: 4px 0;
}
.user .content {
  background: var(--el-color-primary-light-9);
  padding: 10px 14px;
  border-radius: 12px 4px 12px 12px;
}
.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 8px 0;
}
.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--el-text-color-secondary);
  animation: bounce 1.4s infinite ease-in-out both;
}
.dot:nth-child(1) { animation-delay: -0.32s; }
.dot:nth-child(2) { animation-delay: -0.16s; }
@keyframes bounce {
  0%, 80%, 100% { transform: scale(0); }
  40% { transform: scale(1); }
}
.tool-calls {
  margin-top: 8px;
}
.tool-section {
  margin-bottom: 8px;
}
.tool-section pre {
  background: var(--el-fill-color);
  padding: 8px;
  border-radius: 4px;
  font-size: 12px;
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-all;
}
.error-mark {
  color: var(--el-color-danger);
}
</style>