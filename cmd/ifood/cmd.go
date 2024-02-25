package ifood

import (
	"fmt"
	"os"
	"sort"
	"time"

	prettytime "github.com/andanhm/go-prettytime"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/ifood"
)

const dateTimeFormat = "02/01/2006 15:04"

var (
	brazil  *time.Location
	from    = ""
	to      = ""
	Command = &cobra.Command{
		Use:   "ifood",
		Short: "Utilities to deal with ifood's orders",
		Run:   system.Help,
	}
)

func init() {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatal(err)
	}

	brazil = loc
	statusCmd := &cobra.Command{Use: "list", Run: list}
	statusCmd.PersistentFlags().StringVar(&from, "from", "", "")
	statusCmd.PersistentFlags().StringVar(&to, "to", "", "")
	Command.AddCommand(statusCmd)
}

func list(cmd *cobra.Command, args []string) {
	orders, err := ifood.List()
	if err != nil {
		log.Fatal(err)
	}

	orders = orders.FilterByStatus("CONCLUDED")

	if from != "" {
		f, err := time.Parse(dateTimeFormat, from+" 00:00")
		if err != nil {
			log.Fatal(err)
		}

		orders = orders.FilterFrom(f)
	}

	if to != "" {
		t, err := time.Parse(dateTimeFormat, to+" 23:59")
		if err != nil {
			log.Fatal(err)
		}

		orders = orders.FilterTo(t)
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].CreatedAt.After(orders[j].CreatedAt)
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetColWidth(50)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"DATE", "ORDER", "ESTABLISHMENT", "VALUE", "DETAILS"})

	total := 0.0

	for _, o := range orders {
		s := o.ShortID
		v := float64(o.Bag.Total.ValueWithDiscount) / 100.0
		t := fmt.Sprintf("%s (%s)", o.CreatedAt.In(brazil).Format(dateTimeFormat), prettytime.Format(o.CreatedAt))
		table.Append([]string{t, "#" + s[len(s)-4:], o.Merchant.Name, fmt.Sprintf("%.2f", v), "https://www.ifood.com.br/pedido/" + o.ID})
		total = total + v
	}

	table.SetFooter([]string{"", "", "TOTAL", fmt.Sprintf("%.2f", total), ""})
	table.Render()
}
