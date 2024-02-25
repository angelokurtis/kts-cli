package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/terraform"
)

var Command = &cobra.Command{
	Use:   "commands",
	Short: "Utility functions to create Terraform commands",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "apply", Run: apply})
	Command.AddCommand(&cobra.Command{Use: "destroy", Run: destroy})
}

func apply(cmd *cobra.Command, args []string) {
	resources, err := terraform.ListResources()
	if err != nil {
		log.Fatal(err)
	}

	resources, err = resources.SelectMany()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resources.ApplyCommand())
}

func destroy(cmd *cobra.Command, args []string) {
	resources, err := terraform.ListResources()
	if err != nil {
		log.Fatal(err)
	}

	resources, err = resources.SelectMany()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resources.DestroyCommand())
}
