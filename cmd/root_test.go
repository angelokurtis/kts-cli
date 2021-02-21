package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

var expected = []string{
	"kts completion",

	"kts git clone",
	"kts git sign",
	"kts git tags",

	"kts iptv channels",
	"kts iptv groups",
}

func TestCommandsTree(t *testing.T) {
	assert.Equal(t, expected, pathsOf(rootCmd))
}

func pathsOf(cmd *cobra.Command) []string {
	r := make([]string, 0, 0)
	if cmd.HasSubCommands() {
		for _, sub := range cmd.Commands() {
			r = append(r, pathsOf(sub)...)
		}
	} else {
		return []string{cmd.CommandPath()}
	}
	return r
}
