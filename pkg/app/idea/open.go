package idea

import (
	"fmt"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func Open(path string) error {
	_, err := bash.Run(fmt.Sprintf("nohup idea %s >/dev/null 2>&1 &", path))
	return err
}
