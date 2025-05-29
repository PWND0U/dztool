package ServerTool

import (
	"fmt"
	"testing"
)

func TestDzServerSentEvent(t *testing.T) {
	sse := DecodeDzServerSentEvent([]byte("event: update\r\ndata: test\r\n\r\n"))
	fmt.Println(sse.Event.ToString())
	fmt.Println(string(sse.Data))
}
