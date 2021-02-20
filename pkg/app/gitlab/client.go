package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"log"
	"os"
)

var client *gitlab.Client

// GITLAB_BASE_URL
func init() {
	token := os.Getenv("CB_GITLAB_ACCESS_TOKEN")
	baseURL := os.Getenv("CB_GITLAB_BASE_URL")
	c, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	client = c
}
