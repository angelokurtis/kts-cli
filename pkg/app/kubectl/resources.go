package kubectl

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	yamlv3 "gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"

	"github.com/angelokurtis/kts-cli/pkg/app/yq"
	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func ListResourcesOwners(resources, namespace string, includeNoOwners, allNamespaces bool) ([]*ResourceOwner, error) {
	list, err := ListUnstructureds(resources, namespace, allNamespaces)
	if err != nil {
		return nil, err
	}

	counter := make(map[types.UID]int)
	for _, item := range list.Items {
		counter[item.GetUID()] = 0
	}

	for _, item := range list.Items {
		for _, owner := range item.GetOwnerReferences() {
			counter[owner.UID]++
		}
	}

	mapper, err := newMapper()
	if err != nil {
		return nil, err
	}

	items := make([]*ResourceOwner, 0)

	for _, item := range list.Items {
		dependents := counter[item.GetUID()]
		if len(item.GetOwnerReferences()) == 0 && (includeNoOwners || dependents > 0) {
			items = append(items, &ResourceOwner{
				Name:       getFQN(&item, mapper),
				Namespace:  item.GetNamespace(),
				Dependents: dependents,
			})
		}
	}

	return items, nil
}

type ResourceOwner struct {
	Name       string
	Namespace  string
	Dependents int
}

func ListResources(resources, namespace string, allNamespaces bool) ([]string, error) {
	cmd := []string{"get", resources}
	if allNamespaces {
		cmd = append(cmd, "--all-namespaces")
	} else if namespace != "" {
		cmd = append(cmd, "-n", namespace)
	}

	out, err := runAndLogRead(cmd...)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	if err := scanner.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]string, 0)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	return res, nil
}

func SelectResources(resources, namespace string, allNamespaces bool) ([]*resource, error) {
	list, err := ListUnstructureds(resources, namespace, allNamespaces)
	if err != nil {
		return nil, err
	}

	mapper, err := newMapper()
	if err != nil {
		return nil, err
	}

	m := lo.KeyBy(list.Items, func(item unstructured.Unstructured) string {
		return getFQN(&item, mapper)
	})

	var selects []string

	prompt := &survey.MultiSelect{
		Message: "Select the resource:",
		Options: lo.Keys(m),
	}

	err = survey.AskOne(prompt, &selects, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]*resource, 0, len(selects))
	for _, s := range selects {
		u := m[s]

		gvk := u.GroupVersionKind()

		// Use RESTMapper to find the resource
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Fatalf("Failed to get REST mapping: %w", err)
		}

		kind := mapping.Resource.Resource // This gives plural + group
		if len(gvk.Group) > 0 {
			kind = fmt.Sprintf("%s.%s", kind, gvk.Group)
		}

		res = append(res, &resource{
			Name:      u.GetName(),
			Namespace: u.GetNamespace(),
			Kind:      kind,
		})
	}

	return res, nil
}

func SaveResourcesManifests(resources []*resource, keepStatus, sanitize, decodeSecrets, groupByNamespace bool) error {
	for _, r := range resources {
		err := saveResourceManifest(r, keepStatus, sanitize, decodeSecrets, groupByNamespace)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveResourceManifest(resource *resource, keepStatus, sanitize, decodeSecrets, groupByNamespace bool) error {
	cmd := "kubectl get " + resource.Kind + "/" + resource.Name + " -o yaml"
	if resource.Namespace != "" {
		cmd = cmd + " -n " + resource.Namespace
	}

	out, err := bash.Run(cmd)
	if err != nil {
		return err
	}

	if resource.Kind == "secrets" && decodeSecrets {
		sec := make(map[string]interface{})

		err = yamlv3.Unmarshal(out, &sec)
		if err != nil {
			return errors.WithStack(err)
		}

		strdata := make(map[string]string)

		if data := sec["data"]; data != nil {
			if kv, ok := data.(map[string]interface{}); ok {
				for k, v := range kv {
					srtv, err := base64.StdEncoding.DecodeString(fmt.Sprintf("%s", v))
					if err != nil {
						return errors.WithStack(err)
					}

					strdata[k] = string(srtv)
				}
			}
		}

		sec["stringData"] = strdata
		delete(sec, "data")

		out, err = yamlv3.Marshal(&sec)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	yamlFile := resource.Name + ".yaml"

	yamlPath := ""
	if resource.Namespace != "" && groupByNamespace {
		yamlPath = fmt.Sprintf("./manifests/%s/%s", resource.Namespace, resource.Kind)
	} else {
		yamlPath = fmt.Sprintf("./manifests/%s", resource.Kind)
	}

	_, err = bash.Run("mkdir -p " + yamlPath)
	if err != nil {
		return err
	}

	color.Primary.Println(cmd + " > " + yamlPath + "/" + yamlFile)

	if err = ioutil.WriteFile(yamlPath+"/"+yamlFile, out, 0o644); err != nil {
		return errors.WithStack(err)
	}

	if sanitize {
		if err := deleteGeneratedFields(yamlPath+"/"+yamlFile, keepStatus); err != nil {
			return err
		}
	}

	return nil
}

func deleteGeneratedFields(manifestPath string, keepStatus bool) error {
	if !keepStatus {
		if err := yq.DeleteNode(manifestPath, "status"); err != nil {
			return err
		}
	}

	if err := yq.DeleteNode(manifestPath, "metadata.managedFields"); err != nil {
		return err
	}

	if err := yq.DeleteNode(manifestPath, "metadata.generation"); err != nil {
		return err
	}

	if err := yq.DeleteNode(manifestPath, "metadata.selfLink"); err != nil {
		return err
	}

	if err := yq.DeleteNode(manifestPath, `metadata.annotations["kubectl.kubernetes.io/last-applied-configuration"]`); err != nil {
		return err
	}

	if err := yq.DeleteNode(manifestPath, "metadata.creationTimestamp"); err != nil {
		return err
	}

	if err := yq.DeleteNode(manifestPath, "metadata.resourceVersion"); err != nil {
		return err
	}

	if err := yq.DeleteNode(manifestPath, "metadata.uid"); err != nil {
		return err
	}

	if err := yq.DeleteNode(manifestPath, `metadata.annotations["cloud.google.com/neg"]`); err != nil {
		return err
	}

	kind, err := yq.ReadNodeValue(manifestPath, "kind")
	if err != nil {
		return err
	}

	if kind == "Service" {
		if err := yq.DeleteNode(manifestPath, "spec.clusterIP"); err != nil {
			return err
		}

		sessionAffinity, err := yq.ReadNodeValue(manifestPath, "spec.sessionAffinity")
		if err != nil {
			return err
		}

		if sessionAffinity == "None" {
			if err := yq.DeleteNode(manifestPath, "spec.sessionAffinity"); err != nil {
				return err
			}
		}
	}

	if kind == "Deployment" {
		if err := yq.DeleteNode(manifestPath, `metadata.annotations["deployment.kubernetes.io/revision"]`); err != nil {
			return err
		}

		if err := yq.DeleteNode(manifestPath, "spec.template.metadata.creationTimestamp"); err != nil {
			return err
		}
	}

	annotations, err := yq.ReadNodeValue(manifestPath, "metadata.annotations")
	if err != nil {
		return err
	}

	if annotations == "{}" {
		if err := yq.DeleteNode(manifestPath, "metadata.annotations"); err != nil {
			return err
		}
	}

	return nil
}

type resource struct {
	Name      string
	Namespace string
	Kind      string
}

type collection struct {
	Items []Item `json:"items"`
}

type Item struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Dependents int      `json:"-"`
}
