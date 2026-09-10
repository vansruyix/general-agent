import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Session } from '../types/session'
import { generateUUID } from '../utils/uuid'

const STORAGE_KEY = 'general-agent-sessions'
export const MESSAGES_PREFIX = 'general-agent-messages-'

export const useSessionStore = defineStore('session', () => {
  const sessions = ref<Session[]>([])
  const currentSessionId = ref<string | null>(null)

  /** 从 localStorage 加载会话列表 */
  function loadSessions() {
    try {
      const raw = localStorage.getItem(STORAGE_KEY)
      if (raw) {
        sessions.value = JSON.parse(raw)
        // 恢复上次的会话（取第一个即最近创建的会话）
        if (sessions.value.length > 0) {
          currentSessionId.value = sessions.value[0].id
        }
      }
    } catch (e) {
      console.warn('加载会话列表失败:', e)
    }
  }

  /** 持久化会话列表到 localStorage */
  function saveSessions() {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(sessions.value))
    } catch (e) {
      console.warn('保存会话列表失败:', e)
    }
  }

  /** 创建新会话并切换到它 */
  function createSession(): Session {
    const session: Session = {
      id: generateUUID(),
      title: '新会话',
      createdAt: Date.now(),
      updatedAt: Date.now(),
      messageCount: 0,
    }
    sessions.value.unshift(session)
    currentSessionId.value = session.id
    saveSessions()
    return session
  }

  /** 删除会话及其关联的消息数据 */
  function deleteSession(id: string) {
    const index = sessions.value.findIndex((s) => s.id === id)
    if (index === -1) return

    sessions.value.splice(index, 1)
    try {
      localStorage.removeItem(`${MESSAGES_PREFIX}${id}`)
    } catch (e) {
      console.warn('删除会话消息失败:', e)
    }

    if (currentSessionId.value === id) {
      currentSessionId.value = sessions.value[0]?.id ?? null
    }
    saveSessions()
  }

  /** 切换到指定会话 */
  function switchSession(id: string) {
    currentSessionId.value = id
  }

  /** 更新会话标题（取第一条消息前 30 字符） */
  function updateSessionTitle(id: string, title: string) {
    const session = sessions.value.find((s) => s.id === id)
    if (session) {
      session.title = title
      session.updatedAt = Date.now()
      saveSessions()
    }
  }

  /** 递增会话消息计数 */
  function incrementMessageCount(id: string) {
    const session = sessions.value.find((s) => s.id === id)
    if (session) {
      session.messageCount++
      session.updatedAt = Date.now()
      saveSessions()
    }
  }

  // 初始化时从 localStorage 加载
  loadSessions()

  return {
    sessions,
    currentSessionId,
    loadSessions,
    createSession,
    deleteSession,
    switchSession,
    updateSessionTitle,
    incrementMessageCount,
  }
})