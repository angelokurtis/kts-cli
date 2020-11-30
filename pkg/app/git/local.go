package git

import (
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const base = "~/wrkspc/"

type LocalRepo struct {
	Repo *url.URL
}

func NewLocalDir(repo string) (*LocalRepo, error) {
	u, err := url.Parse(repo)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &LocalRepo{Repo: u}, nil
}

func (l *LocalRepo) IsGithub() bool {
	return strings.Contains(l.Repo.Host, "github.com")
}

func (l *LocalRepo) SSHAddress() string {
	h := strings.ReplaceAll(l.Repo.Host, "www.", "")
	p := l.Repo.Path[1:]
	return fmt.Sprintf("git@%s:%s.git", h, p)
}

func (l *LocalRepo) Exist() bool {
	path := l.Path()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (l *LocalRepo) CreateIfNotExist() error {
	path := l.Path()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (l *LocalRepo) Path() string {
	path := base
	usr, _ := user.Current()
	path = filepath.Join(usr.HomeDir, path[2:])

	pn := l.Repo.Path
	if strings.HasSuffix(pn, suffix) {
		pn = pn[:len(pn)-len(suffix)]
	}
	path = filepath.Join(path, l.Repo.Host, pn)
	return path
}
