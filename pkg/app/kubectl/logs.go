package kubectl

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/internal/color"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strings"
)

func SaveLogs(pod *Pod) error {
	pn := pod.Metadata.Name
	ns := pod.Metadata.Namespace

	if len(pod.Spec.Containers) == 1 {
		cn := pod.Spec.Containers[0].Name
		err := saveLogs(pn, cn, ns)
		if err != nil {
			return err
		}
		return nil
	}
	for _, container := range pod.Spec.Containers {
		cn := container.Name
		err := saveLogs(pn, cn, ns)
		if err != nil {
			return err
		}
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
	fmt.Printf(color.Notice, strings.Join(cmd.Args, " ")+" > "+dir+"/"+container+".log"+"\n")

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
