package golang

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/gookit/color"
	"github.com/pkg/errors"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func UnmarshalPackage(data []byte) (Package, error) {
	var r Package
	err := json.Unmarshal(data, &r)

	return r, err
}

type Package struct {
	Dir         string   `json:"Dir"`
	ImportPath  string   `json:"ImportPath"`
	Name        string   `json:"Name"`
	Target      string   `json:"Target"`
	Root        string   `json:"Root"`
	Module      Module   `json:"Module"`
	Match       []string `json:"Match"`
	Stale       bool     `json:"Stale"`
	StaleReason string   `json:"StaleReason"`
	GoFiles     []string `json:"GoFiles"`
	TestGoFiles []string `json:"TestGoFiles"`
	Imports     []string `json:"Imports"`
	Deps        []string `json:"Deps"`
}

func (p *Package) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Package) RelativeImportPath() string {
	if p.ImportPath == p.Module.Path {
		return "/"
	}

	return strings.ReplaceAll(p.ImportPath, p.Module.Path, ".")
}

func (p *Package) InternalImports() []string {
	i := make([]string, 0)

	for _, s := range p.Imports {
		if strings.HasPrefix(s, p.Module.Path) && !strings.HasPrefix(s, p.ImportPath) {
			i = append(i, strings.ReplaceAll(s, p.Module.Path, "."))
		}
	}

	sort.Strings(i)

	return i
}

func (p *Package) ImportsOf(dep string) []string {
	i := make([]string, 0)

	for _, s := range p.Imports {
		if strings.HasPrefix(s, dep) {
			i = append(i, s)
		}
	}

	return i
}

type Module struct {
	Path      string `json:"Path"`
	Main      bool   `json:"Main"`
	Dir       string `json:"Dir"`
	GoMod     string `json:"GoMod"`
	GoVersion string `json:"GoVersion"`
}

type Packages []*Package

func (p Packages) Deps() []string {
	deps := make([]string, 0, 0)
	for _, pkg := range p {
		deps = dedupe(deps, pkg.Deps...)
	}

	return deps
}

func (p Packages) Usages(dep string) Packages {
	owners := make([]*Package, 0, 0)

	for _, pkg := range p {
		if contains(pkg.Imports, dep) {
			owners = append(owners, pkg)
		}
	}

	return owners
}

func DescribePackage(dir string) (*Package, error) {
	j, err := bash.Run(fmt.Sprintf("cd %s && go list -json", dir))
	if err != nil {
		return nil, err
	}

	pkg, err := UnmarshalPackage(j)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &pkg, nil
}

func ListPackages(dir string) (Packages, error) {
	cmd := "cd " + dir + ` && go list -json ./... | jq -s .`
	color.Primary.Println(cmd)

	j, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var pkg []*Package
	if err = json.Unmarshal(j, &pkg); err != nil {
		return nil, errors.WithStack(err)
	}

	return pkg, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.HasPrefix(a, e) {
			return true
		}
	}

	return false
}
