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
	rootCmd.AddCommand(connectCmd)
}
