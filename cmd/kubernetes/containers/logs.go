package containers

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

// kube containers logs -s 2h
func logs(cmd *cobra.Command, args []string) {
	containers, err := kubectl.ListContainers(namespace, allNamespaces, "")
	if err != nil {
		system.Exit(err)
	}

	containers, err = containers.SelectMany()
	if err != nil {
		system.Exit(err)
	}

	if download {
		kubectl.SaveLogs(containers, since, previous)
	} else {
		stern(containers, since)
	}
}

func stern(containers *kubectl.Containers, since string) {
	ns := "--all-namespaces"
	namespaces := containers.Namespaces()

	if len(namespaces) == 1 {
		ns = fmt.Sprintf("-n %s", namespaces[0])
	}

	cmd := []string{"stern", ns}
	cmd = append(cmd, fmt.Sprintf("-c '^(%s)$'", strings.Join(containers.Names(), "|")))
	cmd = append(cmd, fmt.Sprintf("'^(%s)$'", strings.Join(containers.Pods(), "|")))

	if since != "0s" {
		cmd = append(cmd, fmt.Sprintf("--since %s", since))
	}

	color.Secondary.Println(strings.Join(cmd, " "))
}
