package nodes

import (
	"fmt"
	log "log/slog"
	"strings"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

// kube nodes selectors
func selectors(cmd *cobra.Command, args []string) {
	nodes, err := kubectl.ListNodes("")
	if err != nil {
		log.Error(err.Error())
		return
	}

	labels, err := nodes.SelectLabels()
	if err != nil {
		system.Exit(err)
	}

	var lb strings.Builder

	count := 0
	for key, values := range labels {
		if count > 0 {
			_, _ = fmt.Fprintf(&lb, ",")
		}

		if len(values) > 1 {
			_, _ = fmt.Fprintf(&lb, "%s in (", key)

			for i, value := range values {
				if i == 0 {
					_, _ = fmt.Fprintf(&lb, "%s", value)
				} else {
					_, _ = fmt.Fprintf(&lb, ",%s", value)
				}
			}

			_, _ = fmt.Fprintf(&lb, ")")
		} else {
			_, _ = fmt.Fprintf(&lb, "%s=%s", key, values[0])
		}

		count++
	}

	fmt.Println("-l \"" + lb.String() + "\"")
}
