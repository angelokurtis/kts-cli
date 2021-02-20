package yq

import (
	"bufio"
	"bytes"
	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func DeleteNode(yamlPath string, pathExpression string) error {
	_, err := bash.Run("yq delete -i " + yamlPath + " " + pathExpression)
	return err
}

func UpdateNode(yamlPath string, pathExpression string, value string) error {
	_, err := bash.Run("yq write -i " + yamlPath + " " + pathExpression + " " + value)
	return err
}

func UpdateNodeWithQuotes(yamlPath string, pathExpression string, value string) error {
	_, err := bash.Run("yq write -i " + yamlPath + " " + pathExpression + " --style=double " + value)
	return err
}

func ReadNodeValues(yamlPath string, pathExpression string) ([]string, error) {
	out, err := bash.Run("yq r " + yamlPath + " " + pathExpression)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	res := make([]string, 0, 0)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	return res, nil
}

func ReadNodeValue(yamlPath string, pathExpression string) (string, error) {
	out, err := bash.Run("yq r " + yamlPath + " " + pathExpression)
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	if err := scanner.Err(); err != nil {
		return "", err
	}
	res := make([]string, 0, 0)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	if len(res) == 0 {
		return "", nil
	}
	return res[0], nil
}

func DeleteKubernetesNodes(manifestPath string) error {
	if err := DeleteNode(manifestPath, "metadata.generation"); err != nil {
		return err
	}
	if err := DeleteNode(manifestPath, "metadata.selfLink"); err != nil {
		return err
	}
	if err := DeleteNode(manifestPath, "metadata.annotations[kubectl.kubernetes.io/last-applied-configuration]"); err != nil {
		return err
	}
	if err := DeleteNode(manifestPath, "metadata.creationTimestamp"); err != nil {
		return err
	}
	if err := DeleteNode(manifestPath, "metadata.resourceVersion"); err != nil {
		return err
	}
	if err := DeleteNode(manifestPath, "metadata.uid"); err != nil {
		return err
	}
	if err := DeleteNode(manifestPath, "status"); err != nil {
		return err
	}
	if err := DeleteNode(manifestPath, "metadata.annotations[cloud.google.com/neg]"); err != nil {
		return err
	}
	kind, err := ReadNodeValue(manifestPath, "kind")
	if err != nil {
		return err
	}
	if kind == "Service" {
		if err := DeleteNode(manifestPath, "spec.clusterIP"); err != nil {
			return err
		}
		sessionAffinity, err := ReadNodeValue(manifestPath, "spec.sessionAffinity")
		if err != nil {
			return err
		}
		if sessionAffinity == "None" {
			if err := DeleteNode(manifestPath, "spec.sessionAffinity"); err != nil {
				return err
			}
		}
	}
	if kind == "Deployment" {
		if err := DeleteNode(manifestPath, "metadata.annotations[deployment.kubernetes.io/revision]"); err != nil {
			return err
		}
		if err := DeleteNode(manifestPath, "spec.template.metadata.annotations[kubectl.kubernetes.io/restartedAt]"); err != nil {
			return err
		}
		if err := DeleteNode(manifestPath, "spec.template.metadata.creationTimestamp"); err != nil {
			return err
		}
	}
	annotations, err := ReadNodeValue(manifestPath, "metadata.annotations")
	if err != nil {
		return err
	}
	if annotations == "{}" {
		if err := DeleteNode(manifestPath, "metadata.annotations"); err != nil {
			return err
		}
	}
	return nil
}