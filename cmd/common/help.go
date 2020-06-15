package common

import (
	"github.com/spf13/cobra"
)

func Help(cmd *cobra.Command, _ []string) {
	err := cmd.Help()
	if err != nil {
		Exit(err)
	}
}
