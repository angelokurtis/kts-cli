package cmd

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/spf13/cobra"
	"os"
)

var completion = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(kts completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(kts completion)
`,
	Run: func(_ *cobra.Command, _ []string) {
		if err := cmd.GenBashCompletion(os.Stdout); err != nil {
			common.Exit(err)
		}
	},
}
