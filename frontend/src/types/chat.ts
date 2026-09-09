/** SSE 事件类型，与后端 ai/sse/event.go 的 EventType 对齐 */
export type SSEEventType =
  | 'thinking'
  | 'text'
  | 'tool_call'
  | 'tool_result'
  | 'image'
  | 'audio'
  | 'video'
  | 'file'
  | 'error'
  | 'done'

/** 后端 SSE 流式响应中的单条事件 */
export interface SSEEvent {
  type: SSEEventType
  data: string
  mime_type?: string
  base64?: string
  url?: string
  meta?: Record<string, unknown>
}

/** 工具调用记录 */
export interface ToolCall {
  callId: string
  toolName: string
  params: string
  result?: string
}

/** 消息状态 */
export type MessageStatus = 'sending' | 'streaming' | 'done' | 'error'

/** 聊天消息 */
export interface ChatMessage {
  id: string
  role: 'user' | 'assistant'
  content: string
  thinking?: string
  toolCalls?: ToolCall[]
  status: MessageStatus
  timestamp: number
}