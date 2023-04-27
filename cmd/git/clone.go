package git

import (
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/git"
	"github.com/angelokurtis/kts-cli/pkg/app/idea"
	"github.com/spf13/cobra"
)

func clone(cmd *cobra.Command, args []string) {
	repo := args[0]
	dir, err := git.NewLocalDir(repo)
	if err != nil {
		log.Fatal(err)
	}
	if !dir.Exist() {
		log.Infof("cloning repo %s", repo)
		err = git.Clone(repo)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Infof("local repository was found at the path %s", dir.Path())
	}
	if open {
		log.Infof("opening %s on IntelliJ IDEA", dir.Path())
		err = idea.Open(dir.Path())
		if err != nil {
			log.Fatal(err)
		}
	}
}
