package cmd

import (
	"fmt"

	"github.com/angelokurtis/kts-cli/cmd/aws"
	"github.com/angelokurtis/kts-cli/cmd/gcp"
	"github.com/angelokurtis/kts-cli/cmd/git"
	golang "github.com/angelokurtis/kts-cli/cmd/go"
	"github.com/angelokurtis/kts-cli/cmd/ifood"
	"github.com/angelokurtis/kts-cli/cmd/iptv"
	"github.com/angelokurtis/kts-cli/cmd/istio"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes"
	"github.com/angelokurtis/kts-cli/cmd/terraform"
	"github.com/angelokurtis/kts-cli/cmd/terraformer"
	"github.com/angelokurtis/kts-cli/cmd/yaml"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfg string
	cmd = &cobra.Command{
		Use:   "kts",
		Short: "kts is a Kurtis CLI with many daily utility functions",
		Run:   system.Help,
	}
)

func Execute() {
	if err := cmd.Execute(); err != nil {
		system.Exit(err)
	}
}

func init() {
	cobra.OnInitialize(func() {
		if cfg != "" {
			viper.SetConfigFile(cfg)
		} else {
			home, err := homedir.Dir()
			if err != nil {
				system.Exit(err)
			}

			viper.AddConfigPath(home)
			viper.SetConfigName(".kurtis")
		}

		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	})
	cmd.PersistentFlags().StringVar(&cfg, "config", "", "config file (default is $HOME/.kurtis.yaml)")

	cmd.AddCommand(completion)
	cmd.AddCommand(aws.Command)
	cmd.AddCommand(iptv.Command)
	cmd.AddCommand(gcp.Command)
	cmd.AddCommand(ifood.Command)
	cmd.AddCommand(git.Command)
	cmd.AddCommand(kubernetes.Command)
	cmd.AddCommand(terraform.Command)
	cmd.AddCommand(terraformer.Command)
	cmd.AddCommand(yaml.Command)
	cmd.AddCommand(istio.Command)
	cmd.AddCommand(golang.Command)
}
