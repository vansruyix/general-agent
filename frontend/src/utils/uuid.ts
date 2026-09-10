/** 生成 UUID v4 */
export function generateUUID(): string {
  // crypto.randomUUID() 仅在安全上下文（HTTPS/localhost）可用
  // 使用 crypto.getRandomValues() 作为兼容 fallback
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID()
  }
  // fallback: 手动生成 UUID v4
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (crypto.getRandomValues(new Uint8Array(1))[0] & 15) >> (c === 'x' ? 0 : 3)
    return (c === 'x' ? r : (r & 0x3) | 0x8).toString(16)
  })
}