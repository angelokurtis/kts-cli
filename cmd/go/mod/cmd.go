package mod

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use: "upgrade",
	Run: upgrade,
}
