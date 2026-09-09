// Vitest setup: 在 Node.js 环境下提供 localStorage polyfill
// Node.js v26+ 引入了实验性的全局 localStorage（默认为 undefined），
// 会覆盖 jsdom 提供的 localStorage，需要手动 polyfill

const store = new Map<string, string>()

globalThis.localStorage = {
  getItem(key: string): string | null {
    return store.get(key) ?? null
  },
  setItem(key: string, value: string): void {
    store.set(key, value)
  },
  removeItem(key: string): void {
    store.delete(key)
  },
  clear(): void {
    store.clear()
  },
  get length(): number {
    return store.size
  },
  key(index: number): string | null {
    return Array.from(store.keys())[index] ?? null
  },
}