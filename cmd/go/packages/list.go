package packages

import (
	"fmt"
	"log"

	"github.com/gookit/color"

	"github.com/disiqueira/gotree"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/golang"
)

func packages(_ *cobra.Command, args []string) {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	dirs, err := golang.ListDirectories(dir)
	check(err)

	root := gotree.New(color.BgGray.Text(dir))
	for _, ydir := range dirs {
		pkg, err := golang.DescribePackage(ydir)
		check(err)

		imports := pkg.InternalImports()
		if len(imports) > 0 {
			current := func() gotree.Tree {
				if ydir != dir {
					return root.Add(color.BgGray.Text(ydir))
				}
				return root
			}()
			for _, imp := range imports {
				current.Add(imp)
			}
		}
	}

	fmt.Println(root.Print())
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type NodeData struct {
	Key   string `json:"key"`
	Color string `json:"color"`
}

type LinkData struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func dedupe(a []string, b ...string) []string {
	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}
	for letter := range check {
		res = append(res, letter)
	}
	return res
}
