package format

import (
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/spf13/cobra"
)

type commandFunc func(cmd *cobra.Command, args []string) error

type simpleCommandFunc func(cmd *cobra.Command, args []string)

func wrapWithErrorHandler(fn commandFunc) simpleCommandFunc {
	return func(cmd *cobra.Command, args []string) {
		if err := fn(cmd, args); err != nil {
			log.Fatal(err)
		}
	}
}
