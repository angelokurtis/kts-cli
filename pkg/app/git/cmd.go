package git

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/internal/color"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
)

func runAndLog(args ...string) (out []byte, err error) {
	fmt.Printf(color.Notice, "git "+strings.Join(args, " ")+"\n")
	return run(args...)
	//return nil, nil
}

func run(args ...string) (out []byte, err error) {
	cmd := exec.Command("git", args...)
	//for _, v := range cmd.Args {
	//	fmt.Println(v)
	//}
	out, err = cmd.CombinedOutput()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return out, nil
}
