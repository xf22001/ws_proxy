package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// 模拟另一个WebSocket连接
	remoteConn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8389/ws", nil)
	if err != nil {
		log.Println("Error dialing remote websocket:", err)
		return
	}
	defer remoteConn.Close()

	// 转发流量
	go func() {
		defer remoteConn.Close()
		defer conn.Close()
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}
			err = remoteConn.WriteMessage(messageType, p)
			if err != nil {
				log.Println("Error writing message:", err)
				return
			}
		}
	}()

	for {
		messageType, p, err := remoteConn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			log.Println("Error writing message:", err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", websocketHandler)
	fmt.Println("Starting server on port 8390...")
	log.Fatal(http.ListenAndServe(":8390", nil))
}
