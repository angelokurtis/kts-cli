package helm

import (
	"fmt"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func GetNotes(release string, revision int64, options ...OptionFunc) ([]byte, error) {
	o := new(Option)
	if err := o.apply(options...); err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("helm get notes %s --revision %d", release, revision)
	if o.Namespace != "" {
		cmd += " -n " + o.Namespace
	}
	return bash.RunAndLogRead(cmd)
}
