package main

import (
	"bufio"
	"log"
	"net/http"
)

// HTTP sample client
// 1. text/event-stream

func callEventStream(client *http.Client, url string) {
	// body 的流式写入见 https://segmentfault.com/a/1190000022814356
	rep, err := client.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(rep.Body)
	defer rep.Body.Close()

	for {
		// handle chunked data
		data, err := reader.ReadString('\n')
		// ignore the second \n
		_, _ = reader.ReadString('\n')

		// other error
		if err != nil {
			log.Printf("[Error] read error: %v", err)
			continue
		}
		// read data
		log.Printf("receive: [%s]", data)
	}
}

func main() {
	client := &http.Client{}
	rootURL := "http://127.0.0.1:8000/"

	callEventStream(client, rootURL+"event_stream")
}
