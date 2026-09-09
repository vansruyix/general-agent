import type { IAdapter } from './index'
import type { SSEEvent } from '../types/chat'

/**
 * Wails v3 适配器（占位实现）
 * 通过 Wails 的 Go 桥接调用后端，在 Wails 环境就绪前不做实际对接
 */
export class WailsAdapter implements IAdapter {
  chatStream(
    _id: string,
    _question: string,
    _onEvent: (event: SSEEvent) => void,
    onError: (err: Error) => void,
    _onDone: () => void
  ): AbortController {
    const controller = new AbortController()
    onError(new Error('Wails 适配器尚未实现，请使用浏览器模式'))
    return controller
  }
}