package merge_requests

import (
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/gitlab"
	"github.com/spf13/cobra"
)

// cartaobranco projects merge-requests assign
func assign(cmd *cobra.Command, args []string) {
	mrs, err := gitlab.SearchMergeRequestsByUser(username)
	if err != nil {
		log.Fatal(err)
	}

	mr, err := mrs.SelectOne()
	if err != nil {
		log.Fatal(err)
	}

	err = mr.Assign()
	if err != nil {
		log.Fatal(err)
	}
}
