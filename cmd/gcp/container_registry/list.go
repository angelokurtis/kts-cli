package container_registry

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/pkg/app/gcloud"
	"github.com/cheggaaa/pb/v3"
	"github.com/gookit/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

func list(_ *cobra.Command, _ []string) {
	color.Comment.Println("gcloud container images list --repository=gcr.io/<PROJECT_ID>")
	repositories, err := gcloud.SelectContainerRepositories()
	if err != nil {
		common.Exit(err)
	}

	images := make([]*gcloud.ContainerImage, 0, 0)
	if len(repositories) > 0 {
		color.Comment.Println("gcloud container images list-tags gcr.io/<PROJECT_ID>/<IMAGE_PATH>")
		bar := pb.StartNew(len(repositories))
		for _, repository := range repositories {
			img, err := gcloud.ListContainerImages(repository)
			if err != nil {
				common.Exit(err)
			}
			images = append(images, img...)
			bar.Increment()
		}
		bar.Finish()
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Timestamp", "Repository", "Tags", "Digest"})
	table.SetBorder(false)
	for _, image := range images {
		if len(image.Tags) > 0 {
			for _, tag := range image.Tags {
				table.Append([]string{image.Timestamp.Datetime.Format("02-01-2006 15:04"), image.Repository, tag, image.Digest})
			}
		} else {
			table.Append([]string{image.Timestamp.Datetime.Format("02-01-2006 15:04"), image.Repository, "", image.Digest})
		}
	}
	table.Render()
}
