package istioctl

import (
	"encoding/json"
	"fmt"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/angelokurtis/kts-cli/pkg/app/yq"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"io/ioutil"
	"strconv"
	"strings"
	"github.com/yalp/jsonpath"
)

func AddToMesh(deployment *kubectl.Deployment) error {
	name := deployment.Metadata.Name
	namespace := deployment.Metadata.Namespace
	_, err := bash.RunAndLogWrite(fmt.Sprintf("istioctl experimental add-to-mesh deployment %s -n %s", name, namespace))
	if err != nil {
		if !strings.Contains(err.Error(), "0 errors occurred") {
			return err
		}
	}
	return nil
}

func RemoveFromMesh(deployment *kubectl.Deployment) error {
	name := deployment.Metadata.Name
	namespace := deployment.Metadata.Namespace
	_, err := bash.RunAndLogWrite(fmt.Sprintf("istioctl experimental remove-from-mesh deployment %s -n %s", name, namespace))
	if err != nil {
		if !strings.Contains(err.Error(), "0 errors occurred") {
			return err
		}
	}
	return nil
}

func KubeInject(dep *kubectl.Deployment) error {
	n := dep.Metadata.Name
	ns := dep.Metadata.Namespace
	cmd := fmt.Sprintf("kubectl get deployment %s -n %s -o yaml | yq d - spec.template.metadata.annotations[sidecar.istio.io/inject] | istioctl kube-inject -f - ", n, ns)
	out, err := bash.Run(cmd)
	if err != nil {
		return err
	}

	yamlFile := n + ".injected.yaml"
	yamlPath := "./" + ns + "/apps/deployments"

	_, err = bash.Run("mkdir -p " + yamlPath)
	if err != nil {
		return err
	}

	color.Primary.Println(cmd + " > " + yamlPath + "/" + yamlFile)
	if err = ioutil.WriteFile(yamlPath+"/"+yamlFile, out, 0644); err != nil {
		return errors.WithStack(err)
	}

	if err := yq.DeleteKubernetesNodes(yamlPath + "/" + yamlFile); err != nil {
		return err
	}

	return nil
}

func KubeUninject(dep *kubectl.Deployment) error {
	n := dep.Metadata.Name
	ns := dep.Metadata.Namespace
	cmd := fmt.Sprintf("kubectl get deployment %s -n %s -o yaml | istioctl x kube-uninject -f - ", n, ns)
	out, err := bash.Run(cmd)
	if err != nil {
		return err
	}

	yamlFile := n + ".uninjected.yaml"
	yamlPath := "./" + ns + "/apps/deployments"

	_, err = bash.Run("mkdir -p " + yamlPath)
	if err != nil {
		return err
	}

	color.Primary.Println(cmd + " > " + yamlPath + "/" + yamlFile)
	if err = ioutil.WriteFile(yamlPath+"/"+yamlFile, out, 0644); err != nil {
		return errors.WithStack(err)
	}

	err = rollbackPrometheusAnnotations(dep, yamlPath+"/"+yamlFile)
	if err != nil {
		return err
	}

	err = rollbackProbes(dep, yamlPath+"/"+yamlFile)
	if err != nil {
		return err
	}

	if err := yq.DeleteKubernetesNodes(yamlPath + "/" + yamlFile); err != nil {
		return err
	}
	if err := deleteInjectedNodes(yamlPath + "/" + yamlFile); err != nil {
		return err
	}
	return nil
}

