package server

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

func Connect(cmd *cobra.Command, args []string) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:1804/ws", nil)
	if err != nil {
		log.Fatalf("Error connect to server: %v", err)
	}
	defer conn.Close()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error read message:", err)
			}
			fmt.Println("Echo:", string(message))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("send message: ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading scanner:", err)
			}
			break
		}

		input := scanner.Text()
		input = strings.TrimSpace(input)

		if input == "" {
			fmt.Println("Message cannot be empty!")
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("Exiting application. Goodbye!")
			break
		}

		if err := conn.WriteMessage(websocket.TextMessage, []byte(input)); err != nil {
			fmt.Println("Error write message:", err)
			continue
		}
	}
}
