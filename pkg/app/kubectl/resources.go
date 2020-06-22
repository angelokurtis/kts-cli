package kubectl

import (
	"bufio"
	"bytes"
	"strings"
)

var nativeResources = []string{
	"bindings",
	"configmaps",
	"endpoints",
	"events",
	"limitranges",
	"persistentvolumeclaims",
	"pods",
	"podtemplates",
	"replicationcontrollers",
	"resourcequotas",
	"secrets",
	"serviceaccounts",
	"services",
}

func ListResources() ([]string, error) {
	resources, err := ListResourceDefinitions()
	if err != nil {
		return nil, err
	}
	out, err := runAndLog("get", strings.Join(resources, ","), "-o=name")
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

func ListResourceDefinitions() ([]string, error) {
	out, err := run("get", "CustomResourceDefinition", "-o=custom-columns=KIND:spec.names.plural,GROUP:.spec.group,SCOPE:.spec.scope", "--no-headers")
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	resources := make([]string, 0, 0)
	for scanner.Scan() {
		result := strings.Split(scanner.Text(), " ")
		if result[2] == "Namespaced" {
			resources = append(resources, result[0]+"."+result[1])
		}
	}
	resources = append(resources, nativeResources...)
	return resources, nil
}
