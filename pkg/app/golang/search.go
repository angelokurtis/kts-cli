package golang

import (
	"fmt"
	"os/exec"
	"strings"
)

func Search(were, what string) string {
	cmd := fmt.Sprintf(`grep --include \*.go -Hn \"%s\" %s/*`, what, were)
	j, _ := exec.Command("bash", "-c", cmd).Output()
	return strings.TrimSpace(string(j))
}
