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
			def := strings.TrimSpace(string(out))
			if msg == "" && def != "" {
				return nil, errors.New(def)
			} else {
				return nil, errors.New(msg)
			}
		}
		return nil, errors.Wrapf(err, "'%s' execution error", cmd)
	}
	return out, nil
}
