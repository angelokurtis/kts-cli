package branch

import (
	log "log/slog"

	"github.com/spf13/cobra"
)

// git branch list
func list(_ *cobra.Command, _ []string) {
}

func check(err error) {
	if err != nil {
		log.Error(err.Error())
		return
	}
}
