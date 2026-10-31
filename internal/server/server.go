package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var clients = make(map[*websocket.Conn]bool)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			delete(clients, conn)
			break
		}
		log.Println("Received message:", string(message))
		for clientConn := range clients {
			if err := clientConn.WriteMessage(websocket.TextMessage, message); err != nil {
				clientConn.Close()
				delete(clients, conn)
				break
			}
		}
	}
}

func Start(cmd *cobra.Command, args []string) {
	http.HandleFunc("/ws", wsHandler)
	log.Println("Websocket server started on :1804")
	if err := http.ListenAndServe(":1804", nil); err != nil {
		log.Println("Error starting server:", err)
	}
}
