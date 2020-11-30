package git

import (
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const timeLayout = "2006-01-02 15:04:05 -0700"

func GetCommitTime(commit string, dir string) (*time.Time, error) {
	out, err := bash.Run("git -C " + dir + " show -s --format=%ci " + commit)
	if err != nil {
		return nil, err
	}
	split := strings.Split(string(out), "\n")
	val := split[len(split)-2]
	val = val[len(val)-25:]

	t, err := time.Parse(timeLayout, val)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &t, nil
}
