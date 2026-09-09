import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSessionStore } from '../session'

describe('sessionStore', () => {
  beforeEach(() => {
    // 每个测试前重置 Pinia 和 localStorage
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('createSession 应该创建新会话并设为当前会话', () => {
    const store = useSessionStore()
    const session = store.createSession()

    expect(session.id).toBeTruthy()
    expect(session.title).toBe('新会话')
    expect(store.currentSessionId).toBe(session.id)
    expect(store.sessions.length).toBe(1)
  })

  it('deleteSession 应该删除会话并清除关联消息', () => {
    const store = useSessionStore()
    const session = store.createSession()
    // 模拟存储消息
    localStorage.setItem(`general-agent-messages-${session.id}`, '[]')

    store.deleteSession(session.id)

    expect(store.sessions.length).toBe(0)
    expect(store.currentSessionId).toBeNull()
    expect(localStorage.getItem(`general-agent-messages-${session.id}`)).toBeNull()
  })

  it('switchSession 应该切换当前会话', () => {
    const store = useSessionStore()
    const s1 = store.createSession()
    const s2 = store.createSession()

    store.switchSession(s1.id)
    expect(store.currentSessionId).toBe(s1.id)

    store.switchSession(s2.id)
    expect(store.currentSessionId).toBe(s2.id)
  })

  it('loadSessions 应该从 localStorage 恢复会话列表', () => {
    const saved = [
      { id: 'test-1', title: '测试会话', createdAt: 1000, updatedAt: 2000, messageCount: 3 }
    ]
    localStorage.setItem('general-agent-sessions', JSON.stringify(saved))

    // 创建新 store 实例会自动调用 loadSessions
    const store = useSessionStore()
    store.loadSessions()

    expect(store.sessions.length).toBe(1)
    expect(store.sessions[0].id).toBe('test-1')
  })
})