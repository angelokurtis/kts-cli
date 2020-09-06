package ecr

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/aws"
	"github.com/spf13/cobra"
	"strings"
)

var (
	Command = &cobra.Command{
		Use:   "ecr",
		Short: "Utility functions to deal with Amazon Elastic Container Registry (ECR)",
		Run:   system.Help,
	}
)

func init() {
	Command.AddCommand(&cobra.Command{Use: "images", Run: images})
}

func images(cmd *cobra.Command, args []string) {
	repos, err := aws.ListECRRepositories()
	if err != nil {
		log.Fatal(err)
	}
	repos, err = repos.SelectMany()
	if err != nil {
		log.Fatal(err)
	}
	for _, repo := range repos.Items {
		images, err := aws.ListECRImages(repo)
		if err != nil {
			log.Fatal(err)
		}
		for _, image := range images.Items {
			fmt.Println(image.Digest + " " + strings.Join(image.Tags, ", "))
		}
	}
}