func rollbackProbes(dep *kubectl.Deployment, yamlPath string) error {
	name := dep.Metadata.Name
	container := dep.GetContainer("istio-proxy")
	readyz, err := NewReadyzHTTPGet(name, container.GetEnv("ISTIO_KUBE_APP_PROBERS"))
	if err != nil {
		return err
	}
	livez, err := NewLivezHTTPGet(name, container.GetEnv("ISTIO_KUBE_APP_PROBERS"))
	if err != nil {
		return err
	}

	for i, c := range dep.Spec.Template.Spec.Containers {
		containerPath := "spec.template.spec.containers[" + strconv.Itoa(i) + "]"
		found, err := yq.ReadNodeValue(yamlPath, containerPath+".name")
		if err != nil {
			return err
		}
		if found == name {
			if readyz != nil {
				if readyz.Path != "" {
					if err := yq.UpdateNode(yamlPath, containerPath+".readinessProbe.httpGet.path", readyz.Path); err != nil {
						return err
					}
				}
				if readyz.Port != 0 {
					if err := yq.UpdateNode(yamlPath, containerPath+".readinessProbe.httpGet.port", strconv.Itoa(readyz.Port)); err != nil {
						return err
					}
				}
				if readyz.Scheme != "" {
					if err := yq.UpdateNode(yamlPath, containerPath+".readinessProbe.httpGet.scheme", readyz.Scheme); err != nil {
						return err
					}
				}
			}
			if livez != nil {
				if readyz.Path != "" {
					if err := yq.UpdateNode(yamlPath, containerPath+".livenessProbe.httpGet.path", livez.Path); err != nil {
						return err
					}
				}
				if readyz.Port != 0 {
					if err := yq.UpdateNode(yamlPath, containerPath+".livenessProbe.httpGet.port", strconv.Itoa(livez.Port)); err != nil {
						return err
					}
				}
				if readyz.Scheme != "" {
					if err := yq.UpdateNode(yamlPath, containerPath+".livenessProbe.httpGet.scheme", livez.Scheme); err != nil {
						return err
					}
				}
			}

			failureThreshold := c.ReadinessProbe.FailureThreshold
			if failureThreshold != nil {
				if err := yq.UpdateNode(yamlPath, containerPath+".readinessProbe.failureThreshold", strconv.Itoa(*failureThreshold)); err != nil {
					return err
				}
			}
			initialDelay := c.ReadinessProbe.InitialDelaySeconds
			if initialDelay != nil {
				if err := yq.UpdateNode(yamlPath, containerPath+".readinessProbe.initialDelaySeconds", strconv.Itoa(*initialDelay)); err != nil {
					return err
				}
			}
			period := c.ReadinessProbe.PeriodSeconds
			if period != nil {
				if err := yq.UpdateNode(yamlPath, containerPath+".readinessProbe.periodSeconds", strconv.Itoa(*period)); err != nil {
					return err
				}
			}
			successThreshold := c.ReadinessProbe.SuccessThreshold
			if successThreshold != nil {
				if err := yq.UpdateNode(yamlPath, containerPath+".readinessProbe.successThreshold", strconv.Itoa(*successThreshold)); err != nil {
					return err
				}
			}
			timeout := c.ReadinessProbe.TimeoutSeconds
			if timeout != nil {
				if err := yq.UpdateNode(yamlPath, containerPath+".readinessProbe.timeoutSeconds", strconv.Itoa(*timeout)); err != nil {
					return err
				}
			}

			failureThreshold = c.LivenessProbe.FailureThreshold
			if failureThreshold != nil {
				if err = yq.UpdateNode(yamlPath, containerPath+".livenessProbe.failureThreshold", strconv.Itoa(*failureThreshold)); err != nil {
					return err
				}
			}
			initialDelay = c.LivenessProbe.InitialDelaySeconds
			if initialDelay != nil {
				if err = yq.UpdateNode(yamlPath, containerPath+".livenessProbe.initialDelaySeconds", strconv.Itoa(*initialDelay)); err != nil {
					return err
				}
			}
			period = c.LivenessProbe.PeriodSeconds
			if period != nil {
				if err = yq.UpdateNode(yamlPath, containerPath+".livenessProbe.periodSeconds", strconv.Itoa(*period)); err != nil {
					return err
				}
			}
			successThreshold = c.LivenessProbe.SuccessThreshold
			if successThreshold != nil {
				if err = yq.UpdateNode(yamlPath, containerPath+".livenessProbe.successThreshold", strconv.Itoa(*successThreshold)); err != nil {
					return err
				}
			}
			timeout = c.LivenessProbe.TimeoutSeconds
			if timeout != nil {
				if err = yq.UpdateNode(yamlPath, containerPath+".livenessProbe.timeoutSeconds", strconv.Itoa(*timeout)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func rollbackPrometheusAnnotations(dep *kubectl.Deployment, yamlPath string) error {
	container := dep.GetContainer("istio-proxy")
	cfg, err := NewIstioPrometheusAnnotations(container.GetEnv("ISTIO_PROMETHEUS_ANNOTATIONS"))
	if err != nil {
		return err
	}
	if cfg != nil {
		if cfg.Path != "" {
			if err := yq.UpdateNodeWithQuotes(yamlPath, "spec.template.metadata.annotations[prometheus.io/path]", cfg.Path); err != nil {
				return err
			}
		}
		if cfg.Port != "" {
			if err := yq.UpdateNodeWithQuotes(yamlPath, "spec.template.metadata.annotations[prometheus.io/port]", cfg.Port); err != nil {
				return err
			}
		}
		if cfg.Scrape != "" {
			if err := yq.UpdateNodeWithQuotes(yamlPath, "spec.template.metadata.annotations[prometheus.io/scrape]", cfg.Scrape); err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteInjectedNodes(manifestPath string) error {
	if err := yq.DeleteNode(manifestPath, "spec.template.metadata.annotations[sidecar.istio.io/inject]"); err != nil {
		return err
	}
	if err := yq.DeleteNode(manifestPath, "spec.template.metadata.labels[istio.io/rev]"); err != nil {
		return err
	}
	if err := yq.DeleteNode(manifestPath, "spec.template.metadata.labels[security.istio.io/tlsMode]"); err != nil {
		return err
	}
	if err := yq.DeleteNode(manifestPath, "spec.template.metadata.labels[service.istio.io/canonical-name]"); err != nil {
		return err
	}
	if err := yq.DeleteNode(manifestPath, "spec.template.metadata.labels[service.istio.io/canonical-revision]"); err != nil {
		return err
	}
	return nil
}

type HTTPGet struct {
	Path   string `json:"path"`
	Port   int    `json:"port"`
	Scheme string `json:"scheme"`
}

func NewLivezHTTPGet(svc string, data string) (*HTTPGet, error) {
	if data == "" {
		return nil, nil
	}
	var probes interface{}
	err := json.Unmarshal([]byte(data), &probes)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	livez, err := jsonpath.Read(probes, "$[\"/app-health/"+svc+"/livez\"][\"httpGet\"]")
	if err != nil {
		if strings.HasPrefix(err.Error(), "no key ") {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	jsonString, err := json.Marshal(livez)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	probe := &HTTPGet{}
	err = json.Unmarshal(jsonString, probe)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return probe, nil
}

func NewReadyzHTTPGet(svc string, data string) (*HTTPGet, error) {
	if data == "" {
		return nil, nil
	}
	var probes interface{}
	err := json.Unmarshal([]byte(data), &probes)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	readyz, err := jsonpath.Read(probes, "$[\"/app-health/"+svc+"/readyz\"][\"httpGet\"]")
	if err != nil {
		if strings.HasPrefix(err.Error(), "no key ") {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	jsonString, err := json.Marshal(readyz)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	probe := &HTTPGet{}
	err = json.Unmarshal(jsonString, probe)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return probe, nil
}

type IstioPrometheusCfg struct {
	Scrape string `json:"scrape"`
	Path   string `json:"path"`
	Port   string `json:"port"`
}

func NewIstioPrometheusAnnotations(s string) (*IstioPrometheusCfg, error) {
	if s == "" {
		return nil, nil
	}
	data := &IstioPrometheusCfg{}
	err := json.Unmarshal([]byte(s), data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}
