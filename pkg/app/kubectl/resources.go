package kubectl

import (
	"bufio"
	"bytes"
)

func ListResources(resources string, allNamespaces bool) ([]string, error) {
	cmd := []string{"get", resources}
	if allNamespaces {
		cmd = append(cmd, "--all-namespaces")
	}
	out, err := runAndLog(cmd...)
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

//func (s *Services) SelectResources(resources string, allNamespaces bool) (map[string][]string, error) {
//	cmd := []string{"get", resources}
//	if allNamespaces {
//		cmd = append(cmd, "--all-namespaces")
//	}
//	out, err := runAndLog(cmd...)
//	if err != nil {
//		return nil, err
//	}
//
//
//	var options []string
//	for key, values := range s.Labels() {
//		for _, value := range values {
//			options = append(options, key+"="+value)
//		}
//	}
//
//	var selects []string
//	prompt := &survey.MultiSelect{
//		Message: "Select the service labels:",
//		Options: options,
//	}
//
//	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10))
//	if err != nil {
//		return nil, err
//	}
//
//	labels := make(map[string][]string, 0)
//	for _, s := range selects {
//		spt := strings.Split(s, "=")
//		key := spt[0]
//		value := spt[len(spt)-1]
//
//		values := labels[key]
//		values = append(values, value)
//		labels[key] = values
//	}
//
//	return labels, nil
//}
