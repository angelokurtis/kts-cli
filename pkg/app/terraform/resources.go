package terraform

import (
	"fmt"
	"io/ioutil"

	"github.com/AlecAivazis/survey/v2"
	changeCase "github.com/ku/go-change-case"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func ListProviderResources(provider string) []string {
	switch provider {
	case "helm":
		return []string{
			"helm_release",
		}
	case "google":
		return []string{
			"google_container_cluster",
			"google_container_node_pool",
		}
	case "aws":
		return nil
	case "kubernetes":
		return []string{
			// "kubernetes_api_service",
			// "kubernetes_certificate_signing_request",
			// "kubernetes_cluster_role",
			// "kubernetes_cluster_role_binding",
			"kubernetes_config_map",
			// "kubernetes_cron_job",
			// "kubernetes_csi_driver",
			// "kubernetes_daemonset",
			// "kubernetes_default_service_account",
			"kubernetes_deployment",
			// "kubernetes_endpoints",
			// "kubernetes_horizontal_pod_autoscaler",
			// "kubernetes_ingress",
			// "kubernetes_job",
			// "kubernetes_limit_range",
			// "kubernetes_mutating_webhook_configuration",
			"kubernetes_namespace",
			// "kubernetes_network_policy",
			// "kubernetes_persistent_volume",
			// "kubernetes_persistent_volume_claim",
			// "kubernetes_pod",
			// "kubernetes_pod_disruption_budget",
			// "kubernetes_pod_security_policy",
			// "kubernetes_priority_class",
			// "kubernetes_replication_controller",
			// "kubernetes_resource_quota",
			// "kubernetes_role",
			// "kubernetes_role_binding",
			// "kubernetes_secret",
			// "kubernetes_service",
			// "kubernetes_service_account",
			// "kubernetes_stateful_set",
			// "kubernetes_storage_class",
			// "kubernetes_validating_webhook_configuration",
		}
	}
	return nil
}

func SelectResource(provider string) (*Resource, error) {
	r := ListProviderResources(provider)

	var selected string
	if len(r) == 0 {
		return nil, errors.New("no resources where found")
	} else if len(r) > 1 {
		prompt := &survey.Select{
			Message: "Select the Terraform Resource:",
			Options: r,
		}

		err := survey.AskOne(prompt, &selected, survey.WithPageSize(10), survey.WithKeepFilter(true))
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		selected = r[0]
	}

	return BuildResource(selected)
}

func YAMLDecode(filename string) ([]byte, error) {
	out, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	r := &kubeResource{}
	err = yaml.Unmarshal(out, r)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	out, err = bash.RunAndLogRead(fmt.Sprintf("echo 'yamldecode(file(\"%s\"))' | terraform console", filename))
	if err != nil {
		return nil, err
	}

	prefix := fmt.Sprintf("resource \"kubernetes_manifest\" \"%s\" {\n provider = kubernetes-alpha\n\n manifest =\n", changeCase.Snake(r.Metadata.Name+"_"+r.Kind))
	suffix := "\n}"
	return []byte(prefix + string(out) + suffix), err
}

type kubeResource struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec interface{} `yaml:"spec"`
}
