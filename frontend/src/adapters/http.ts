import type { IAdapter } from './types'
import type { SSEEvent } from '../types/chat'

const API_BASE = '/api/v1/general-agent'

export class HttpAdapter implements IAdapter {
  chatStream(
    id: string,
    question: string,
    onEvent: (event: SSEEvent) => void,
    onError: (err: Error) => void,
    onDone: () => void
  ): AbortController {
    const controller = new AbortController()

    fetch(`${API_BASE}/agent/chat`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id, question }),
      signal: controller.signal,
    })
      .then(async (response) => {
        if (!response.ok) {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`)
        }
        const reader = response.body?.getReader()
        if (!reader) {
          throw new Error('Response body is not readable')
        }

        const decoder = new TextDecoder()
        let buffer = ''

        while (true) {
          const { done, value } = await reader.read()
          if (done) break

          buffer += decoder.decode(value, { stream: true })
          const lines = buffer.split('\n')
          buffer = lines.pop() || ''

          for (const line of lines) {
            if (line.startsWith('data: ')) {
              const jsonStr = line.slice(6).trim()
              if (!jsonStr) continue
              try {
                const event: SSEEvent = JSON.parse(jsonStr)
                onEvent(event)
              } catch {
                console.warn('SSE 解析错误，跳过:', jsonStr)
              }
            }
          }
        }
        onDone()
      })
      .catch((err: Error) => {
        if (err.name === 'AbortError') {
          onDone()
        } else {
          onError(err)
        }
      })

    return controller
  }
}