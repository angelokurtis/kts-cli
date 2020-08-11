package bash

import (
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"os/exec"
)

func RunAndLogRead(cmd string) (out []byte, err error) {
	color.Comment.Println(cmd)
	return Run(cmd)
}

func RunAndLogWrite(cmd string) (out []byte, err error) {
	color.Primary.Println(cmd)
	return Run(cmd)
}

func Run(cmd string) ([]byte, error) {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to execute command: %s", cmd)
	}
	return out, nil
}
