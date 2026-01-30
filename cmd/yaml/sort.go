package yaml

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/gookit/color"
)

func sort(cmd *cobra.Command, args []string) {
	filename := args[0]

	// Run: go run github.com/mikefarah/yq/v4/cmd@latest -i -P 'sort_keys(..)' file.yml
	c := exec.CommandContext(
		cmd.Context(),
		"go", "run", "github.com/mikefarah/yq/v4/cmd@latest",
		"-i", "-P",
		"'sort_keys(..)'",
		filename,
	)

	var stderr bytes.Buffer

	c.Stdout = cmd.OutOrStdout()
	c.Stderr = &stderr

	color.Primary.Println(c.String())

	if err := c.Run(); err != nil {
		cmd.PrintErrf("yq sort failed: %v\n%s", err, stderr.String())
		os.Exit(1)
	}
}
