package golang

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/pkg/errors"
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

type Module struct {
	Path      string `json:"Path"`
	Main      bool   `json:"Main"`
	Dir       string `json:"Dir"`
	GoMod     string `json:"GoMod"`
	GoVersion string `json:"GoVersion"`
}

func ListPackages(dir string) (*Package, error) {
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
