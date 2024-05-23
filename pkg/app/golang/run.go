package golang

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

func Run(dir, dep string, arg ...string) error {
	cmd := exec.Command("go", append([]string{"run", "-mod=mod", dep}, arg...)...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
