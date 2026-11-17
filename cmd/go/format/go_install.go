package format

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func installImportsReviser() error {
	shellScript := `
	#!/bin/bash
	
	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	# Check if goimports-reviser is installed
	if ! command -v goimports-reviser &> /dev/null
	then
		echo -e "${BLUE}go install github.com/incu6us/goimports-reviser/v3@latest${NC}"
		go install github.com/incu6us/goimports-reviser/v3@latest
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
		return errors.Errorf("failed to install goimports-reviser: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}

func installGofumpt() error {
	shellScript := `
	#!/bin/bash
	
	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	# Check if gofumpt is installed
	if ! command -v gofumpt &> /dev/null
	then
		echo -e "${BLUE}go install mvdan.cc/gofumpt@latest${NC}"
		go install mvdan.cc/gofumpt@latest
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
		return errors.Errorf("failed to install gofumpt: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}

func installWsl() error {
	shellScript := `
	#!/bin/bash
	
	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	# Check if wsl is installed
	if ! command -v wsl &> /dev/null
	then
		echo -e "${BLUE}go install github.com/bombsimon/wsl/v5/cmd/wsl@latest${NC}"
		go install github.com/bombsimon/wsl/v5/cmd/wsl@latest
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
		return errors.Errorf("failed to install wsl: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}

func installGoreorder() error {
	shellScript := `
	#!/bin/bash
	set -e

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	# Check if goreorder is installed
	if ! command -v goreorder &> /dev/null
	then
		echo -e "${BLUE}go install github.com/metal3d/goreorder/cmd/goreorder@latest${NC}"
		go install github.com/metal3d/goreorder/cmd/goreorder@latest
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
		return errors.Errorf("failed to install goreorder: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}

func installCallorder() error {
	shellScript := `
	#!/bin/bash
	set -e

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	# Check if callorder is installed
	if ! command -v callorder &> /dev/null
	then
		echo -e "${BLUE}go install github.com/tho/callorder/cmd/callorder@latest${NC}"
		go install github.com/tho/callorder/cmd/callorder@latest
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
		return errors.Errorf("failed to install callorder: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}

func installAll() error {
	if !disableImportsReviser {
		if err := installImportsReviser(); err != nil {
			return err
		}
	}

	if !disableGofumpt {
		if err := installGofumpt(); err != nil {
			return err
		}
	}

	if !disableWsl {
		if err := installWsl(); err != nil {
			return err
		}
	}

	if enableGoreorder {
		if err := installGoreorder(); err != nil {
			return err
		}
	}

	if enableCallorder {
		if err := installCallorder(); err != nil {
			return err
		}
	}

	return nil
}
