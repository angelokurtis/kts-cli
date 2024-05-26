package golang

import (
	"os"
	"os/exec"
	"strings"

	"github.com/gookit/color"
	"github.com/pkg/errors"
)

func Run(dir, dep string, arg ...string) error {
	color.Primary.Println(strings.Join(append([]string{"go", "run", "-mod=readonly", dep}, arg...), " "))
	cmd := exec.Command("go", append([]string{"run", "-mod=readonly", dep}, arg...)...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func RunSilently(dir, dep string, arg ...string) error {
	color.Primary.Println(strings.Join(append([]string{"go", "run", "-mod=readonly", dep}, arg...), " "))
	cmd := exec.Command("go", append([]string{"run", "-mod=readonly", dep}, arg...)...)
	cmd.Dir = dir

	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
