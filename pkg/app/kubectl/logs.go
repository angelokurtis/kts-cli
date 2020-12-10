package kubectl

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/internal/colors"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

func Logs(container *Container, since string, previous bool) ([]byte, error) {
	p := container.Pod
	ns := container.Namespace
	c := container.Name

	var cmd string
	if container.Single {
		cmd = fmt.Sprintf("kubectl logs %s --since=%s -n %s", p, since, ns)
	} else {
		cmd = fmt.Sprintf("kubectl logs %s -c %s --since=%s -n %s", p, c, since, ns)
	}
	if previous {
		cmd += " --previous"
	}

	logs, err := bash.RunAndLogRead(cmd)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func SaveLogs(containers *Containers, since string, previous bool) {
	for _, container := range containers.Items {
		err := saveLog(container, since, previous)
		if err != nil {
			color.Yellow.Println("[WARN] " + err.Error())
		}
	}
}

func saveLog(container *Container, since string, previous bool) error {
	p := container.Pod
	ns := container.Namespace
	c := container.Name

	dir := fmt.Sprintf("./logs/%s/pods/%s", ns, p)
	filename := fmt.Sprintf("%s/%s.log", dir, c)

	logs, err := Logs(container, since, previous)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}
	err = ioutil.WriteFile(filename, colors.Strip(logs), os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
