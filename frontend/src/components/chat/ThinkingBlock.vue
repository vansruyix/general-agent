<template>
  <el-collapse class="thinking-block">
    <el-collapse-item title="💭 思考过程">
      <div class="thinking-content">{{ content }}</div>
    </el-collapse-item>
  </el-collapse>
</template>

<script setup lang="ts">
defineProps<{
  content: string
}>()
</script>

<style scoped>
.thinking-block {
  margin-bottom: 6px;
}
.thinking-block :deep(.el-collapse) {
  border: 1px solid var(--ga-border);
  border-radius: var(--ga-radius-sm);
  overflow: hidden;
  --el-collapse-header-height: 30px;
}
.thinking-block :deep(.el-collapse-item__header) {
  padding: 0 10px;
  font-size: 12px;
  font-weight: 500;
  color: var(--ga-text-secondary);
  border: none;
  background: transparent;
}
.thinking-block :deep(.el-collapse-item__wrap) {
  border: none;
  background: transparent;
}
.thinking-block :deep(.el-collapse-item__content) {
  padding: 0 10px 8px;
}
.thinking-content {
  font-size: 12px;
  color: var(--ga-text-secondary);
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.7;
  font-family: var(--ga-font-mono);
}
</style>