package kubectl

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func SaveLogs(container *Container) error {
	p := container.Pod
	ns := container.Namespace
	c := container.Name

	dir := fmt.Sprintf("./logs/%s/%s", ns, p)
	filename := fmt.Sprintf("%s/%s.log", dir, c)
	cmd := fmt.Sprintf("kubectl logs %s -c %s -n %s", p, c, ns)

	color.Primary.Printf("%s > %s\n", cmd, filename)
	logs, err := bash.Run(cmd)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}
	err = ioutil.WriteFile(filename, logs, os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func saveLogs(pod string, container string, namespace string) error {
	dir := fmt.Sprintf("./logs/%s/%s", namespace, pod)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}
	//cmd := exec.Command("kubectl", "logs", pod, "-c", container, "-n", namespace, ">", dir+"/"+container+".log")
	cmd := exec.Command("kubectl", "logs", pod, "-c", container, "-n", namespace)
	color.Primary.Println(strings.Join(cmd.Args, " ") + " > " + dir + "/" + container + ".log")

	// open the out file for writing
	outfile, err := os.Create(dir + "/" + container + ".log")
	if err != nil {
		return errors.WithStack(err)
	}
	defer outfile.Close()
	cmd.Stdout = outfile

	err = cmd.Start()
	if err != nil {
		return errors.WithStack(err)
	}
	cmd.Wait()
	return nil
}
