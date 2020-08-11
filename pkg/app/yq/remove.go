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
