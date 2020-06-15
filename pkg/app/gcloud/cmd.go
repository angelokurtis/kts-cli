package gcloud

import (
	"errors"
	"fmt"
	"github.com/angelokurtis/kts-cli/internal/color"
	"os/exec"
	"strings"
)

func runAndLog(args ...string) (out []byte, err error) {
	fmt.Printf(color.Notice, "gcloud "+strings.Join(args, " ")+"\n")
	return run(args...)
}

func run(args ...string) (out []byte, err error) {
	args = append(args, "--format=json")
	command := exec.Command("gcloud", args...)
	out, err = command.CombinedOutput()
	if err != nil {
		return nil, errors.New(err.Error() + ":\n" + string(out))
	}
	return out, nil
}
