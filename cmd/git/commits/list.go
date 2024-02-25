package commits

import (
	"fmt"
	"os"
	"time"

	prettytime "github.com/andanhm/go-prettytime"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/git"
)

var brazil *time.Location

func init() {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatal(err)
	}

	brazil = loc
}

// git commits list
func list(_ *cobra.Command, _ []string) {
	commits, err := git.ListCommits(dir)
	check(err)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"COMMIT", "CREATED", "SIGNED", "SIGNER"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetColWidth(100)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	for _, commit := range commits {
		commitTime := commit.Time
		t := fmt.Sprintf("%s (%s)", commitTime.In(brazil).Format("02/01/2006 15:04"), prettytime.Format(commitTime))
		signed := fmt.Sprintf("%v %s", commit.VerificationStatus(), commit.Verification())
		table.Append([]string{commit.Commit, t, signed, commit.Signer})
	}

	table.Render()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
