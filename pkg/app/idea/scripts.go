package idea

import (
	"os"
	"path/filepath"
)

// listJetBrainsScripts lists files in ~/.local/share/JetBrains/Toolbox/scripts
func listJetBrainsScripts() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, ".local", "share", "JetBrains", "Toolbox", "scripts")

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var scriptFiles []string

	for _, file := range files {
		if !file.IsDir() {
			scriptFiles = append(scriptFiles, file.Name())
		}
	}

	return scriptFiles, nil
}

func chooseJetBrainsScript(scripts []string) string {
	for _, name := range []string{"idea", "goland"} {
		for _, script := range scripts {
			if script == name {
				return name
			}
		}
	}

	return ""
}
