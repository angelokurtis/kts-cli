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

	c := containers.Names()
	p := containers.Pods()
	cmd := fmt.Sprintf("stern %s -c '%s' '%s' --since %s", ns, strings.Join(c, "|"), strings.Join(p, "|"), since)
	color.Secondary.Println(cmd)
}
