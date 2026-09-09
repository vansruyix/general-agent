import type { SSEEvent } from '../types/chat'
import { HttpAdapter } from './http'
import { WailsAdapter } from './wails'

/** 适配器接口：统一浏览器和桌面端的通信方式 */
export interface IAdapter {
  /**
   * 发送聊天消息，通过 SSE 流式接收响应
   * @returns AbortController，调用方可通过 .abort() 取消请求
   */
  chatStream(
    id: string,
    question: string,
    onEvent: (event: SSEEvent) => void,
    onError: (err: Error) => void,
    onDone: () => void
  ): AbortController
}

/** 运行时检测：Wails 环境会注入 window.__WAILS__ */
const adapter: IAdapter = window.__WAILS__
  ? new WailsAdapter()
  : new HttpAdapter()

export default adapter