package golangci_lint

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"os/exec"
	"path/filepath"
	"strings"
)

func Run() (*Result, error) {
	cmd := "go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --issues-exit-code 0 --max-issues-per-linter 0 --max-same-issues 0 --out-format json ./..."
	color.Secondary.Println(cmd)
	data, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	res, err := UnmarshalResult(data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &res, nil
}

func UnmarshalResult(data []byte) (Result, error) {
	var r Result
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Result) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Result struct {
	Issues []*Issue `json:"Issues"`
	Report *Report  `json:"Report"`
}

func (r *Result) PrettyPrint() error {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0),
	)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, issue := range r.Issues {
		if issue.Pos == nil {
			continue
		}
		pos := color.Bold.Sprintf("%s:%d:%d", issue.Pos.Filename, issue.Pos.Line, issue.Pos.Column)
		text := color.Red.Sprint(issue.Text)
		_ = renderer
		output, err := renderer.Render(fmt.Sprintf("```%s\n%s\n```", "golang", strings.Join(issue.SourceLines, "\n")))
		if err != nil {
			return errors.WithStack(err)
		}
		fmt.Printf("%s: %s (%s)%s\n", pos, text, issue.FromLinter, output)
	}
	return nil
}

func (r *Result) FilterByFiles(files []string) (*Result, error) {
	issues := make([]*Issue, 0)
	for _, issue := range r.Issues {
		if issue.Pos == nil {
			continue
		}
		file, err := filepath.Abs(issue.Pos.Filename)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if lo.Contains(files, file) {
			issues = append(issues, issue)
		}
	}
	return &Result{
		Issues: issues,
		Report: r.Report,
	}, nil
}

type Issue struct {
	FromLinter           string      `json:"FromLinter"`
	Text                 string      `json:"Text"`
	Severity             string      `json:"Severity"`
	SourceLines          []string    `json:"SourceLines"`
	Replacement          interface{} `json:"Replacement"`
	Pos                  *Pos        `json:"Pos"`
	ExpectNoLint         bool        `json:"ExpectNoLint"`
	ExpectedNoLintLinter string      `json:"ExpectedNoLintLinter"`
}

type Pos struct {
	Filename string `json:"Filename"`
	Offset   int64  `json:"Offset"`
	Line     int64  `json:"Line"`
	Column   int64  `json:"Column"`
}

type Report struct {
	Linters []*Linter `json:"Linters"`
}

type Linter struct {
	Name             string `json:"Name"`
	Enabled          *bool  `json:"Enabled,omitempty"`
	EnabledByDefault *bool  `json:"EnabledByDefault,omitempty"`
}
