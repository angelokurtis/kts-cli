package kubectl

import (
	"encoding/json"
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func ListTLSSecrets() (*Secrets, error) {
	out, err := bash.RunAndLogRead("kubectl get secret --field-selector type=kubernetes.io/tls --all-namespaces -o=json")
	if err != nil {
		return nil, err
	}

	var secrets *Secrets
	if err := json.Unmarshal(out, &secrets); err != nil {
		return nil, errors.WithStack(err)
	}

	return secrets, nil
}

func ListSecrets() (*Secrets, error) {
	out, err := bash.RunAndLogRead("kubectl get secret --all-namespaces -o=json")
	if err != nil {
		return nil, err
	}

	var secrets *Secrets
	if err := json.Unmarshal(out, &secrets); err != nil {
		return nil, errors.WithStack(err)
	}

	return secrets, nil
}

func SearchSecret(label string) (*Secrets, error) {
	out, err := runAndLogRead("get", "Secret", "-o=json", "--all-namespaces", "-l", label)
	if err != nil {
		return nil, err
	}

	var secrets *Secrets
	if err := json.Unmarshal(out, &secrets); err != nil {
		return nil, errors.WithStack(err)
	}

	return secrets, nil
}

func GetSecretKeyValue(ref *KeyRef, namespace string) (string, error) {
	out, err := bash.RunAndLogRead(fmt.Sprintf("kubectl get secret --namespace %s %s -o jsonpath=\"{.data.%s}\" | base64 --decode", namespace, ref.Name, ref.Key))
	if err != nil {
		return "", err
	}

	return string(out), nil
}

type Secrets struct {
	Items []*Secret `json:"items"`
}

func (m *Secrets) Names() []string {
	secrets := m.Items
	names := make([]string, 0, len(secrets))

	for _, release := range secrets {
		names = append(names, release.Metadata.Namespace+"/"+release.Metadata.Name)
	}

	return names
}

func (m *Secrets) Get(name string) *Secret {
	for _, secret := range m.Items {
		if secret.Metadata.Namespace+"/"+secret.Metadata.Name == name {
			return secret
		}
	}

	return nil
}

func (m *Secrets) SelectMany() (*Secrets, error) {
	if len(m.Items) == 0 {
		return &Secrets{}, nil
	}

	names := m.FullNames()
	prompt := &survey.MultiSelect{
		Message: "Select the Secrets:",
		Options: names,
	}

	var selects []string

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	secrets := make([]*Secret, 0, len(selects))
	for _, name := range selects {
		secrets = append(secrets, m.Get(name))
	}

	return &Secrets{Items: secrets}, nil
}

func (m *Secrets) FullNames() []string {
	secrets := m.Items
	names := make([]string, 0, len(secrets))

	for _, release := range secrets {
		names = append(names, release.Metadata.Namespace+"/"+release.Metadata.Name)
	}

	return names
}

func (m *Secrets) SelectOne() (*Secret, error) {
	names := m.Names()

	if len(names) == 1 {
		return m.Get(names[0]), nil
	}

	var selected string

	prompt := &survey.Select{
		Message: "Select the Secret:",
		Options: names,
	}

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return m.Get(selected), nil
}

func (m *Secrets) SingleResult() (*Secret, error) {
	if len(m.Items) == 0 {
		return nil, nil
	}

	if len(m.Items) == 1 {
		return m.Items[0], nil
	}

	return nil, errors.New("found more than one Secret")
}

type Secret struct {
	Data     map[string]string `json:"data"`
	Kind     string            `json:"kind"`
	Metadata struct {
		Labels    map[string]string `json:"labels"`
		Name      string            `json:"name"`
		Namespace string            `json:"namespace"`
	} `json:"metadata"`
}
