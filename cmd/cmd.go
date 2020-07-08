package cmd

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/cmd/gcp"
	"github.com/angelokurtis/kts-cli/cmd/git"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes"
	"github.com/angelokurtis/kts-cli/cmd/sensedia"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfg string
	cmd = &cobra.Command{
		Use:   "kts",
		Short: "kts is a Kurtis CLI with many daily utility functions",
		Run:   common.Help,
	}
)

func Execute() {
	if err := cmd.Execute(); err != nil {
		common.Exit(err)
	}
}

func init() {
	cobra.OnInitialize(func() {
		if cfg != "" {
			viper.SetConfigFile(cfg)
		} else {
			home, err := homedir.Dir()
			if err != nil {
				common.Exit(err)
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

	cmd.AddCommand(gcp.Command)
	cmd.AddCommand(git.Command)
	cmd.AddCommand(kubernetes.Command)
	cmd.AddCommand(sensedia.Command)
	cmd.AddCommand(completion)
}
