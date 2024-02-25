package deployments

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

// kube deployments logs -s 2h
func logs(cmd *cobra.Command, args []string) {
	deploys, err := kubectl.ListDeployments(namespace, allNamespaces)
	if err != nil {
		log.Fatal(err)
	}

	deploy, err := deploys.SelectOne()
	if err != nil {
		log.Fatal(err)
	}

	if download {
		containers, err := kubectl.ListContainersByDeployment(deploy)
		if err != nil {
			log.Fatal(err)
		}

		kubectl.SaveLogs(containers, since, previous)
	} else {
		follow(deploy, since)
	}
}

func follow(deploy *kubectl.Deployment, since string) {
	ns := fmt.Sprintf("-n %s", deploy.Metadata.Namespace)

	labels := make([]string, 0, 0)
	for key, value := range deploy.Spec.Selector.MatchLabels {
		labels = append(labels, key+"="+value)
	}

	if since == "0s" {
		since = ""
	} else {
		since = "--since " + since
	}

	cmd := fmt.Sprintf("stern %s -l %s %s", ns, strings.Join(labels, ","), since)
	color.Secondary.Println(cmd)
}
