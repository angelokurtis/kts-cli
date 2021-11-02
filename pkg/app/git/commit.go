package git

import (
	"encoding/json"
	"os/exec"
	"strings"
	"time"

	"github.com/enescakir/emoji"
	"github.com/pkg/errors"

	"github.com/angelokurtis/kts-cli/pkg/bash"
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

func GetCommitMessage(commit string, dir string) (string, error) {
	out, err := bash.Run("git -C " + dir + " show -s --format=%B " + commit)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func GetCommitVerificationMessage(commit string, dir string) (string, error) {
	out, err := bash.Run("git -C " + dir + " show -s --format=%GG " + commit)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func ListCommits(dir string) (Commits, error) {
	cmd := "git -C " + dir + ` log --pretty=format:'{"commit": "%H","time": "%cI","message": "%f","verification_flag":"%G?","signer":"%GS","signer_key":"%GK","author":{"name":"%aN","email":"%aE","date":"%aD"},"commiter":{"name":"%cN","email":"%cE","date":"%cD"}}' | jq -s .`
	j, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	commits, err := UnmarshalCommits(j)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return commits, nil
}

type Commits []Commit

func UnmarshalCommits(data []byte) (Commits, error) {
	var r Commits
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Commits) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Commit struct {
	Commit           string    `json:"commit"`
	Time             time.Time `json:"time"`
	Message          string    `json:"message"`
	VerificationFlag string    `json:"verification_flag"`
	Signer           string    `json:"signer"`
	SignerKey        string    `json:"signer_key"`
	Author           Author    `json:"author"`
	Commiter         Author    `json:"commiter"`
}

// show "G" for a good (valid) signature, "B" for a bad signature, "U" for a good signature with unknown validity, "X"
// for a good signature that has expired, "Y" for a good signature made by an expired key, "R" for a good signature
// made by a revoked key, "E" if the signature cannot be checked (e.g. missing key) and "N" for no signature
func (c *Commit) VerificationStatus() emoji.Emoji {
	switch c.VerificationFlag {
	case "G", "U", "X":
		return emoji.CheckMarkButton
	default:
		return emoji.CrossMark
	}
}

func (c *Commit) Verification() string {
	switch c.VerificationFlag {
	case "G":
		return "ok"
	case "B":
		return "bad"
	case "U":
		return "unknown"
	case "X", "Y":
		return "expired"
	case "R":
		return "revoked"
	case "E":
		return "missing key"
	default:
		return "no signature"
	}
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}
