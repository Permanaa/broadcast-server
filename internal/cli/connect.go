package cli

import (
	"github.com/permanaa/broadcast-server/internal/server"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect the client to the server.",
	Run:   server.Connect,
}

func init() {
	connectCmd.Flags().String("user", "", "username to show in text message")
	rootCmd.AddCommand(connectCmd)
}
