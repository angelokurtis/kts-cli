package container_registry

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/internal/color"
	"github.com/angelokurtis/kts-cli/pkg/app/gcloud"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

func clean(_ *cobra.Command, _ []string) {
	fmt.Printf(color.Notice, "gcloud container images list\n")
	repositories, err := gcloud.SelectContainerRepositories()
	if err != nil {
		common.Exit(err)
	}

	images := make([]*gcloud.ContainerImage, 0, 0)
	if len(repositories) > 0 {
		fmt.Printf(color.Notice, "gcloud container images list-tags gcr.io/<PROJECT_ID>/<IMAGE_PATH> --filter=\"NOT tags:*\"\n")
		tagBar := pb.StartNew(len(repositories))
		for _, repository := range repositories {
			img, err := gcloud.ListContainerImagesWithoutTags(repository)
			if err != nil {
				common.Exit(err)
			}
			images = append(images, img...)
			tagBar.Increment()
		}
		tagBar.Finish()
	}

	if len(images) > 0 {
		fmt.Printf(color.Notice, "gcloud container images delete gcr.io/<PROJECT_ID>/<IMAGE_PATH>@<DIGEST>\n")
		delBar := pb.StartNew(len(images))
		for _, image := range images {
			err := gcloud.DeleteContainerImage(image)
			if err != nil {
				common.Exit(err)
			}
			delBar.Increment()
		}
		delBar.Finish()
	}
}
