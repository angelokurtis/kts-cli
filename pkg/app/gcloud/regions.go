package gcloud

import (
	"strings"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func CurrentRegion() (string, error) {
	out, err := bash.RunAndLogRead("gcloud config get-value compute/zone")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
