package linux

import (
	"bufio"
	"fmt"
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/pkg/errors"
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
		if strings.HasPrefix(text, "# Starting the map of external accesses to the ") {
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

func (h *HostsFile) Add(context string, ingresses []*kubectl.Ingress, gateways []*kubectl.IstioGateway) (err error) {
	if len(ingresses) < 1 && len(gateways) < 1 {
		return nil
	}

	hosts := make(map[string]*hostObj, 0)

	for _, ing := range ingresses {
		n := ing.Metadata.Name
		ns := ing.Metadata.Namespace

		for _, rule := range ing.Spec.Rules {
			name := rule.Host
			eip := ing.ExternalIP()
			if len(eip) > 0 && len(name) > 0 && name != "*" {
				hosts[name] = &hostObj{
					name:        name,
					address:     eip,
					description: fmt.Sprintf("Ingress %s/%s", ns, n),
				}
			}
		}
	}

	var istioIngress string
	if len(gateways) > 0 {
		istioIngress, err = kubectl.IstioIngress()
		if err != nil {
			return err
		}
	}
	for _, gtw := range gateways {
		n := gtw.Metadata.Name
		ns := gtw.Metadata.Namespace

		for _, svr := range gtw.Spec.Servers {
			for _, name := range svr.Hosts {
				if len(istioIngress) > 0 && len(name) > 0 && name != "*" {
					hosts[name] = &hostObj{
						name:        name,
						address:     istioIngress,
						description: fmt.Sprintf("Istio Gateway %s/%s", ns, n),
					}
				}
			}
		}
	}

	var sb strings.Builder
	for _, host := range hosts {
		if sb.Len() == 0 {
			sb.WriteString(fmt.Sprintf("# Starting the map of external accesses to the %s\n", context))
		}
		sb.WriteString(fmt.Sprintf("%s\t\t%s # %s\n", host.address, host.name, host.description))
	}
	if sb.Len() > 0 {
		sb.WriteString(fmt.Sprintf("# Ending the map of external accesses to the %s\n\n", context))
	}

	h.ingresses = sb.String()
	return nil
}

func (h *HostsFile) Write() error {
	content := []byte(h.String())
	log.Debugf("\n%s", content)
	err := ioutil.WriteFile("/etc/hosts", content, 0644)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (h *HostsFile) String() string {
	last := h.contents[len(h.contents)-1]
	if len(last) > 1 {
		h.contents = append(h.contents, "\n")
	}
	return strings.Join(h.contents, "\n") + h.ingresses
}

type hostObj struct {
	name        string
	address     string
	description string
}
