package agent

import (
	"testing"

	"github.com/google/uuid"
)

func TestDeepAgents(t *testing.T) {
	NewSimpleDeepAgent(uuid.NewString(), "当前项目目录的结构是什么样的,全程使用中文交流")
}
