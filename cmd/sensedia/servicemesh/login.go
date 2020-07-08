package servicemesh

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/internal/color"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/angelokurtis/kts-cli/pkg/app/sensedia"
	"github.com/spf13/cobra"
	"strings"
)

func login(cmd *cobra.Command, args []string) {
	configMaps, err := kubectl.SearchConfigMap("app.kubernetes.io/component=manager")
	if err != nil {
		common.Exit(err)
	}
	configMap, err := configMaps.SingleResult()
	if err != nil {
		common.Exit(err)
	}
	if configMap == nil {
		fmt.Printf(color.Debug, "the ConfigMap was not found, we'll try with other label\n")
		configMaps, err := kubectl.SearchConfigMap("app.sensedia.com/component=routes-manager")
		if err != nil {
			common.Exit(err)
		}
		configMap, err = configMaps.SingleResult()
		if err != nil {
			common.Exit(err)
		}
	}
	oauthProvider := configMap.Data["OAUTH2_URL_PROVIDER"]
	oauthProvider = strings.TrimSuffix(oauthProvider, "/") + "/"
	sensediaManager := configMap.Data["SENSEDIA_MANAGER_URL"]
	fmt.Printf(color.Debug, "trying to authenticate on "+sensediaManager+"\n")
	login, password, err := sensedia.SelectGatewayUser()
	if err != nil {
		common.Exit(err)
	}
	sensediaAuth, xsrf, err := sensedia.Login(sensediaManager, login, password)
	if err != nil {
		common.Exit(err)
	}
	fmt.Printf(color.Notice, oauthProvider+"oauth/authorize?sensedia_auth="+sensediaAuth+"&xsrf="+xsrf+"\n")
}
