package format

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/golang"
)

var Command = &cobra.Command{
	Use: "format",
	Run: wrapWithErrorHandler(runFormat),
}

func runFormat(cmd *cobra.Command, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	}

	if path == "" {
		return errors.New("no path provided")
	}

	pkgDetails, err := golang.DescribePackage(path)
	if err != nil {
		return errors.Errorf("failed to describe package: %v", err)
	}

	_ = pkgDetails
	// TODO:
	return nil
}
