package ingresses

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "ingresses",
	Short: "Utility functions for Kubernetes Ingresses",
	Run:   common.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "hosts", Run: hosts})
}
