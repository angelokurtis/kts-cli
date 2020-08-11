package kubectl

import (
	"bufio"
	"bytes"
	"github.com/AlecAivazis/survey/v2"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/pkg/errors"
	"strings"
)

func ListResources(resources string, allNamespaces bool) ([]string, error) {
	cmd := []string{"get", resources}
	if allNamespaces {
		cmd = append(cmd, "--all-namespaces")
	}
	out, err := runAndLogRead(cmd...)
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

func SelectResources(resources string, allNamespaces bool) ([]*resource, error) {
	cmd := "kubectl get " + resources + " -o=jsonpath='{.items[*].metadata.selfLink}'"
	if allNamespaces {
		cmd = cmd + " --all-namespaces"
	}
	out, err := bash.RunAndLogRead(cmd)
	if err != nil {
		return nil, err
	}

	links := make(map[string]*resource, 0)
	var options []string
	for _, value := range strings.Fields(string(out)) {
		r, err := newResource(value)
		if err != nil {
			return nil, err
		}
		key := ""
		if allNamespaces {
			key = key + r.Namespace + "/"
		}
		key = key + r.Kind + "/" + r.Name
		links[key] = r
		options = append(options, key)
	}

	var selects []string
	prompt := &survey.MultiSelect{
		Message: "Select the resource:",
		Options: options,
	}

	err = survey.AskOne(prompt, &selects, survey.WithPageSize(10))
	if err != nil {
		return nil, err
	}

	res := make([]*resource, 0, len(selects))
	for _, s := range selects {
		res = append(res, links[s])
	}

	return res, nil
}

type resource struct {
	Name      string
	Kind      string
	Namespace string
}

func newResource(l string) (*resource, error) {
	splitted := strings.Split(l, "/")
	size := len(splitted)
	if size == 8 {
		return &resource{
			Name:      splitted[7],
			Kind:      splitted[6] + "." + splitted[2],
			Namespace: splitted[5],
		}, nil
	} else if size == 7 {
		return &resource{
			Name:      splitted[6],
			Kind:      splitted[5],
			Namespace: splitted[4],
		}, nil
	}
	return nil, errors.New("unrecognized selfLink format")
}
