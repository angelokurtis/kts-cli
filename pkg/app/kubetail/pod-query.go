package kubetail

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/pkg/errors"
	"strings"
)

func PodQuery(containers *kubectl.Containers, since string) error {
	namespaces := containers.Namespaces()
	if len(namespaces) > 1 {
		return errors.New("kubetail do not support multiple namespaces")
	}
	ns := namespaces[0]
	c := containers.Names()
	p := containers.Pods()
	cmd := fmt.Sprintf("kubetail -k pod -n %s -c '%s' '%s' --since %s", ns, strings.Join(c, "|"), strings.Join(p, "|"), since)
	return bash.Follow(cmd)
}
