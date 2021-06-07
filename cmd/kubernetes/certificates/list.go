package certificates

import (
	"fmt"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/spf13/cobra"
)

func list(cmd *cobra.Command, args []string) {
	secrets, err := kubectl.ListTLSSecrets()
	check(err)

	for _, sec := range secrets.Items {
		fmt.Println(sec.Metadata.Name)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
