package networks

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	prettytime "github.com/andanhm/go-prettytime"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func list(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	dieOnErr(err)
	defer apiClient.Close()

	networks, err := apiClient.NetworkList(ctx, network.ListOptions{})
	dieOnErr(err)

	sort.Slice(networks, func(i, j int) bool {
		t1 := networks[i].Created
		t2 := networks[j].Created

		return t1.Before(t2)
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetBorder(false)
	table.SetHeader([]string{"NETWORK ID", "NAME", "DRIVER", "SCOPE", "SUBNET", "GATEWAY", "CREATED"})

	for _, n := range networks {
		subnet := make([]string, 0)
		gateway := make([]string, 0)

		for _, config := range n.IPAM.Config {
			subnet = append(subnet, config.Subnet)
			gateway = append(gateway, config.Gateway)
		}

		created := fmt.Sprintf("%s (%s)", n.Created.Format("02/01/2006 15:04"), prettytime.Format(n.Created))
		table.Append([]string{
			n.ID,
			n.Name,
			n.Driver,
			n.Scope,
			strings.Join(subnet, ", "),
			strings.Join(gateway, ", "),
			created,
		})
	}

	table.Render()
}
