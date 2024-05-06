package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Http sample server
// 1. octed_stream
// 2. text/event-stream

func registFunc(mux *http.ServeMux) {
	mux.HandleFunc("/octed_stream", octedStream)
	mux.HandleFunc("/event_stream", eventStream)
}

// for production, can use https://github.com/alexandrevicenzi/go-sse
func eventStream(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	if "GET" == method {
		log.Printf("handle eventStream begin")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		flusher, ok := w.(http.Flusher)
		if !ok {
			log.Panic("server not support")
		}

		// TODO: should take over the http connection and keep it in memory
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		// 若通道为空，则阻塞
		// 若通道有数据，则读取
		// 若通道关闭，则退出
		for range ticker.C {
			// 事件字符串应该遵循text/event-stream数据格式的结构，可以被 前端 SSE 框架读取
			_, err := fmt.Fprint(w, "data: hello, user!\n\n")
			// occurs error, break the ticket
			if err != nil {
				log.Printf("[Error] write data error %v", err)
				break
			}
			// 刷新数据，以通知请求端
			flusher.Flush()
		}
		log.Println("handle eventStream over")
		return
	}

	w.Write([]byte("404 not found"))
}

func octedStream(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	if "GET" == method {
		log.Printf("handle octedStream")
		w.Header().Set("Content-Type", "application/octed-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\"picture.png\"")
		w.Write([]byte("hello, user!"))
		return
	}

	w.Write([]byte("404 not found"))
}

func main() {
	server := &http.Server{
		Addr: "0.0.0.0:8000",
		// No Write Timeout and Idle Timeout as server should keep the connection alive
	}
	mux := http.NewServeMux()
	registFunc(mux)
	server.Handler = mux
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("server start error %v", err)
	}
}
