package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// websocket server demo
// for more usage, see https://github.com/gorilla/websocket/tree/main/examples/chat

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	// use upgrade field to distinguish http or websocket request
	upgrade := r.Header.Get("Upgrade")
	if len(upgrade) == 0 {
		log.Printf("not websocket, use http")
		w.Write([]byte("it is a normal http request"))
		return
	}

	log.Printf("upgrade is %s, using websocket", upgrade)
	// 升级HTTP连接为WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// 处理WebSocket连接
	for {
		// 读取消息
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error", err)
			return
		}
		log.Println("Received message:", string(p))

		// 发送消息
		err = conn.WriteMessage(messageType, []byte("Hello, world!"))
		if err != nil {
			log.Println("write error", err)
			return
		}
	}
}

func main() {
	// 创建HTTP服务器
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
