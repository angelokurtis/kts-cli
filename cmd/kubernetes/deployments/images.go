package deployments

import (
	"fmt"
	log "log/slog"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

// kube deployments update-images
func updateImages(cmd *cobra.Command, args []string) {
	deploys, err := kubectl.ListDeployments(namespace, allNamespaces)
	if err != nil {
		log.Error(err.Error())
		return
	}

	deploys, err = deploys.SelectMany()
	if err != nil {
		log.Error(err.Error())
		return
	}

	containers, err := deploys.SelectContainers()
	if err != nil {
		log.Error(err.Error())
		return
	}

	tag := ""
	prompt := &survey.Input{Message: "Inform the new tag:"}

	err = survey.AskOne(prompt, &tag, survey.WithKeepFilter(true))
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, deploy := range deploys.Items {
		for _, container := range deploy.Spec.Template.Spec.Containers {
			if containers.Contains(container.Name) {
				updateImage(deploy, container, tag)
			}
		}
	}
}

func updateImage(deploy *kubectl.Deployment, container *kubectl.Container, tag string) {
	c := deploy.GetContainer(container.Name)
	if c != nil {
		name := deploy.Metadata.Name
		ns := deploy.Metadata.Namespace
		image := strings.Split(c.Image, ":")[0]
		fmt.Printf("kubectl set image deployment/%s %s=%s:%s -n %s\n", name, c.Name, image, tag, ns)
	}
}
