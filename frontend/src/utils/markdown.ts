import { marked } from 'marked'

/** 将 Markdown 文本渲染为 HTML 字符串 */
export function renderMarkdown(text: string): string {
  if (!text) return ''
  return marked.parse(text) as string
}