package repositories

import (
	"fmt"
	log "log/slog"
	"os"
	"sort"
	"time"

	prettytime "github.com/andanhm/go-prettytime"
	ptr "github.com/gotidy/ptr"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	brazil  *time.Location
	printer *message.Printer
)

func init() {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Error(err.Error())
		return
	}

	brazil = loc
	printer = message.NewPrinter(language.BrazilianPortuguese)
}

func list(cmd *cobra.Command, args []string) {
	dockerhub := newDockerhubClient()
	hubuser := args[0]

	repos, err := dockerhub.ListRepositories(hubuser)
	if err != nil {
		log.Error(err.Error())
		return
	}

	sort.Slice(repos, func(i, j int) bool {
		t1 := ptr.To(repos[i].LastUpdated)
		t2 := ptr.To(repos[j].LastUpdated)

		return t1.After(t2)
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"REPOSITORY", "UPDATED", "PULL COUNT", "WEBPAGE"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetColWidth(100)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	for _, repo := range repos {
		repository := fmt.Sprintf("%s/%s", repo.Namespace, repo.Name)

		var updated string

		if repo.LastUpdated != nil {
			t := *repo.LastUpdated
			updated = fmt.Sprintf("%s (%s)", t.In(brazil).Format("02/01/2006 15:04"), prettytime.Format(t))
		}

		webpage := fmt.Sprintf("https://hub.docker.com/r/%s/%s", repo.Namespace, repo.Name)
		pulls := printer.Sprintf("%d", repo.PullCount)
		table.Append([]string{repository, updated, pulls, webpage})
	}

	table.Render()
}
