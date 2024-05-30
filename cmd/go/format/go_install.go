package format

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func installImportsReviser() error {
	shellScript := `
	#!/bin/bash
	
	# Check if goimports-reviser is installed
	if ! command -v goimports-reviser &> /dev/null
	then
		echo "goimports-reviser is not installed. Installing..."
		go install github.com/incu6us/goimports-reviser/v3@latest
	fi
	`

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var out bytes.Buffer

	var stderr bytes.Buffer

	cmd.Stdout = &out
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
	
	# Check if gofumpt is installed
	if ! command -v gofumpt &> /dev/null
	then
		echo "gofumpt is not installed. Installing..."
		go install mvdan.cc/gofumpt@latest
	fi
	`

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var out bytes.Buffer

	var stderr bytes.Buffer

	cmd.Stdout = &out
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
	
	# Check if wsl is installed
	if ! command -v wsl &> /dev/null
	then
		echo "wsl is not installed. Installing..."
		go install mvdan.cc/wsl@latest
	fi
	`

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var out bytes.Buffer

	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return errors.Errorf("failed to install wsl: %s", strings.TrimSpace(stderr.String()))
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

	return nil
}
