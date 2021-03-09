package cmd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestDockerImagesCmd(t *testing.T) {
	os.Args = strings.Split("kts docker images golangci/golangci-lint", " ")
	err := dockerImagesCmd.Execute()
	assert.NoError(t, err)
}
