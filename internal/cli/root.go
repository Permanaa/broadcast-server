package cli

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "broadcast-server",
	Short: "A server that can broadcast messages to connected clients.",
	Long:  "A simple broadcast server that will allow clients to connect to it, send messages that will be broadcasted to all connected clients.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println("Error execute root command:", err)
		os.Exit(1)
	}
}
