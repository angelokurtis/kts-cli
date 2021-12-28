package blacklists

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/golang"
)

func analyze(cmd *cobra.Command, args []string) {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	packages, err := golang.ListPackages(dir)
	check(err)

	dat, err := os.ReadFile(filepath.Join(dir, "blacklists.json"))
	check(err)

	blacklists, err := UnmarshalBlacklists(dat)
	check(err)

	found := make([]string, 0, 0)
	for _, blacklist := range blacklists.Set {
		for _, imp := range packages.Usages(blacklist.Import) {
			if !contains(blacklist.Except, imp.Dir) {
				// if !contains(blacklist.Except, imp.Dir) && !strings.HasPrefix(blacklist.Import, imp.ImportPath) {
				for _, s := range imp.ImportsOf(blacklist.Import) {
					found = dedupe(found, golang.Search(imp.Dir, s))
				}
			}
		}
	}
	sort.Strings(found)
	for _, s := range found {
		fmt.Println(s)
	}
	if len(found) > 0 {
		os.Exit(2)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if e == a {
			// if strings.HasPrefix(e, a) {
			return true
		}
	}
	return false
}

func dedupe(a []string, b ...string) []string {
	c := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		c[val] = 1
	}
	for letter := range c {
		res = append(res, letter)
	}
	return res
}
