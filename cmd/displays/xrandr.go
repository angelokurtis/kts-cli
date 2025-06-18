package displays

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/pkg/errors"
)

func listDisplays(ctx context.Context) (DisplayList, error) {
	// Define the shell script as a string
	shellScript := `
	#!/bin/bash
	set -e
	xrandr --query
	`

	color.Secondary.Println("xrandr --query")

	// Create a new command to run the script
	cmd := exec.Command("bash", "-c", shellScript)

	// Capture the output and error
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = io.MultiWriter(os.Stdout, &stderr)

	// Run the command
	if err := cmd.Run(); err != nil {
		return nil, errors.Errorf("failed to execute xrandr: %v\nstderr: %s", err, strings.TrimSpace(stderr.String()))
	}

	return parseXrandrOutput(ctx, stdout)
}

func disableDisplay(ctx context.Context, displayName string) error {
	script := fmt.Sprintf(`
	#!/bin/bash
	set -e
	xrandr --output %s --off
	`, displayName)

	color.Primary.Printf("xrandr --output %s --off", displayName)

	cmd := exec.Command("bash", "-c", script)

	var stderr bytes.Buffer

	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.Errorf("failed to disable display: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}

func enableDisplay(ctx context.Context, displayName string, width, height, xOffset int64, isPrimary bool) error {
	primaryFlag := ""
	if isPrimary {
		primaryFlag = "--primary"
	}

	script := fmt.Sprintf(`
	#!/bin/bash
	set -e
	xrandr --output %s %s --mode %dx%d --pos %dx0 --rotate normal
	`, displayName, primaryFlag, width, height, xOffset)

	color.Primary.Printf("xrandr --output %s %s --mode %dx%d --pos %dx0 --rotate normal", displayName, primaryFlag, width, height, xOffset)

	cmd := exec.Command("bash", "-c", script)

	var stderr bytes.Buffer

	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.Errorf("failed to configure display: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}

var (
	displayRegex = regexp.MustCompile(`^([^\s]+)\s+(connected|disconnected)(\s+primary)?(?:\s+(\d+)x(\d+)\+(\d+)\+(\d+))?`)
	modeRegex    = regexp.MustCompile(`^\s+(\d+)x(\d+)\s+(.+)$`)
	refreshRegex = regexp.MustCompile(`([\d.]+)([*+ ]*)`)
)

func parseXrandrOutput(ctx context.Context, stdout bytes.Buffer) (DisplayList, error) {
	var displays DisplayList
	var currentDisplay *Display

	lines := strings.Split(stdout.String(), "\n")
	for _, line := range lines {
		if m := displayRegex.FindStringSubmatch(line); m != nil {
			if currentDisplay != nil {
				displays = append(displays, *currentDisplay)
			}

			name := m[1]
			connected := m[2] == "connected"
			currentDisplay = &Display{
				Name:      name,
				Connected: connected,
			}

			if m[3] != "" {
				val := true
				currentDisplay.Primary = &val
			}

			if m[4] != "" && connected {
				width, err := strconv.ParseInt(m[4], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid resolution width: %v", err)
				}

				height, err := strconv.ParseInt(m[5], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid resolution height: %v", err)
				}

				x, err := strconv.ParseInt(m[6], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid position X: %v", err)
				}

				y, err := strconv.ParseInt(m[7], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid position Y: %v", err)
				}

				currentDisplay.Resolution = &Resolution{Width: width, Height: height}
				currentDisplay.Position = &Position{X: x, Y: y}
			}

			continue
		}

		if currentDisplay == nil || !currentDisplay.Connected {
			continue
		}

		if m := modeRegex.FindStringSubmatch(line); m != nil {
			width, err := strconv.ParseInt(m[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid mode width: %v", err)
			}

			height, err := strconv.ParseInt(m[2], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid mode height: %v", err)
			}

			var refreshRates []float64
			mode := Mode{
				Width:  width,
				Height: height,
			}

			refreshParts := refreshRegex.FindAllStringSubmatch(m[3], -1)
			for _, part := range refreshParts {
				rateStr := part[1]
				flags := part[2]

				rate, err := strconv.ParseFloat(rateStr, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid refresh rate: %v", err)
				}

				refreshRates = append(refreshRates, rate)

				if strings.Contains(flags, "*") {
					rateInt := int64(rate)
					currentDisplay.RefreshRate = &rateInt
					mode.Active = true
				}

				if strings.Contains(flags, "+") {
					mode.Preferred = true
				}
			}

			mode.RefreshRates = refreshRates
			currentDisplay.Modes = append(currentDisplay.Modes, mode)
		}
	}

	if currentDisplay != nil {
		displays = append(displays, *currentDisplay)
	}

	return displays, nil
}
