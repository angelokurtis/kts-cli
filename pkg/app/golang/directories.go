package golang

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

func ListDirectories(dir string) ([]string, error) {
	folders := make([]string, 0)
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !strings.Contains(path, "/vendor/") && !strings.HasPrefix(path, "vendor/") && strings.HasSuffix(path, ".go") {
				folders = dedupe(folders, filepath.Dir(path))
			}
			return nil
		})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sort.Strings(folders)
	return folders, nil
}

func dedupe(a []string, b ...string) []string {
	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}
	for letter := range check {
		res = append(res, letter)
	}
	return res
}
