package files

import (
	"fmt"
	"os"
	"path/filepath"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func join(_ *cobra.Command, args []string) {
	wd := must(os.Getwd())
	paths := must(listAll(wd))
	paths = must(choose(paths))

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

func choose(paths []string) ([]string, error) {
	if len(paths) == 0 {
		return paths, nil
	}

	prompt := &survey.MultiSelect{
		Message: "Select the files to join:",
		Options: paths,
	}
	var selects []string

	if err := survey.AskOne(prompt, &selects, survey.WithPageSize(10), survey.WithKeepFilter(true)); err != nil {
		return nil, errors.WithStack(err)
	}

	return selects, nil
}

func write(paths []string) error {
	f, err := os.Create("all.txt")
	if err != nil {
		return fmt.Errorf("failed to create all.txt: %w", err)
	}
	defer f.Close()

	for _, path := range paths {
		if _, err := fmt.Fprintf(f, "# %s\n", path); err != nil {
			return fmt.Errorf("failed to write path header for %s: %w", path, err)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		if _, err := f.Write(content); err != nil {
			return fmt.Errorf("failed to write content for %s: %w", path, err)
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
