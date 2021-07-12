package packages

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use: "packages",
	Run: packages,
}
