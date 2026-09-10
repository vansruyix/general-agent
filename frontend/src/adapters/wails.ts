import type { IAdapter } from './types'
import type { SSEEvent } from '../types/chat'
import { HttpAdapter } from './http'

/**
 * Wails v3 桌面适配器 — 浅集成模式
 * 复用 HttpAdapter 通过 localhost HTTP 与内嵌后端通信，
 * 与浏览器模式唯一区别是运行时检测标记 (window.__WAILS__)
 */
export class WailsAdapter implements IAdapter {
  private http: HttpAdapter = new HttpAdapter()

  chatStream(
    id: string,
    question: string,
    onEvent: (event: SSEEvent) => void,
    onError: (err: Error) => void,
    onDone: () => void
  ): AbortController {
    return this.http.chatStream(id, question, onEvent, onError, onDone)
  }
}