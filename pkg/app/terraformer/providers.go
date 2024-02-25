package terraformer

import (
	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

func ListProviders() Providers {
	return []string{
		// "alicloud",
		"aws",
		// "azure",
		// "cloudflare",
		// "commercetools",
		// "datadog",
		// "digitalocean",
		// "fastly",
		// "github",
		// "gmailfilter",
		"google",
		// "heroku",
		// "keycloak",
		"kubernetes",
		// "linode",
		// "logzio",
		// "mikrotik",
		// "newrelic",
		// "ns1",
		// "octopusdeploy",
		// "openstack",
		// "plan",
		// "rabbitmq",
		// "vultr",
		// "yandex",
	}
}

func (p Providers) SelectProvider() (string, error) {
	var selected string

	if len(p) == 0 {
		return "", nil
	} else if len(p) > 1 {
		prompt := &survey.Select{
			Message: "Select the Terraformer Provider:",
			Options: p,
		}

		err := survey.AskOne(prompt, &selected, survey.WithPageSize(10), survey.WithKeepFilter(true))
		if err != nil {
			return "", errors.WithStack(err)
		}
	} else {
		selected = p[0]
	}

	return selected, nil
}

type Providers []string

type Provider string
