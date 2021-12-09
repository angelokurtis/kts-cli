package gpg

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

func SelectSecretKey() (*SecretKey, error) {
	keys, err := ListSecretKeys()
	if err != nil {
		return nil, err
	}
	if keys == nil || len(keys.Items) == 0 {
		return nil, nil
	} else if len(keys.Items) == 1 {
		return keys.Items[0], nil
	}

	options := make([]string, 0, 0)
	m := make(map[string]*SecretKey)
	for _, key := range keys.Items {
		options = append(options, key.UID)
		m[key.UID] = key
	}

	var k string
	prompt := &survey.Select{
		Message: "Select the GnuPG key:",
		Options: options,
	}

	err = survey.AskOne(prompt, &k, survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return m[k], nil
}

func ListSecretKeys() (*Keys, error) {
	out, err := runAndLogRead("--list-secret-keys", "--keyid-format", "LONG")
	if err != nil {
		return nil, err
	}
	return NewKeys(out)
}

type Keys struct {
	Items []*SecretKey
}

func NewKeys(out []byte) (*Keys, error) {
	items := make([]*SecretKey, 0, 0)

	var current *SecretKey
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "sec ") {
			sec := strings.Fields(line)[1]
			current = &SecretKey{Sec: sec}
		}
		if strings.HasPrefix(line, "uid ") && current != nil {
			uid := strings.Fields(line)[2:]
			current.UID = strings.Join(uid, " ")
			items = append(items, &SecretKey{Sec: current.Sec, UID: current.UID})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Keys{Items: items}, nil
}

type SecretKey struct {
	Sec string
	UID string
}
