package golang

import (
	"bufio"
	"context"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"

	intlformat "github.com/angelokurtis/kts-cli/cmd/go/format"
	"github.com/angelokurtis/kts-cli/pkg/app/golang"
)

func clean(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}

	pkgs, err := golang.ListPackages(wd)
	if err != nil {
		return fmt.Errorf("listing packages: %w", err)
	}

	source, err := intlformat.NewSourceCodes(wd, pkgs)
	if err != nil {
		return fmt.Errorf("loading source code: %w", err)
	}

	files, err := source.SelectMany()
	if err != nil {
		return err
	}

	for _, path := range files.RelativeFilePaths() {
		if err := removeComments(path); err != nil {
			return err
		}

		slog.InfoContext(ctx, "Removed comments", slog.String("file", path))
	}

	return nil
}

func removeComments(path string) error {
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return err
	}

	file.Comments = nil

	var buf strings.Builder
	if err := format.Node(&buf, fset, file); err != nil {
		return err
	}

	cleaned := removeEmptyLines(buf.String())

	return os.WriteFile(path, []byte(cleaned), 0o644)
}

func removeEmptyLines(s string) string {
	var result strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			result.WriteString(line + "\n")
		}
	}

	return result.String()
}
