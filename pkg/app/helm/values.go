package helm

import (
	"fmt"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func GetSuppliedValues(release string, revision int64, options ...OptionFunc) ([]byte, error) {
	o := new(Option)
	if err := o.apply(options...); err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("helm get values %s --revision %d", release, revision)
	if o.Namespace != "" {
		cmd += " -n " + o.Namespace
	}
	return bash.RunAndLogRead(cmd)
}

func GetComputedValues(release string, revision int64, options ...OptionFunc) ([]byte, error) {
	o := new(Option)
	if err := o.apply(options...); err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("helm get values %s --revision %d --all ", release, revision)
	if o.Namespace != "" {
		cmd += " -n " + o.Namespace
	}
	return bash.RunAndLogRead(cmd)
}
