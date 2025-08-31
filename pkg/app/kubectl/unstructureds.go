package kubectl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func ListUnstructureds(resources, namespace string, allNamespaces bool) (*unstructured.UnstructuredList, error) {
	cmd := []string{"get", resources, "-o", "json"}
	if allNamespaces {
		cmd = append(cmd, "--all-namespaces")
	} else if namespace != "" {
		cmd = append(cmd, "-n", namespace)
	}

	out, err := runAndLogRead(cmd...)
	if err != nil {
		return nil, err
	}

	// Be defensive in case kubectl prints headers/noise before JSON
	idx := strings.IndexByte(string(out), '{')
	if idx == -1 {
		return nil, fmt.Errorf("no JSON object found in kubectl output")
	}

	var raw runtime.RawExtension
	raw.Raw = out[idx:]

	var list unstructured.UnstructuredList
	if err := json.NewDecoder(bytes.NewReader(raw.Raw)).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed to decode JSON into UnstructuredList: %w", err)
	}

	return &list, nil
}
