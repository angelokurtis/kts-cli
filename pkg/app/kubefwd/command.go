package kubefwd

import (
	"fmt"
	"strings"
)

type Command struct {
	Labels     map[string][]string
	Namespaces []string
}

func NewCommand(labels map[string][]string, namespaces []string) *Command {
	return &Command{Labels: labels, Namespaces: namespaces}
}

func (c *Command) String() string {
	var nb strings.Builder
	for _, namespace := range c.Namespaces {
		_, _ = fmt.Fprintf(&nb, " -n %s", namespace)
	}

	var lb strings.Builder

	count := 0
	for key, values := range c.Labels {
		if count > 0 {
			_, _ = fmt.Fprintf(&lb, ",")
		}

		if len(values) > 1 {
			_, _ = fmt.Fprintf(&lb, "%s in (", key)

			for i, value := range values {
				if i == 0 {
					_, _ = fmt.Fprintf(&lb, "%s", value)
				} else {
					_, _ = fmt.Fprintf(&lb, ",%s", value)
				}
			}

			_, _ = fmt.Fprintf(&lb, ")")
		} else {
			_, _ = fmt.Fprintf(&lb, "%s=%s", key, values[0])
		}

		count++
	}

	return "sudo -E kubefwd svc" + nb.String() + " -l \"" + lb.String() + "\""
}
