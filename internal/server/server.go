package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

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

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Println("Received message:\n", message)
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

func Start(cmd *cobra.Command, args []string) {
	log.Println("cmd", cmd.Use)
	log.Println("args", args)
	http.HandleFunc("/ws", wsHandler)
	log.Println("Websocket server started on :1804")
	if err := http.ListenAndServe(":1804", nil); err != nil {
		log.Println("Error starting server:", err)
	}
}
