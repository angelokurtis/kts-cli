package container_registry

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/pkg/app/gcloud"
	"github.com/cheggaaa/pb/v3"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func clean(_ *cobra.Command, _ []string) {
	color.Comment.Println("gcloud container images list")
	repositories, err := gcloud.SelectContainerRepositories()
	if err != nil {
		common.Exit(err)
	}

	images := make([]*gcloud.ContainerImage, 0, 0)
	if len(repositories) > 0 {
		color.Comment.Println("gcloud container images list-tags gcr.io/<PROJECT_ID>/<IMAGE_PATH> --filter=\"NOT tags:*\"")
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
		color.Primary.Println("gcloud container images delete gcr.io/<PROJECT_ID>/<IMAGE_PATH>@<DIGEST>")
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
