package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var completion = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(kts completion bash)

# To load completions for each session, execute once:
Linux:
  $ kts completion bash > /etc/bash_completion.d/kts
MacOS:
  $ kts completion bash > /usr/local/etc/bash_completion.d/kts

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ kts completion zsh > "${fpath[1]}/_kts"

# You will need to start a new shell for this setup to take effect.

Fish:

$ kts completion fish | source

# To load completions for each session, execute once:
$ kts completion fish > ~/.config/fish/completions/kts.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}
