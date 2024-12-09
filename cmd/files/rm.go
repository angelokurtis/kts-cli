package files

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func remove(_ *cobra.Command, args []string) {
	ctx := context.Background()

	var dir string

	for i, arg := range args {
		switch i {
		case 0:
			dir = arg
		}
	}

	fileList, err := runGetBySize(ctx, dir)
	check(err)

	fileList = lo.Filter(fileList, func(item *file, index int) bool {
		return strings.TrimSpace(item.path) != "."
	})

	fileList, err = fileList.SelectFiles()
	check(err)

	for _, f := range fileList {
		err = runRemove(ctx, dir, f.path)
		check(err)
	}
}
