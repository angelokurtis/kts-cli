package container_registry

import (
	pb "github.com/cheggaaa/pb/v3"
	"github.com/gookit/color"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/gcloud"
)

func clean(_ *cobra.Command, _ []string) {
	color.Secondary.Println("gcloud container images list")

	repositories, err := gcloud.SelectContainerRepositories()
	if err != nil {
		system.Exit(err)
	}

	images := make([]*gcloud.ContainerImage, 0, 0)

	if len(repositories) > 0 {
		color.Secondary.Println("gcloud container images list-tags gcr.io/<PROJECT_ID>/<IMAGE_PATH> --filter=\"NOT tags:*\"")

		tagBar := pb.StartNew(len(repositories))

		for _, repository := range repositories {
			img, err := gcloud.ListContainerImagesWithoutTags(repository)
			if err != nil {
				system.Exit(err)
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
				system.Exit(err)
			}

			delBar.Increment()
		}

		delBar.Finish()
	}
}
