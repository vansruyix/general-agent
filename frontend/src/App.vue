<template>
  <AppLayout />
</template>

<script setup lang="ts">
import AppLayout from './components/layout/AppLayout.vue'
</script>

<style>
:root {
  --ga-bg-sidebar: #f7f7f8;
  --ga-bg-chat: #ffffff;
  --ga-bg-user-bubble: #eff6ff;
  --ga-bg-code: #f4f5f6;
  --ga-border: #ecedf0;
  --ga-radius-sm: 6px;
  --ga-radius: 8px;
  --ga-radius-lg: 12px;
  --ga-font: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC", "Microsoft YaHei", "Helvetica Neue", Arial, sans-serif;
  --ga-font-mono: "Cascadia Code", "Fira Code", "JetBrains Mono", "SF Mono", Consolas, Monaco, monospace;
  --ga-text-primary: #1a1a1a;
  --ga-text-secondary: #666;
  --ga-text-muted: #999;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: var(--ga-font);
  font-size: 14px;
  line-height: 1.6;
  color: var(--ga-text-primary);
  background: var(--ga-bg-chat);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  width: 100vw;
  height: 100vh;
  overflow: hidden;
}

/* 键盘焦点 — 对所有可聚焦元素提供清晰可见的焦点环 */
:focus-visible {
  outline: 2px solid var(--el-color-primary);
  outline-offset: 2px;
  border-radius: 2px;
}

/* 减少动画 — 尊重用户系统偏好 */
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
</style>