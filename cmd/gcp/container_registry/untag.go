package container_registry

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/internal/color"
	"github.com/angelokurtis/kts-cli/pkg/app/gcloud"
	"github.com/spf13/cobra"
)

func untag(_ *cobra.Command, _ []string) {
	tags, err := gcloud.SelectTags()
	fmt.Printf(color.Notice, "gcloud container images untag gcr.io/<PROJECT_ID>/<IMAGE_PATH>:<TAG>\n")
	err = gcloud.UntagImages(tags)
	if err != nil {
		common.Exit(err)
	}
}
