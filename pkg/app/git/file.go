package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Files []*File

func (f Files) SelectFiles() (Files, error) {
	staged := f.StagedFiles()

	defaults, err := staged.RelativePaths()
	if err != nil {
		return nil, err
	}

	files, err := f.RelativePaths()
	if err != nil {
		return nil, err
	}

	var selects []string

	prompt := &survey.MultiSelect{
		Message:  "Choose the files you want to commit:",
		Options:  files,
		Default:  defaults,
		PageSize: 10,
	}
	if err = survey.AskOne(prompt, &selects, survey.WithPageSize(20), survey.WithKeepFilter(true)); err != nil {
		return nil, errors.WithStack(err)
	}

	return f.FilterByRelativePaths(selects)
}

func (f Files) StagedFiles() Files {
	return lo.Filter(f, func(file *File, _ int) bool {
		return isStaged(file.Status)
	})
}

func (f Files) UnStagedFiles() Files {
	return lo.Filter(f, func(file *File, _ int) bool {
		return isUnstaged(file.Status)
	})
}

func (f Files) RelativePaths() ([]string, error) {
	files := make([]string, 0)

	for _, file := range f {
		relpath, err := file.RelativePath()
		if err != nil {
			return nil, err
		}

		files = append(files, relpath)
	}

	return files, nil
}

func (f Files) FilterByRelativePaths(relpaths []string) (Files, error) {
	result := make(Files, 0)

	for _, file := range f {
		relpath, err := file.RelativePath()
		if err != nil {
			return nil, err
		}

		if lo.Contains(relpaths, relpath) {
			result = append(result, file)
		}
	}

	return result, nil
}

type File struct {
	Path   string
	Status string
}

func NewFileFromShortStatus(text string) (*File, error) {
	if len(text) <= 2 {
		return nil, fmt.Errorf("unexpected file status: %q", text)
	}

	// Extract status and filename from the text.
	status := text[:2]
	filename := strings.TrimSpace(text[2:]) // Trim leading and trailing spaces from filename.

	// Split the filename into parts, and filter out any empty strings.
	parts := strings.Fields(filename) // Fields automatically trims and splits on whitespace.

	// If there are no valid parts, return an error.
	if len(parts) == 0 {
		return nil, fmt.Errorf("unexpected file status: %q", text)
	}

	// The last part is considered the relative path.
	relpath := parts[len(parts)-1]

	// Convert relative path to absolute path.
	abspath, err := filepath.Abs(relpath)
	if err != nil {
		return nil, err
	}

	return &File{Path: abspath, Status: status}, nil
}

func (f *File) RelativePath() (string, error) {
	current, err := os.Getwd()
	if err != nil {
		return "", errors.WithStack(err)
	}

	relpath, err := filepath.Rel(current, f.Path)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return relpath, nil
}
