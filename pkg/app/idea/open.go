package idea

import (
	"fmt"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func Open(path string) error {
	scripts, err := listJetBrainsScripts()
	if err != nil {
		return err
	}

	script := chooseJetBrainsScript(scripts)
	_, err = bash.Run(fmt.Sprintf("nohup %s %s >/dev/null 2>&1 &", script, path))

	return err
}
