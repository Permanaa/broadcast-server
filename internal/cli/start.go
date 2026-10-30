package cli

import (
	"github.com/permanaa/broadcast-server/internal/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server.",
	Run:   server.Start,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
