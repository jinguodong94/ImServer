package socket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func StartServer(addr string) {

	log.Println("启动websocket服务")

	ClientMgr.Init()

	http.HandleFunc("/websocket", func(writer http.ResponseWriter, request *http.Request) {
		// 升级协议
		conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(writer, request, nil)
		if err != nil {
			http.NotFound(writer, request)
			return
		}

		client := NewClient(conn.RemoteAddr().String(), conn)

		client.Open()
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
