package git

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func ListTags(dir string) ([]*Tag, error) {
	out, err := bash.RunAndLogRead(fmt.Sprintf("git -C %s show-ref --tags", dir))
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	tags := make([]*Tag, 0, 0)

	for scanner.Scan() {
		split := strings.Split(scanner.Text(), " ")

		tag, err := newTag(split[0], strings.ReplaceAll(split[1], "refs/tags/", ""), dir)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

type Tag struct {
	Time     *time.Time
	CommitID string
	Name     string
}

func newTag(commitID, name, dir string) (*Tag, error) {
	t, err := GetCommitTime(commitID, dir)
	if err != nil {
		return nil, err
	}

	return &Tag{Time: t, CommitID: commitID, Name: name}, nil
}
