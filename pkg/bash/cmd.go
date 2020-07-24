package bash

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/internal/color"
	"github.com/pkg/errors"
	"os/exec"
)

func RunAndLog(cmd string) (out []byte, err error) {
	fmt.Printf(color.Notice, cmd+"\n")
	return Run(cmd)
}

func Run(cmd string) ([]byte, error) {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to execute command: %s", cmd)
	}
	return out, nil
}
