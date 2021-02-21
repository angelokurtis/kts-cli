package bash

import (
	"fmt"
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
			} else if msg == "" && def == "" {
				return nil, nil
			} else {
				return nil, errors.New(msg)
			}
		}
		return nil, errors.Wrapf(err, "'%s' execution error", cmd)
	}
	return out, nil
}

func Follow(command string) error {
	color.Primary.Println(command)
	cmdArgs := strings.Fields(command)

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	stdout, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		return err
	}

	oneByte := make([]byte, 1)
	for {
		_, err := stdout.Read(oneByte)
		if err != nil {
			break
		}
		fmt.Printf("%s", oneByte)
	}
	return nil
}
