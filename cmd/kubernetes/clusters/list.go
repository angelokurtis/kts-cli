package clusters

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/aws"
	"github.com/angelokurtis/kts-cli/pkg/app/gcloud"
)

func list(cmd *cobra.Command, args []string) {
	awsc, err := aws.ListEKSClusters()
	if err != nil {
		color.Yellow.Println("[WARN] " + err.Error())
	}

	gc, err := gcloud.ListGKEClustersNames()
	if err != nil {
		color.Yellow.Println("[WARN] " + err.Error())
	}

	clusters := append(awsc, gc...)
	if len(clusters) == 0 {
		log.Debug("no clusters where found")
	} else {
		for _, cluster := range clusters {
			fmt.Println(cluster)
		}
	}
}
