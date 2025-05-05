package gpg

import (
	"bufio"
	"bytes"
	"regexp"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

var (
	secRegex    = regexp.MustCompile(`^sec\s+(\S+)/(\S+)\s+(\d{4}-\d{2}-\d{2})(?:\s+\[SC])?(?:\s+\[expires:\s+(\d{4}-\d{2}-\d{2})])?`)
	uidRegex    = regexp.MustCompile(`^uid\s+\[(.*?)]\s+(.+?)\s*<(.+?)>`)
	subkeyRegex = regexp.MustCompile(`^ssb\s+(\S+)/(\S+)\s+(\d{4}-\d{2}-\d{2})(?:\s+\[(\w*)])?`)
)

type SecretKeys []*SecretKey

type SecretKey struct {
	KeyType string    `json:"keyType"`
	KeyID   string    `json:"keyID"`
	Created string    `json:"created"`
	Expires *string   `json:"expires,omitempty"`
	Uids    []*Uid    `json:"uids"`
	Subkeys []*Subkey `json:"subkeys"`
}

type Subkey struct {
	KeyType string   `json:"keyType"`
	KeyID   string   `json:"keyID"`
	Created string   `json:"created"`
	Usage   []string `json:"usage"`
}

type Uid struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Trust string `json:"trust"`
}

func NewSecretKeys(out []byte) (SecretKeys, error) {
	var keys SecretKeys
	var current *SecretKey

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()

		if matches := secRegex.FindStringSubmatch(line); matches != nil {
			current = &SecretKey{
				KeyType: matches[1],
				KeyID:   matches[2],
				Created: matches[3],
			}

			if matches[4] != "" {
				expires := matches[4]
				current.Expires = &expires
			}

			keys = append(keys, current)

			continue
		}

		if matches := uidRegex.FindStringSubmatch(line); matches != nil && current != nil {
			trust := matches[1]
			name := matches[2]
			email := matches[3]
			current.Uids = append(current.Uids, &Uid{
				Name:  name,
				Email: email,
				Trust: trust,
			})

			continue
		}

		if matches := subkeyRegex.FindStringSubmatch(line); matches != nil && current != nil {
			var usage []string
			if matches[4] != "" {
				usage = append(usage, matches[4])
			}

			current.Subkeys = append(current.Subkeys, &Subkey{
				KeyType: matches[1],
				KeyID:   matches[2],
				Created: matches[3],
				Usage:   usage,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	return keys, nil
}

func SelectSecretKey() (*SecretKey, error) {
	keys, err := ListSecretKeys()
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, nil
	}

	m := lo.KeyBy(keys, func(key *SecretKey) string { return key.KeyID })
	var k string
	prompt := &survey.Select{
		Message: "Select the GnuPG key:",
		Options: lo.Keys(m),
	}

	err = survey.AskOne(prompt, &k, survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return m[k], nil
}

func ListSecretKeys() (SecretKeys, error) {
	out, err := runAndLogRead("--list-secret-keys", "--keyid-format", "LONG")
	if err != nil {
		return nil, err
	}

	return NewSecretKeys(out)
}
