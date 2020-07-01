package portforward

import (
	"fmt"
	"strconv"
	"strings"
)

type Command struct {
	Labels    map[string][]string
	Namespace string
	Port      int
}

func NewCommand(namespace string, labels map[string][]string, port int) *Command {
	return &Command{Labels: labels, Namespace: namespace, Port: port}
}

func (c *Command) String() string {
	var nb strings.Builder
	_, _ = fmt.Fprintf(&nb, " -n %s", c.Namespace)
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

	port := strconv.Itoa(c.Port)
	return "kubectl port-forward" + nb.String() + " -l \"" + lb.String() + "\" " + port + ":" + port
}
