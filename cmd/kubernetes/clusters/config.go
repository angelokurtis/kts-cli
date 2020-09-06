package clusters

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/aws"
	"github.com/spf13/cobra"
)

func config(cmd *cobra.Command, args []string) {
	cluster, err := aws.SelectEKSCluster()
	if err != nil {
		system.Exit(err)
	}
	err = aws.ConnectToEKSCluster(cluster)
	if err != nil {
		system.Exit(err)
	}
}
