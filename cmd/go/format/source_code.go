package format

import (
	"fmt"
	"path/filepath"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/samber/lo"

	"github.com/angelokurtis/kts-cli/pkg/app/golang"
)

type SourceCodes []*SourceCode

func NewSourceCodes(currentDir string, packages golang.Packages) (SourceCodes, error) {
	sources := make(SourceCodes, 0)

	for _, pkg := range packages {
		rel, err := filepath.Rel(currentDir, pkg.Dir)
		if err != nil {
			return nil, fmt.Errorf(": %w", err)
		}

		for _, file := range pkg.GoFiles {
			sources = append(sources, &SourceCode{
				FileName:     file,
				FileDir:      pkg.Dir,
				RelativePath: rel,
			})
		}

		for _, file := range pkg.TestGoFiles {
			sources = append(sources, &SourceCode{
				FileName:     file,
				FileDir:      pkg.Dir,
				RelativePath: rel,
			})
		}
	}

	return sources, nil
}

type SourceCode struct {
	FileName     string
	FileDir      string
	RelativePath string
}

func (c *SourceCode) FullFilePath() string {
	return filepath.Join(c.FileDir, c.FileName)
}

func (c *SourceCode) FullDirPath() string {
	return c.FileDir
}

func (c *SourceCode) RelativeFilePath() string {
	return filepath.Join(c.RelativePath, c.FileName)
}

func (c *SourceCode) RelativeDirPath() string {
	return c.RelativePath
}

func (c SourceCodes) SelectMany() (SourceCodes, error) {
	result := make(SourceCodes, 0)

	if len(c) == 0 {
		return result, nil
	}

	prompt := &survey.MultiSelect{
		Message:  "Choose the files you want to format:",
		Options:  c.RelativeFilePaths(),
		PageSize: 10,
	}

	var selects []string

	if err := survey.AskOne(prompt, &selects, survey.WithPageSize(20), survey.WithKeepFilter(true)); err != nil {
		return nil, fmt.Errorf(": %w", err)
	}

	for _, s := range selects {
		sourceCode, err := c.GetByRelativeFilePath(s)
		if err != nil {
			return nil, fmt.Errorf(": %w", err)
		}

		result = append(result, sourceCode)
	}

	return result, nil
}

func (c SourceCodes) RelativeFilePaths() []string {
	result := make([]string, 0)
	for _, code := range c {
		result = append(result, code.RelativeFilePath())
	}

	return lo.Uniq(result)
}

func (c SourceCodes) GetByRelativeFilePath(relativeFilePath string) (*SourceCode, error) {
	for _, code := range c {
		if code.RelativeFilePath() == relativeFilePath {
			return code, nil
		}
	}

	return nil, fmt.Errorf(": ")
}

func (c SourceCodes) RelativeDirPaths() []string {
	result := make([]string, 0)
	for _, code := range c {
		result = append(result, code.RelativeDirPath())
	}

	return lo.Uniq(result)
}

func (c SourceCodes) FilterByRelativeDirPaths(relativeDirPaths []string) SourceCodes {
	result := make(SourceCodes, 0)

	for _, code := range c {
		for _, relativeDirPath := range relativeDirPaths {
			if code.RelativeDirPath() == relativeDirPath {
				result = append(result, code)
			}
		}
	}

	return result
}

func (c SourceCodes) FilterByRelativeFilePaths(relativeFilePaths []string) SourceCodes {
	result := make(SourceCodes, 0)

	for _, code := range c {
		for _, relativeFilePath := range relativeFilePaths {
			if code.RelativeFilePath() == relativeFilePath {
				result = append(result, code)
			}
		}
	}

	return result
}

func (c SourceCodes) Contains(srcCode *SourceCode) bool {
	for _, code := range c {
		if code.RelativeFilePath() == srcCode.RelativeFilePath() {
			return true
		}
	}

	return false
}
