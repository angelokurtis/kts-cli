package git

import (
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const base = "~/wrkspc/"

var (
	githubHostRegex   = regexp.MustCompile(`^([a-zA-Z0-9-]+\.)?github\.com$`)
	gitlabHostRegex   = regexp.MustCompile(`^([a-zA-Z0-9-]+\.)?gitlab\.com$`)
	codebergHostRegex = regexp.MustCompile(`^([a-zA-Z0-9-]+\.)?codeberg\.org$`)
)

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
	return githubHostRegex.MatchString(l.Repo.Host)
}

func (l *LocalRepo) IsGitlab() bool {
	return gitlabHostRegex.MatchString(l.Repo.Host)
}

func (l *LocalRepo) IsCodeberg() bool {
	return codebergHostRegex.MatchString(l.Repo.Host)
}

func (l *LocalRepo) IsGoogleOpenSource() bool {
	return strings.Contains(l.Repo.Host, "opensource.google")
} // TODO: transform https://cs.opensource.google/go/x/tools to git@source.developers.google.com:p/go/x/tools

func (l *LocalRepo) SSHAddress() string {
	h := strings.ReplaceAll(l.Repo.Host, "www.", "")
	p := l.Repo.Path[1:]

	switch {
	case isTotvsGithubRepo(l.Repo):
		h = "github-totvs"

	case isTotvsGitlabRepo(l.Repo):
		h = "gitlab-totvs"

	case l.IsGithub():
		h = "github-personal"

	case l.IsCodeberg():
		h = "codeberg-personal"
	}

	return fmt.Sprintf("git@%s:%s.git", h, p)
}

func isTotvsGithubRepo(repo *url.URL) bool {
	if !githubHostRegex.MatchString(repo.Host) {
		return false
	}

	path := strings.TrimPrefix(repo.Path, "/")
	prefixes := []string{
		"cloud104/",
		"totvs-cloud/",
		"tiagoangelototvs/",
		"BugExtermination-Co/",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func isTotvsGitlabRepo(repo *url.URL) bool {
	if !gitlabHostRegex.MatchString(repo.Host) {
		return false
	}

	path := strings.TrimPrefix(repo.Path, "/")
	return strings.HasPrefix(path, "ascenty/")
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
		err = os.MkdirAll(path, 0o755)
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
