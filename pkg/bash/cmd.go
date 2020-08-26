package bash

import (
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
)

func RunAndLogRead(cmd string) (out []byte, err error) {
	color.Secondary.Println(cmd)
	return Run(cmd)
}

func RunAndLogWrite(cmd string) (out []byte, err error) {
	color.Primary.Println(cmd)
	return Run(cmd)
}

func Run(cmd string) ([]byte, error) {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		if eerr, ok := err.(*exec.ExitError); ok {
			msg := strings.TrimSpace(string(eerr.Stderr))
			return nil, errors.Wrapf(eerr, "'%s' execution error: %s", cmd, msg)
		}
		return nil, errors.Wrapf(err, "'%s' execution error", cmd)
	}
	return out, nil
}
