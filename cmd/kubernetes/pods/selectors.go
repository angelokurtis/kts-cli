package pods

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/spf13/cobra"
	"strings"
)

// kube pods selectors
func selectors(cmd *cobra.Command, args []string) {
	pods, err := kubectl.ListPods(namespace, allNamespaces, "")
	if err != nil {
		log.Fatal(err)
	}
	labels, err := pods.SelectLabels()
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

	if namespace != "" {
		fmt.Println("-l \"" + lb.String() + "\" -n " + namespace)
	} else {
		fmt.Println("-l \"" + lb.String() + "\"")
	}
}
