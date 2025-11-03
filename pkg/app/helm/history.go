package helm

import (
	"encoding/json"
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func GetHistory(release string, options ...OptionFunc) (Revisions, error) {
	o := new(Option)
	if err := o.apply(options...); err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("helm history %s -o json", release)
	if o.Namespace != "" {
		cmd += " -n " + o.Namespace
	}

	out, err := bash.RunAndLogRead(cmd)
	if err != nil {
		return nil, err
	}

	return UnmarshalHistory(out)
}

type Revisions []*Revision

func UnmarshalHistory(data []byte) (Revisions, error) {
	var r Revisions
	err := json.Unmarshal(data, &r)

	return r, errors.WithStack(err)
}

func (r *Revisions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r Revisions) SelectMany() (Revisions, error) {
	if len(r) == 0 {
		return make([]*Revision, 0, 0), nil
	} else if len(r) == 1 {
		return []*Revision{r[0]}, nil
	}

	numbers := r.Numbers()
	prompt := &survey.MultiSelect{
		Message: "Select the Helm release revisions:",
		Options: numbers,
	}

	var selects []string

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return GetRevisions(r, selects...), nil
}

func (r Revisions) Numbers() []string {
	numbers := make([]string, 0, len(r))
	for _, revision := range r {
		numbers = append(numbers, fmt.Sprintf("%d", revision.Number))
	}

	return numbers
}

func GetRevisions[T ~int64 | ~string](r Revisions, numbers ...T) []*Revision {
	revisions := make([]*Revision, 0, len(numbers))
	for _, n := range numbers {
		for _, revision := range r {
			var match bool
			switch v := any(n).(type) {
			case int64:
				match = revision.Number == v
			case string:
				match = fmt.Sprintf("%d", revision.Number) == v
			}
			if match {
				revisions = append(revisions, revision)
			}
		}
	}

	return revisions
}

type Revision struct {
	Number      int64  `json:"revision"`
	Updated     string `json:"updated"`
	Status      string `json:"status"`
	Chart       string `json:"chart"`
	AppVersion  string `json:"app_version"`
	Description string `json:"description"`
}
