package main

import (
	"fmt"
	"github.com/PWND0U/dztool/ServerTool"
	"net/http"
	"time"
)

func main() {
	//fmt.Println("create pro")
	//fmt.Println(StringTool.NewDzString("1").ToString())
	http.HandleFunc("/sse", func(writer http.ResponseWriter, request *http.Request) {
		for {
			select {
			case <-request.Context().Done():
				fmt.Println("客户端离开")
				return
			default:
				fmt.Println("发送信息")
				event := ServerTool.NewDzServerSentEvent([]byte("Hello,World"), "update", "1", "哈哈哈哈", 1)
				event.SSEDataFlush(writer)
				time.Sleep(time.Second * 2)
			}
		}
	})
	http.ListenAndServe(":8888", nil)
}
