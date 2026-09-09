import { describe, it, expect } from 'vitest'
import { generateUUID } from '../uuid'

describe('generateUUID', () => {
  it('应该返回字符串', () => {
    const id = generateUUID()
    expect(typeof id).toBe('string')
  })

  it('应该返回有效的 UUID 格式（36 字符，含 4 个连字符）', () => {
    const id = generateUUID()
    expect(id.length).toBe(36)
    expect(id.split('-').length).toBe(5)
  })

  it('每次调用应该返回不同的值', () => {
    const id1 = generateUUID()
    const id2 = generateUUID()
    expect(id1).not.toBe(id2)
  })
})