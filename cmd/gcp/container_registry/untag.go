package container_registry

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/gcloud"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func untag(_ *cobra.Command, _ []string) {
	tags, err := gcloud.SelectTags()
	color.Primary.Println("gcloud container images untag gcr.io/<PROJECT_ID>/<IMAGE_PATH>:<TAG>")
	err = gcloud.UntagImages(tags)
	if err != nil {
		system.Exit(err)
	}
}
