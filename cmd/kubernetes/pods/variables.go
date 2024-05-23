package pods

import (
	"fmt"
	log "log/slog"
	"sort"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

// kube pods variables
func variables(cmd *cobra.Command, args []string) {
	pods, err := kubectl.ListPods(namespace, allNamespaces, selector)
	if err != nil {
		log.Error(err.Error())
		return
	}

	pod, err := pods.SelectOne()
	if err != nil {
		log.Error(err.Error())
		return
	}

	vars := make(map[string]string)

	for _, env := range pod.EnvironmentVariables() {
		if env.Value != "" {
			vars[env.Name] = env.Value
		} else if env.ValueFrom != nil && env.ValueFrom.ConfigMapKeyRef != nil {
			value, err := kubectl.GetConfigMapKeyValue(env.ValueFrom.ConfigMapKeyRef, pod.Metadata.Namespace)
			if err != nil {
				log.Error(err.Error())
				return
			}

			vars[env.Name] = value
		} else if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil {
			value, err := kubectl.GetSecretKeyValue(env.ValueFrom.SecretKeyRef, pod.Metadata.Namespace)
			if err != nil {
				log.Error(err.Error())
				return
			}

			vars[env.Name] = value
		}
	}

	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%s=%s\n", k, vars[k])
	}
}
