package versions

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use: "versions",
	Run: versions,
}
