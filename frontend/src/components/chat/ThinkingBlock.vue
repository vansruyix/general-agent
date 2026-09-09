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
  margin-bottom: 8px;
}
.thinking-content {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.6;
}
</style>