import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useChatStore } from '../chat'
import { useSessionStore } from '../session'

// Mock 适配器
vi.mock('../../adapters', () => ({
  default: {
    chatStream: vi.fn((_id: string, _question: string, onEvent: any, onError: any, onDone: any) => {
      // 模拟异步 SSE 事件
      setTimeout(() => onEvent({ type: 'text', data: '你好' }), 0)
      setTimeout(() => onEvent({ type: 'text', data: '，世界' }), 0)
      setTimeout(() => onEvent({ type: 'done', data: '' }), 0)
      setTimeout(() => onDone(), 0)
      return new AbortController()
    }),
  },
}))

describe('chatStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    // 先创建一个会话
    const sessionStore = useSessionStore()
    sessionStore.createSession()
  })

  it('sendMessage 应该添加用户消息和 AI 占位消息', async () => {
    const store = useChatStore()
    store.sendMessage('你好')

    // 用户消息立即添加
    expect(store.messages.length).toBe(2)
    expect(store.messages[0].role).toBe('user')
    expect(store.messages[0].content).toBe('你好')
    expect(store.messages[1].role).toBe('assistant')
    expect(store.isStreaming).toBe(true)
  })

  it('stopStreaming 应该停止流并设置状态', () => {
    const store = useChatStore()
    store.sendMessage('test')
    store.stopStreaming()

    expect(store.isStreaming).toBe(false)
  })

  it('clearMessages 应该清空消息列表', () => {
    const store = useChatStore()
    store.sendMessage('test')
    store.clearMessages()

    expect(store.messages.length).toBe(0)
  })

  it('loadMessages 应该从 localStorage 加载消息', () => {
    const sessionStore = useSessionStore()
    const sessionId = sessionStore.currentSessionId!
    const saved = [
      {
        id: 'msg-1',
        role: 'user',
        content: '你好',
        status: 'done',
        timestamp: 1000,
      },
    ]
    localStorage.setItem(`general-agent-messages-${sessionId}`, JSON.stringify(saved))

    const store = useChatStore()
    store.loadMessages()

    expect(store.messages.length).toBe(1)
    expect(store.messages[0].content).toBe('你好')
  })
})