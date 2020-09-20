package helm

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/pkg/app/helm"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/gookit/color"
	changeCase "github.com/ku/go-change-case"
	"strings"
)

type Release struct {
	Chart     string `hcl:"chart"`
	Name      string `hcl:"name"`
	Namespace string `hcl:"namespace"`
}

func NewRelease() (*Release, error) {
	releases, err := helm.ListReleases()
	if err != nil {
		return nil, err
	}
	release, err := releases.SelectOne()
	if err != nil {
		return nil, err
	}
	return &Release{Chart: release.Chart, Name: release.Name, Namespace: release.Namespace}, nil
}

func (r *Release) GetType() string {
	return "helm_release"
}

func (r *Release) GetID() string {
	return fmt.Sprintf("%s.%s", r.GetType(), r.GetName())
}

func (r *Release) GetName() string {
	return changeCase.Snake(r.Name)
}

func (r *Release) Import() error {
	cmd := fmt.Sprintf("terraform import %s %s/%s", r.GetID(), r.Namespace, r.Name)
	if _, err := bash.RunAndLogWrite(cmd); err != nil {
		if strings.Contains(err.Error(), "Resource already managed by Terraform") {
			color.Yellow.Printf("[WARN] the %s is already managed by Terraform\n", r.GetID())
		} else {
			return err
		}
	}
	return nil
}
