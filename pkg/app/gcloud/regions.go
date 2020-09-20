package gcloud

import (
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"strings"
)

func CurrentRegion() (string, error) {
	out, err := bash.RunAndLogRead("gcloud config get-value compute/zone")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
