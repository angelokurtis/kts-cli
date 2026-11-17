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

func installFuncorder() error {
	shellScript := `
	#!/bin/bash
	set -e

	# Define colors
	BLUE='\033[0;34m'
	NC='\033[0m' # No Color

	# Check if funcorder is installed
	if ! command -v funcorder &> /dev/null
	then
		echo -e "${BLUE}go install github.com/manuelarte/funcorder@latest${NC}"
		go install github.com/manuelarte/funcorder@latest
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
		return errors.Errorf("failed to install funcorder: %s", strings.TrimSpace(stderr.String()))
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

	if enableFuncorder {
		if err := installFuncorder(); err != nil {
			return err
		}
	}

	return nil
}
