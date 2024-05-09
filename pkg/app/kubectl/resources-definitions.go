package kubectl

import (
	"sort"
	"strconv"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/angelokurtis/kts-cli/internal/kube"
)

var emptyChar int32 = 32

const ignoreEventsResource = true

func ListResourceDefinitions() (*ResourcesDefinitions, error) {
	resources := make([]*ResourceDefinition, 0)

	lists, err := kube.GetDiscoveryClient().ServerPreferredResources()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, list := range lists {
		if len(list.APIResources) == 0 {
			continue
		}

		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			continue
		}

		for _, rsrc := range list.APIResources {
			if len(rsrc.Verbs) == 0 {
				continue
			}

			if rsrc.Name == "events" && (gv.Group == "" || gv.Group == "events.k8s.io") && gv.Version == "v1" {
				continue
			}

			resources = append(resources, &ResourceDefinition{
				Name:       rsrc.Name,
				ShortNames: rsrc.ShortNames,
				APIGroup:   gv.Group,
				Namespaced: rsrc.Namespaced,
				Kind:       rsrc.Kind,
				Verbs:      rsrc.Verbs,
			})
		}
	}

	return &ResourcesDefinitions{Items: resources}, nil
}

type ResourcesDefinitions struct {
	Items []*ResourceDefinition
}

func (r *ResourcesDefinitions) add(rd *ResourceDefinition) {
	r.Items = append(r.Items, rd)
}

func (r *ResourcesDefinitions) FilterVerbs(verb string) *ResourcesDefinitions {
	definitions := &ResourcesDefinitions{}

	for _, definition := range r.Items {
		if contains(definition.Verbs, verb) {
			definitions.add(definition)
		}
	}

	return definitions
}

func (r *ResourcesDefinitions) FilterNamespaced() *ResourcesDefinitions {
	definitions := &ResourcesDefinitions{}

	for _, definition := range r.Items {
		if definition.Namespaced {
			definitions.add(definition)
		}
	}

	return definitions
}

func (r *ResourcesDefinitions) FilterAPIGroup(group string) *ResourcesDefinitions {
	definitions := &ResourcesDefinitions{}

	for _, definition := range r.Items {
		if strings.Contains(definition.APIGroup, group) {
			definitions.add(definition)
		}
	}

	return definitions
}

func (r *ResourcesDefinitions) SelectGroups() ([]string, error) {
	groups := make([]string, 0, 0)
	for _, item := range r.Items {
		groups = dedupeStr(groups, item.APIGroup)
	}

	var selects []string

	prompt := &survey.MultiSelect{
		Message: "Select the APIGroups:",
		Options: groups,
	}

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return selects, nil
}

func (r *ResourcesDefinitions) APIGroups() []string {
	m := make(map[string]string, 0)

	for _, item := range r.Items {
		k := reverse(item.APIGroup)
		m[k] = item.APIGroup
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	groups := make([]string, 0, len(m))
	for _, k := range keys {
		groups = append(groups, m[k])
	}

	return groups
}

func (r *ResourcesDefinitions) String() string {
	sb := strings.Builder{}
	for i, resource := range r.Items {
		sb.WriteString(resource.String())

		if i < len(r.Items)-1 {
			sb.WriteString(",")
		}
	}

	return sb.String()
}

type ResourceDefinition struct {
	Name       string
	ShortNames []string
	APIGroup   string
	Namespaced bool
	Kind       string
	Verbs      []string
}

func NewResourceDefinition(line string, indexes []int) (*ResourceDefinition, error) {
	var name, shortNames, apiGroup, kind, verbs string

	var namespaced bool

	for i, index := range indexes {
		switch i {
		case 0:
			name = strings.TrimSpace(line[:index])
		case 1:
			shortNames = strings.TrimSpace(line[indexes[i-1]:index])
		case 2:
			g := strings.TrimSpace(line[indexes[i-1]:index])
			if g == "v1" {
				apiGroup = ""
			} else if strings.Contains(g, "/") {
				apiGroup = strings.Split(g, "/")[0]
			} else {
				apiGroup = g
			}
		case 3:
			var err error

			namespaced, err = strconv.ParseBool(strings.TrimSpace(line[indexes[i-1]:index]))
			if err != nil {
				return nil, errors.WithStack(err)
			}
		case 4:
			kind = strings.TrimSpace(line[indexes[i-1]:index])
			verbs = strings.TrimSpace(line[index:])
			verbs = strings.Replace(verbs, "[", "", -1)
			verbs = strings.Replace(verbs, "]", "", -1)
		}
	}

	return &ResourceDefinition{
		Name:       name,
		ShortNames: strings.Split(shortNames, ","),
		APIGroup:   apiGroup,
		Namespaced: namespaced,
		Kind:       kind,
		Verbs:      strings.Split(verbs, " "),
	}, nil
}

func (r *ResourceDefinition) String() string {
	if r.APIGroup == "" {
		return r.Name
	}

	return r.Name + "." + r.APIGroup
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func reverse(input string) string {
	n := 0

	runes := make([]rune, len(input))
	for _, r := range input {
		runes[n] = r
		n++
	}

	runes = runes[0:n]
	// Reverse
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}
	// Convert back to UTF-8.
	return string(runes)
}
