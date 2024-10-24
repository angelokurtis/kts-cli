package git

import (
	"fmt"
	"io"
	log "log/slog"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/git"
)

func funnyCommit(cmd *cobra.Command, args []string) {
	files, err := git.ListStagedFiles()
	if err != nil {
		log.Error(err.Error())
		return
	}

	if len(files) == 0 {
		log.Info("Nothing to commit")
		return
	}

	message, err := whatTheCommit()
	if err != nil {
		log.Error(err.Error())
		return
	}

	if err = git.DoCommitStagedFiles(message); err != nil {
		log.Error(err.Error())
		return
	}
}

func whatTheCommit() (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://whatthecommit.com/index.txt", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return strings.TrimSpace(string(bodyText)), nil
}
