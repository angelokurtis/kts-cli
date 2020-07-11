package kubectl

import (
	"bufio"
	"bytes"
)

func CurrentContext() (string, error) {
	out, err := run("config", "current-context")
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", nil
}
