package istio

import (
	"fmt"
	log "log/slog"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/kiali"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

// istio relationships
func relationships(cmd *cobra.Command, args []string) {
	var nss []string

	if namespace == "" {
		namespaces, err := kubectl.ListNamespaces()
		if err != nil {
			log.Error(err.Error())
			return
		}

		namespaces, err = namespaces.SelectMany()
		if err != nil {
			log.Error(err.Error())
			return
		}

		nss = make([]string, 0, len(namespaces.Items))
		for _, ns := range namespaces.Items {
			nss = append(nss, ns.Metadata.Name)
		}
	} else {
		nss = []string{namespace}
	}

	graph, err := kiali.LoadGraphInfo(nss...)
	if err != nil {
		log.Error(err.Error())
		return
	}

	nodes := graph.GetNodes()

	node, err := nodes.SelectOne()
	if err != nil {
		log.Error(err.Error())
		return
	}

	res := graph.Inbound(node).Join(graph.Outbound(node))
	i := 0

	for _, re := range res {
		if i > 0 {
			fmt.Print(" AND ")
		}

		fmt.Print(re.Selector())

		i++
	}

	fmt.Println()
}
