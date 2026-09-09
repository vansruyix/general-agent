import type { IAdapter } from './types'
import { HttpAdapter } from './http'
import { WailsAdapter } from './wails'

/** 运行时检测：Wails 环境会注入 window.__WAILS__ */
const adapter: IAdapter = window.__WAILS__
  ? new WailsAdapter()
  : new HttpAdapter()

export default adapter