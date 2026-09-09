import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { ChatMessage, SSEEvent, ToolCall } from '../types/chat'
import { generateUUID } from '../utils/uuid'
import adapter from '../adapters'
import { useSessionStore, MESSAGES_PREFIX } from './session'

export const useChatStore = defineStore('chat', () => {
  const messages = ref<ChatMessage[]>([])
  const isStreaming = ref(false)
  const error = ref<string | null>(null)

  let abortController: AbortController | null = null

  const sessionStore = useSessionStore()

  /** 从 localStorage 加载当前会话的消息 */
  function loadMessages() {
    const sessionId = sessionStore.currentSessionId
    if (!sessionId) {
      messages.value = []
      return
    }
    try {
      const raw = localStorage.getItem(`${MESSAGES_PREFIX}${sessionId}`)
      messages.value = raw ? JSON.parse(raw) : []
    } catch (e) {
      console.warn('加载消息失败:', e)
      messages.value = []
    }
  }

  /** 持久化当前会话消息到 localStorage */
  function saveMessages() {
    const sessionId = sessionStore.currentSessionId
    if (!sessionId) return
    try {
      localStorage.setItem(`${MESSAGES_PREFIX}${sessionId}`, JSON.stringify(messages.value))
    } catch (e) {
      console.warn('保存消息失败:', e)
    }
  }

  /** 发送消息并启动 SSE 流式接收 */
  function sendMessage(question: string) {
    const sessionId = sessionStore.currentSessionId
    if (!sessionId || !question.trim()) return

    // 如果正在接收流式响应，先停止旧流
    if (isStreaming.value) stopStreaming()

    error.value = null

    // 第一条消息：用问题作为会话标题
    if (messages.value.length === 0) {
      sessionStore.updateSessionTitle(sessionId, question.slice(0, 30))
    }

    // 添加用户消息
    const userMsg: ChatMessage = {
      id: generateUUID(),
      role: 'user',
      content: question,
      status: 'done',
      timestamp: Date.now(),
    }
    messages.value.push(userMsg)

    // 创建 AI 消息占位
    const aiMsg: ChatMessage = {
      id: generateUUID(),
      role: 'assistant',
      content: '',
      status: 'streaming',
      timestamp: Date.now(),
    }
    messages.value.push(aiMsg)

    isStreaming.value = true
    saveMessages()
    sessionStore.incrementMessageCount(sessionId)

    abortController = adapter.chatStream(
      sessionId,
      question,
      // onEvent: 分发 SSE 事件到消息
      (event: SSEEvent) => handleSSEEvent(aiMsg, event),
      // onError: 标记错误
      (err: Error) => {
        aiMsg.status = 'error'
        aiMsg.content = `请求失败: ${err.message}`
        error.value = err.message
        isStreaming.value = false
        saveMessages()
      },
      // onDone: 标记完成
      () => {
        if (aiMsg.status === 'streaming') {
          aiMsg.status = 'done'
        }
        isStreaming.value = false
        abortController = null
        saveMessages()
      }
    )
  }

  /** 处理单条 SSE 事件，更新 AI 消息状态 */
  function handleSSEEvent(msg: ChatMessage, event: SSEEvent) {
    switch (event.type) {
      case 'thinking':
        msg.thinking = (msg.thinking || '') + event.data
        break
      case 'text':
        msg.content += event.data
        break
      case 'tool_call': {
        if (!msg.toolCalls) msg.toolCalls = []
        const tc: ToolCall = {
          callId: (event.meta?.call_id as string) || '',
          toolName: (event.meta?.tool_name as string) || '未知工具',
          params: event.data || '',
        }
        msg.toolCalls.push(tc)
        break
      }
      case 'tool_result': {
        const callId = event.meta?.call_id as string
        if (msg.toolCalls && callId) {
          const tc = msg.toolCalls.find((t) => t.callId === callId)
          if (tc) tc.result = event.data || ''
        }
        break
      }
      case 'error':
        msg.content += `\n\n❌ 错误: ${event.data}`
        msg.status = 'error'
        error.value = event.data
        break
      case 'done':
        msg.status = 'done'
        break
      // image, audio, video, file 暂不做特殊处理
      default:
        break
    }
    saveMessages()
  }

  /** 停止当前 SSE 流 */
  function stopStreaming() {
    if (abortController) {
      abortController.abort()
      abortController = null
    }
    isStreaming.value = false
    // 将最后一条 assistant 消息从 streaming 置为 done，避免 UI 持续显示加载动画
    const lastMsg = messages.value[messages.value.length - 1]
    if (lastMsg && lastMsg.role === 'assistant' && lastMsg.status === 'streaming') {
      lastMsg.status = 'done'
      saveMessages()
    }
  }

  /** 清空当前会话消息 */
  function clearMessages() {
    messages.value = []
    saveMessages()
  }

  return {
    messages,
    isStreaming,
    error,
    loadMessages,
    sendMessage,
    stopStreaming,
    clearMessages,
  }
})