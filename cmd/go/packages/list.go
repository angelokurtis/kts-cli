package packages

import (
	"fmt"

	"github.com/disiqueira/gotree"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"

	"github.com/angelokurtis/kts-cli/pkg/app/golang"
)

func packages(_ *cobra.Command, args []string) {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	dirs, err := golang.ListDirectories(dir)
	if err != nil {
		slog.Error("failed to list directories", "error", err)
		return
	}

	var opts []golang.DescribePackagesOption
	if tests {
		opts = append(opts, golang.WithTests())
	}

	if len(tags) > 0 {
		opts = append(opts, golang.WithTags(tags...))
	}

	root := gotree.New(color.BgGray.Text(dir))

	for _, subdir := range dirs {
		pkgs, err := golang.DescribePackages(subdir, opts...)
		if err != nil {
			slog.Error("failed to describe package", "dir", subdir, "error", err)
			continue
		}

		for _, pkg := range pkgs {
			imports := pkg.AllImports()
			if internal {
				imports = pkg.InternalImports()
			}

			if len(imports) == 0 {
				continue
			}

			node := root
			if subdir != dir {
				node = root.Add(color.BgGray.Text(subdir))
			}

			for _, imp := range imports {
				node.Add(imp)
			}
		}
	}

	fmt.Println(root.Print())
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
