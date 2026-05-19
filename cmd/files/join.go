package files

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	ignore "github.com/sabhiram/go-gitignore"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func join(_ *cobra.Command, args []string) {
	wd := must(os.Getwd())
	paths := must(listAll(wd))
	paths = must(filterIgnored(wd, paths))
	paths = must(choose(wd, paths))

	if err := write(paths); err != nil {
		panic(err)
	}
}

// listAll returns all file paths under the given directory recursively.
func listAll(root string) ([]string, error) {
	var paths []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // handle error, skip this path
		}

		if !info.IsDir() {
			paths = append(paths, path)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list all file paths under %s: %w", root, err)
	}

	return paths, nil
}

func choose(root string, paths []string) ([]string, error) {
	if len(paths) == 0 {
		return paths, nil
	}

	options := lo.Map(paths, func(p string, _ int) string {
		return lo.Must(filepath.Rel(root, p))
	})

	index := lo.SliceToMap(paths, func(p string) (string, string) {
		return lo.Must(filepath.Rel(root, p)), p
	})

	prompt := &survey.MultiSelect{
		Message: "Select the files to join:",
		Options: options,
	}

	var selects []string
	if err := survey.AskOne(prompt, &selects, survey.WithPageSize(20), survey.WithKeepFilter(true)); err != nil {
		return nil, errors.WithStack(err)
	}

	return lo.Map(selects, func(rel string, _ int) string {
		return index[rel]
	}), nil
}

func filterIgnored(root string, paths []string) ([]string, error) {
	gitignorePath := filepath.Join(root, ".gitignore")

	var ig *ignore.GitIgnore

	if _, err := os.Stat(gitignorePath); err == nil {
		compiled, err := ignore.CompileIgnoreFile(gitignorePath)
		if err != nil {
			return nil, err
		}

		ig = compiled
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	filtered := lo.Filter(paths, func(p string, _ int) bool {
		rel := filepath.ToSlash(lo.Must(filepath.Rel(root, p)))

		// Always ignore .git directory
		if rel == ".git" || strings.HasPrefix(rel, ".git/") {
			return false
		}

		if ig != nil && ig.MatchesPath(rel) {
			return false
		}

		return true
	})

	return filtered, nil
}

func write(paths []string) error {
	f, err := os.Create("all.txt")
	if err != nil {
		return fmt.Errorf("failed to create all.txt: %w", err)
	}
	defer f.Close()

	for _, path := range paths {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			relPath = path
		}

		if _, err := fmt.Fprintf(f, "# %s\n", relPath); err != nil {
			return fmt.Errorf("failed to write path header for %s: %w", path, err)
		}

		in, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", path, err)
		}

		scanner := bufio.NewScanner(in)
		lineNo := 1

		for scanner.Scan() {
			if _, err := fmt.Fprintf(f, "%d: %s\n", lineNo, scanner.Text()); err != nil {
				in.Close()
				return fmt.Errorf("failed to write numbered line for %s: %w", path, err)
			}

			lineNo++
		}

		if err := scanner.Err(); err != nil {
			in.Close()
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		if err := in.Close(); err != nil {
			return fmt.Errorf("failed to close file %s: %w", path, err)
		}

		if _, err := f.Write([]byte("\n")); err != nil {
			return fmt.Errorf("failed to write newline after %s: %w", path, err)
		}
	}

	return nil
}

func must[T any](out T, err error) T {
	if err != nil {
		panic(err)
	}

	return out
}
