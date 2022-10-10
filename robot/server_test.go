package robot

import (
	"testing"

	"after/chat"
)

func TestResponse(t *testing.T) {
	chat.InitChat("../chat/corpus.gob", 10)
	t.Log(chat.GetResponse("你好鸭！"))
}
