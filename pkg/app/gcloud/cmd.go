package gcloud

import (
	"os/exec"
	"strings"

	"github.com/gookit/color"
	"github.com/pkg/errors"
)

func runAndLogRead(args ...string) (out []byte, err error) {
	color.Secondary.Println("gcloud " + strings.Join(args, " "))
	return run(args...)
}

func runAndLogWrite(args ...string) (out []byte, err error) {
	color.Primary.Println("gcloud " + strings.Join(args, " "))
	return run(args...)
}

func run(args ...string) (out []byte, err error) {
	args = append(args, "--format=json")
	command := exec.Command("gcloud", args...)

	out, err = command.CombinedOutput()
	if err != nil {
		return nil, errors.WithMessage(err, string(out))
	}

	return out, nil
}
