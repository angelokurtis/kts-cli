package mod

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func installGoModUpgrade() error {
	shellScript := `
	#!/bin/bash
	
	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	# Check if go-mod-upgrade is installed
	if ! command -v go-mod-upgrade &> /dev/null
	then
		echo -e "${BLUE}go install github.com/oligot/go-mod-upgrade@latest${NC}"
		go install github.com/oligot/go-mod-upgrade@latest
	fi
	`

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var stderr bytes.Buffer

	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return errors.Errorf("failed to install go-mod-upgrade: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}
