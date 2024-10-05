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
		echo -e "${BLUE}go install github.com/bombsimon/wsl/v4/cmd...@latest${NC}"
		go install github.com/bombsimon/wsl/v4/cmd...@latest
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

func installUnconvert() error {
	shellScript := `
	#!/bin/bash
	
	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	# Check if unconvert is installed
	if ! command -v unconvert &> /dev/null
	then
		echo -e "${BLUE}go install github.com/mdempsky/unconvert@latest${NC}"
		go install github.com/mdempsky/unconvert@latest
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
		return errors.Errorf("failed to install unconvert: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}

func installAll() error {
	if err := installImportsReviser(); err != nil {
		return err
	}

	if err := installGofumpt(); err != nil {
		return err
	}

	if err := installWsl(); err != nil {
		return err
	}

	if err := installUnconvert(); err != nil {
		return err
	}

	return nil
}
