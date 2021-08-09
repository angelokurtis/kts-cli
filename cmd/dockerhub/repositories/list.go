package repositories

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
)

func list(cmd *cobra.Command, args []string) {
	dockerhub := newDockerhubClient()
	repos, err := dockerhub.ListRepositories("kurtis")
	if err != nil {
		log.Fatal(err)
	}

	for _, repo := range repos {
		fmt.Println(repo.Name)
	}
}
