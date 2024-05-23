package format

import (
	log "log/slog"

	"github.com/spf13/cobra"
)

type commandFunc func(cmd *cobra.Command, args []string) error

type simpleCommandFunc func(cmd *cobra.Command, args []string)

func wrapWithErrorHandler(fn commandFunc) simpleCommandFunc {
	return func(cmd *cobra.Command, args []string) {
		if err := fn(cmd, args); err != nil {
			log.Error(err.Error())
			return
		}
	}
}
