package gpg

import (
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
)

func runAndLogRead(args ...string) (out []byte, err error) {
	color.Secondary.Println("gpg " + strings.Join(args, " "))
	return run(args...)
}

func runAndLogWrite(args ...string) (out []byte, err error) {
	color.Primary.Println("gpg " + strings.Join(args, " "))
	return run(args...)
}

func run(args ...string) (out []byte, err error) {
	command := exec.Command("gpg", args...)
	out, err = command.CombinedOutput()
	if err != nil {
		return nil, errors.WithMessage(err, string(out))
	}
	return out, nil
}
