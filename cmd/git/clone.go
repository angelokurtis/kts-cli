package git

import (
	log "log/slog"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/git"
	"github.com/angelokurtis/kts-cli/pkg/app/idea"
)

func clone(cmd *cobra.Command, args []string) {
	repo := args[0]

	dir, err := git.NewLocalDir(repo)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if !dir.Exist() {
		log.Info("cloning repo", "repo", repo)

		err = git.Clone(repo)
		if err != nil {
			log.Error(err.Error())
			return
		}
	} else {
		log.Info("local repository was found at the path", "path", dir.Path())
	}

	if open {
		log.Info("opening on IntelliJ IDEA", "path", dir.Path())

		err = idea.Open(dir.Path())
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
}
