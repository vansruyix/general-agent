package sse

// Consumer 是 SSEvent 的消费函数类型。
// 调用方通过传入不同的 Consumer 实现来决定 SSEvent 的消费方式：
//   - Web 场景：将 event 序列化为 JSON 后写入 HTTP SSE 响应
//   - 控制台场景：将 event 格式化后打印到终端
//
// 返回 error 时，调用方应停止后续事件的消费。
type Consumer func(event SSEvent) error
