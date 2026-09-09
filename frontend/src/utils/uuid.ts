/** 生成 UUID v4，使用浏览器原生 crypto API */
export function generateUUID(): string {
  return crypto.randomUUID()
}