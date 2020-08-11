package kubectl

import (
	"strings"
)

func MeshesHosts() ([]string, error) {
	out, err := runAndLogRead("get", "--all-namespaces", "meshes.management.sensedia.com", "-o=jsonpath='{.items[*].spec.host}'")
	if err != nil {
		return nil, err
	}
	str := string(out)
	str = str[1:]
	str = str[:len(str)-1]

	meshes := strings.Split(str, " ")
	return meshes, nil
}
