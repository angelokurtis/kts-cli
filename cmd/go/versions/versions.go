package versions

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func versions(_ *cobra.Command, args []string) {
	vs, err := bash.RunAndLogRead("gvm listall")
	check(err)

	for _, v := range strings.Split(string(vs), "\n") {
		v = strings.TrimSpace(v)
		if strings.HasPrefix(v, "go") {
			fmt.Println(v[2:])
		}
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
