package kubectl

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	yamlv3 "gopkg.in/yaml.v3"

	"github.com/angelokurtis/kts-cli/pkg/app/yq"
	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func ListResourcesOwners(resources, namespace string, allNamespaces bool) ([]Item, error) {
	cmd := "kubectl get " + resources + " -o=json"
	if allNamespaces {
		cmd += " --all-namespaces"
	} else if namespace != "" {
		cmd = cmd + " -n " + namespace
	}

	out, err := bash.RunAndLogRead(cmd)
	if err != nil {
		return nil, err
	}

	var col *collection
	if err := json.Unmarshal(out, &col); err != nil {
		return nil, errors.WithStack(err)
	}

	counter := make(map[string]int)
	for _, item := range col.Items {
		counter[item.Metadata.UID] = 0
	}

	for _, item := range col.Items {
		for _, owner := range item.Metadata.OwnerReferences {
			counter[owner.UID]++
		}
	}

	items := make([]Item, 0)

	for _, item := range col.Items {
		item.Dependents = counter[item.Metadata.UID]
		if len(item.Metadata.OwnerReferences) == 0 && item.Dependents > 0 {
			items = append(items, item)
		}
	}

	return items, nil
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

	res := make([]string, 0, 0)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	return res, nil
}

func SelectResources(resources, namespace string, allNamespaces bool) ([]*resource, error) {
	cmd := "kubectl get " + resources + " -o=json"
	if allNamespaces {
		cmd = cmd + " --all-namespaces"
	} else if namespace != "" {
		cmd = cmd + " -n " + namespace
	}

	out, err := bash.RunAndLogRead(cmd)
	if err != nil {
		return nil, err
	}

	var col *collection
	if err := json.Unmarshal(out, &col); err != nil {
		return nil, errors.WithStack(err)
	}

	links := make(map[string]*resource, 0)

	var options []string

	for _, item := range col.Items {
		split := strings.Split(item.APIVersion, "/")

		var fullKindName, group string

		if len(split) <= 1 {
			group = ""
			fullKindName = item.Kind
		} else {
			group = split[0]
			fullKindName = item.Kind + "." + group
		}

		r := &resource{
			Name:         item.Metadata.Name,
			FullKindName: fullKindName,
			Kind:         item.Kind,
			Group:        group,
			Namespace:    item.Metadata.Namespace,
		}
		key := ""

		if allNamespaces {
			key = key + r.Namespace + "/"
		}

		key = key + r.FullKindName + "/" + r.Name
		links[key] = r

		options = append(options, key)
	}

	var selects []string

	prompt := &survey.MultiSelect{
		Message: "Select the resource:",
		Options: options,
	}

	err = survey.AskOne(prompt, &selects, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]*resource, 0, len(selects))
	for _, s := range selects {
		res = append(res, links[s])
	}

	return res, nil
}

func SaveResourcesManifests(resources []*resource, keepStatus, decodeSecrets bool) error {
	for _, r := range resources {
		err := saveResourceManifest(r, keepStatus, decodeSecrets)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveResourceManifest(resource *resource, keepStatus, decodeSecrets bool) error {
	cmd := "kubectl get " + resource.FullKindName + " " + resource.Name + " -o yaml"
	if resource.Namespace != "" {
		cmd = cmd + " -n " + resource.Namespace
	}

	out, err := bash.Run(cmd)
	if err != nil {
		return err
	}

	if resource.Kind == "Secret" && resource.Group == "" && decodeSecrets {
		sec := make(map[string]interface{}, 0)

		err = yamlv3.Unmarshal(out, &sec)
		if err != nil {
			return errors.WithStack(err)
		}

		strdata := make(map[string]string, 0)

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

		sec["data"] = strdata

		out, err = yamlv3.Marshal(&sec)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	yamlFile := resource.Name + ".yaml"

	yamlPath := ""
	if resource.Namespace != "" && resource.Group != "" {
		yamlPath = fmt.Sprintf("./manifests/%s/%s.%s", resource.Namespace, resource.Kind, resource.Group)
	} else if resource.Namespace != "" && resource.Group == "" {
		yamlPath = fmt.Sprintf("./manifests/%s/%s", resource.Namespace, resource.Kind)
	} else if resource.Namespace == "" && resource.Group != "" {
		yamlPath = fmt.Sprintf("./manifests/%s.%s", resource.Kind, resource.Group)
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

	if err := deleteGeneratedFields(yamlPath+"/"+yamlFile, keepStatus); err != nil {
		return err
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

	if err := yq.DeleteNode(manifestPath, "metadata.annotations[kubectl.kubernetes.io/last-applied-configuration]"); err != nil {
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

	if err := yq.DeleteNode(manifestPath, "metadata.annotations[cloud.google.com/neg]"); err != nil {
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
		if err := yq.DeleteNode(manifestPath, "metadata.annotations[deployment.kubernetes.io/revision]"); err != nil {
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
	Name         string
	FullKindName string
	Kind         string
	Group        string
	Namespace    string
}

func newResource(l string) (*resource, error) {
	splitted := strings.Split(l, "/")
	size := len(splitted)

	if size == 8 {
		return &resource{
			Name:         splitted[7],
			Kind:         splitted[6],
			Group:        splitted[2],
			FullKindName: splitted[6] + "." + splitted[2],
			Namespace:    splitted[5],
		}, nil
	} else if size == 7 {
		return &resource{
			Name:         splitted[6],
			FullKindName: splitted[5],
			Kind:         splitted[5],
			Namespace:    splitted[4],
		}, nil
	} else if size == 6 {
		return &resource{
			Name:         splitted[5],
			Kind:         splitted[4],
			Group:        splitted[2],
			FullKindName: splitted[4] + "." + splitted[2],
			Namespace:    "",
		}, nil
	} else if size == 5 {
		return &resource{
			Name:         splitted[4],
			FullKindName: splitted[3],
			Kind:         splitted[3],
			Namespace:    "",
		}, nil
	}

	return nil, errors.New("unrecognized selfLink format: " + l)
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
