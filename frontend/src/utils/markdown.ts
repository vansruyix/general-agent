import { marked } from 'marked'
import DOMPurify from 'dompurify'

/** 将 Markdown 文本渲染为安全的 HTML 字符串 */
export function renderMarkdown(text: string): string {
  if (!text) return ''
  const rawHtml = marked.parse(text) as string
  return DOMPurify.sanitize(rawHtml)
}

/** 对 HTML 字符串做 XSS 净化，供其他组件使用 */
export function sanitizeHtml(html: string): string {
  return DOMPurify.sanitize(html)
}