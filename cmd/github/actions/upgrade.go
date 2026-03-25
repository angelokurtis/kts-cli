package actions

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/google/go-github/v84/github"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/spf13/cobra"
)

func upgrade(cmd *cobra.Command, args []string) error {
	client := &http.Client{Transport: &RequestLogger{
		DefaultTransport: cleanhttp.DefaultTransport(),
	}}

	gh := github.NewClient(client)

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		gh = gh.WithAuthToken(token)
	}

	return newUpdater(gh).run(cmd.Context())
}

type updater struct {
	gh *github.Client

	mu    sync.Mutex
	cache map[string]version
}

type version struct {
	tag string
	sha string
}

func newUpdater(gh *github.Client) *updater {
	return &updater{
		gh:    gh,
		cache: make(map[string]version),
	}
}

func (u *updater) run(ctx context.Context) error {
	return u.updateWorkflows(ctx, ".")
}

var usesRe = regexp.MustCompile(`uses:\s*([a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+)@([^\s#]+)`)

func (u *updater) updateWorkflows(ctx context.Context, root string) error {
	dir := filepath.Join(root, ".github", "workflows")

	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".yml") && !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var out bytes.Buffer

		scanner := bufio.NewScanner(bytes.NewReader(data))

		for scanner.Scan() {
			line := scanner.Text()
			m := usesRe.FindStringSubmatch(line)

			if len(m) == 3 {
				repo := m[1]

				owner, name, err := splitRepo(repo)
				if err != nil {
					return err
				}

				key := owner + "/" + name

				u.mu.Lock()
				cached, ok := u.cache[key]
				u.mu.Unlock()

				if !ok {
					tag, sha, err := u.latest(ctx, owner, name)
					if err != nil {
						return err
					}

					cached = version{tag: tag, sha: sha}

					u.mu.Lock()
					u.cache[key] = cached
					u.mu.Unlock()
				}

				newLine := fmt.Sprintf("uses: %s@%s # %s", repo, cached.sha, cached.tag)
				prefix := line[:strings.Index(line, "uses:")]
				line = prefix + newLine
			}

			out.WriteString(line + "\n")
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		return os.WriteFile(path, out.Bytes(), 0o644)
	})
}

func splitRepo(s string) (owner, repo string, err error) {
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return "", "", errors.New("invalid action format")
	}

	owner = strings.TrimSpace(parts[0])

	repo = strings.TrimSpace(parts[1])
	if owner == "" || repo == "" {
		return "", "", errors.New("invalid action format")
	}

	return owner, repo, nil
}

func (u *updater) latest(ctx context.Context, owner, repo string) (tag, sha string, err error) {
	rel, _, err := u.gh.Repositories.GetLatestRelease(ctx, owner, repo)
	if err == nil && rel != nil && rel.TagName != nil {
		tag = rel.GetTagName()
	} else {
		tags, _, err := u.gh.Repositories.ListTags(ctx, owner, repo, nil)
		if err != nil {
			return "", "", fmt.Errorf("list tags: %w", err)
		}

		if len(tags) == 0 {
			return "", "", errors.New("no tags found")
		}

		tag = tags[0].GetName()
	}

	ref, _, err := u.gh.Git.GetRef(ctx, owner, repo, "refs/tags/"+tag)
	if err != nil {
		return "", "", fmt.Errorf("get ref: %w", err)
	}

	obj := ref.Object
	if obj == nil {
		return "", "", errors.New("empty ref object")
	}

	switch obj.GetType() {
	case "commit":
		return tag, obj.GetSHA(), nil
	case "tag":
		tagObj, _, err := u.gh.Git.GetTag(ctx, owner, repo, obj.GetSHA())
		if err != nil {
			return "", "", fmt.Errorf("get tag: %w", err)
		}

		if tagObj.Object == nil {
			return "", "", errors.New("tag object empty")
		}

		return tag, tagObj.Object.GetSHA(), nil
	default:
		return "", "", fmt.Errorf("unsupported object type: %s", obj.GetType())
	}
}
