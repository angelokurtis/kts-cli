package linux

import (
	"bufio"
	"fmt"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"io/ioutil"
	"os"
	"strings"
)

type HostsFile struct {
	contents  []string
	ingresses string
}

func LoadHostsFile() (*HostsFile, error) {
	file, err := os.Open("/etc/hosts")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	contents := make([]string, 0, 0)
	var ib strings.Builder

	var ingressBlock bool

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "\n# Starting the map of external accesses to the ") {
			ingressBlock = true
		}
		if ingressBlock {
			ib.WriteString(text)
			ib.WriteString("\n")
		} else {
			contents = append(contents, text)
		}
		if strings.HasPrefix(text, "# Ending the map of external accesses to the ") {
			ingressBlock = false
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &HostsFile{contents: contents, ingresses: ib.String()}, nil
}

func (h *HostsFile) Add(context string, ingresses *kubectl.Ingresses) {
	if len(ingresses.Items) < 1 {
		return
	}

	var sb strings.Builder
	for _, ing := range ingresses.Items {
		n := ing.Metadata.Name
		ns := ing.Metadata.Namespace

		for _, ingress := range ing.Status.LoadBalancer.Ingresses {
			for _, rule := range ing.Spec.Rules {
				host := rule.Host
				if len(ingress.IP) > 0 && len(host) > 0 {
					if sb.Len() == 0 {
						sb.WriteString(fmt.Sprintf("\n# Starting the map of external accesses to the %s\n", context))
					}
					sb.WriteString(fmt.Sprintf("%s\t\t%s # Ingress %s/%s\n", ingress.IP, host, ns, n))
				}
			}
		}
	}
	if sb.Len() > 0 {
		sb.WriteString(fmt.Sprintf("# Ending the map of external accesses to the '%s'\n\n", context))
	}
	h.ingresses = sb.String()
}

func (h *HostsFile) Write() error {
	content := []byte(h.String())
	err := ioutil.WriteFile("/etc/hosts", content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (h *HostsFile) String() string {
	last := h.contents[len(h.contents)-1]
	if len(last) > 1 {
		h.contents = append(h.contents, "")
	}
	return strings.Join(h.contents, "\n") + h.ingresses
}
