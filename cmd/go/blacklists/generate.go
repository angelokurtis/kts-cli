package blacklists

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/golang"
)

// go blacklists generate
func generate(cmd *cobra.Command, args []string) {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	packages, err := golang.ListPackages(dir)
	check(err)

	b := new(Blacklists)
	for _, pkg := range packages {
		b.Allow(pkg.Imports, pkg.Dir)
	}
	bytes, err := b.Marshal()
	check(err)

	fmt.Printf("%s\n", bytes)
}
