package tags

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/Masterminds/semver"
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

// git tags list
func list(_ *cobra.Command, _ []string) {
	tags, err := git.ListTags(dir)
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(tags, func(i, j int) bool {
		a := tags[i]
		b := tags[j]
		av, aerr := semver.NewVersion(a.Name)
		bv, berr := semver.NewVersion(b.Name)

		if aerr != nil || berr != nil {
			return a.Time.After(*b.Time)
		}

		return av.GreaterThan(bv)
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"TAG", "COMMIT", "CREATED"})
	table.SetRowLine(false)
	table.SetBorder(false)
	table.SetColWidth(50)

	for _, tag := range tags {
		t := fmt.Sprintf("%s (%s)", tag.Time.In(brazil).Format("02/01/2006 15:04"), prettytime.Format(*tag.Time))
		table.Append([]string{tag.Name, tag.CommitID, t})
	}

	table.Render()
}
