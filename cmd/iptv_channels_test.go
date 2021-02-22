package cmd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestIPTVChannelsCmd(t *testing.T) {
	os.Args = strings.Split("kts iptv channels", " ")
	err := iptvChannelsCmd.Execute()
	assert.NoError(t, err)
}
