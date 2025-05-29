package ServerTool

import (
	"bytes"
	"fmt"
	"github.com/PWND0U/dztool/StringTool"
	"net/http"
	"regexp"
)

var (
	lineSepExpr, _ = regexp.Compile("\r\n|\r|\n")
)

type DzServerSentEvent struct {
	Data             []byte
	Event            StringTool.DzString
	Id               StringTool.DzString
	Retry            int
	Comment          StringTool.DzString
	DefaultSeparator string
}

func NewDzServerSentEvent(data []byte, event, id, comment string, retry int) *DzServerSentEvent {
	return &DzServerSentEvent{
		Data:             data,
		Event:            StringTool.NewDzString(event),
		Id:               StringTool.NewDzString(id),
		Retry:            retry,
		Comment:          StringTool.NewDzString(comment),
		DefaultSeparator: "\r\n",
	}
}

func (dSse *DzServerSentEvent) Encode() []byte {
	buffer := bytes.NewBuffer(nil)
	if !dSse.Comment.IsEmpty() {
		for _, s := range lineSepExpr.Split(dSse.Comment.ToString(), -1) {
			buffer.WriteString(fmt.Sprintf(": %s%s", s, dSse.DefaultSeparator))
		}
	}

	if !dSse.Id.IsEmpty() {
		buffer.WriteString(fmt.Sprintf("id: %s%s", string(lineSepExpr.ReplaceAll([]byte(dSse.Id), []byte(""))), dSse.DefaultSeparator))
	}

	if !dSse.Event.IsEmpty() {
		buffer.WriteString(fmt.Sprintf("event: %s%s", string(lineSepExpr.ReplaceAll([]byte(dSse.Event), []byte(""))), dSse.DefaultSeparator))
	}

	if len(dSse.Data) > 0 {
		for _, s := range lineSepExpr.Split(string(dSse.Data), -1) {
			buffer.WriteString(fmt.Sprintf("data: %s%s", s, dSse.DefaultSeparator))
		}
	}

	if dSse.Retry > 0 {
		buffer.WriteString(fmt.Sprintf("retry: %d%s", dSse.Retry, dSse.DefaultSeparator))
	}
	buffer.WriteString(dSse.DefaultSeparator)
	return buffer.Bytes()
}

func DecodeDzServerSentEvent(b []byte) *DzServerSentEvent {

	return nil
}

func (dSse *DzServerSentEvent) SSEDataFlush(resp http.ResponseWriter) bool {
	resp.Header().Set("Content-Type", "text/event-stream")
	resp.Header().Set("Cache-Control", "no-cache")
	resp.Header().Set("Connection", "keep-alive")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	write, err := resp.Write(dSse.Encode())
	if err == nil && write > 0 {
		resp.(http.Flusher).Flush()
		return true
	}
	return false
}
