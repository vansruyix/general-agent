package memory

import (
	"sync"

	"github.com/cloudwego/eino/schema"
)

// 全局会话存储（内存 map）
var SimpleMemoryMap = &sync.Map{}

type SimpleMemory struct {
	ID       string            // 会话ID
	Messages []*schema.Message // 消息历史
}

func NewMemory(id string) (*SimpleMemory, error) {
	var simpleMemory *SimpleMemory
	simpleMemory = &SimpleMemory{
		ID:       id,
		Messages: []*schema.Message{},
	}
	SimpleMemoryMap.Store(id, simpleMemory)
	return simpleMemory, nil
}

func AppendMessage(id string, msg *schema.Message) {
	memory, _ := SimpleMemoryMap.Load(id)
	if memory == nil {
		memory, _ = NewMemory(id)
	}
	mem := memory.(*SimpleMemory)
	mem.Messages = append(mem.Messages, msg)
	SimpleMemoryMap.Store(id, mem)
}

func GetMemory(id string) *SimpleMemory {
	memory, _ := SimpleMemoryMap.Load(id)
	if memory == nil {
		memory, _ = NewMemory(id)
		SimpleMemoryMap.Store(id, memory)
		return memory.(*SimpleMemory)
	}
	return memory.(*SimpleMemory)
}

func Clear(id string) error {
	SimpleMemoryMap.Delete(id)
	return nil
}
