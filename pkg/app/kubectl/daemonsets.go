package kubectl

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"strings"
	"time"
)

func ListDaemonSets(namespace string, allNamespaces bool) (*DaemonSets, error) {
	cmd := []string{"get", "daemonsets", "-o=json"}
	if allNamespaces {
		cmd = append(cmd, "--all-namespaces")
	} else if namespace != "" {
		cmd = append(cmd, "-n", namespace)
	}

	out, err := run(cmd...)
	if err != nil {
		return nil, err
	}

	var daemonSets *DaemonSets
	if err = json.Unmarshal(out, &daemonSets); err != nil {
		return nil, errors.WithStack(err)
	}

	return daemonSets, nil
}

type DaemonSets struct {
	Items []*DaemonSet `json:"items"`
}

func (s *DaemonSets) FilterSuspended() (*DaemonSets, error) {
	if len(s.Items) == 0 {
		return &DaemonSets{}, nil
	}

	daemonSets := lo.Filter(s.Items, func(item *DaemonSet, index int) bool {
		nodeSelector, ok := item.Spec.Template.Spec.NodeSelector["non-existing"]
		return ok && nodeSelector == "true"
	})

	return &DaemonSets{Items: daemonSets}, nil
}

func (s *DaemonSets) SelectMany() (*DaemonSets, error) {
	if len(s.Items) == 0 {
		return &DaemonSets{}, nil
	}

	names := s.FullNames()
	prompt := &survey.MultiSelect{
		Message: "Select DaemonSets:",
		Options: names,
	}

	var selects []string

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	daemonSets := make([]*DaemonSet, 0, len(selects))
	for _, name := range selects {
		daemonSets = append(daemonSets, s.Get(name))
	}

	return &DaemonSets{Items: daemonSets}, nil
}

func (s *DaemonSets) FullNames() []string {
	items := s.Items
	names := make([]string, 0, len(items))

	for _, release := range items {
		names = append(names, release.Metadata.Namespace+"/"+release.Metadata.Name)
	}

	return names
}

func (s *DaemonSets) Get(name string) *DaemonSet {
	for _, deployment := range s.Items {
		if deployment.Metadata.Namespace+"/"+deployment.Metadata.Name == name {
			return deployment
		}
	}

	return nil
}

func (s *DaemonSets) Suspend() error {
	byNamespace := lo.GroupBy(s.Items, func(item *DaemonSet) string {
		return item.Metadata.Namespace
	})

	for namespace, daemonSets := range byNamespace {
		names := lo.Map(daemonSets, func(daemonSet *DaemonSet, index int) string {
			return daemonSet.Metadata.Name
		})

		_, err := bash.RunAndLogWrite(fmt.Sprintf(`kubectl -n %s patch daemonset %s -p '{"spec": {"template": {"spec": {"nodeSelector": {"non-existing": "true"}}}}}'`, namespace, strings.Join(names, " ")))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DaemonSets) Activate() error {
	byNamespace := lo.GroupBy(s.Items, func(item *DaemonSet) string {
		return item.Metadata.Namespace
	})

	for namespace, daemonSets := range byNamespace {
		names := lo.Map(daemonSets, func(daemonSet *DaemonSet, index int) string {
			return daemonSet.Metadata.Name
		})

		_, err := bash.RunAndLogWrite(fmt.Sprintf(`kubectl -n %s patch daemonset %s --type json -p='[{"op": "remove", "path": "/spec/template/spec/nodeSelector/non-existing"}]'`, namespace, strings.Join(names, " ")))
		if err != nil {
			return err
		}
	}

	return nil
}

type DaemonSet struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   DaemonSetMetadata `json:"metadata"`
	Spec       DaemonSetSpec     `json:"spec"`
	Status     Status            `json:"status"`
}

type DaemonSetMetadata struct {
	Annotations       PurpleAnnotations `json:"annotations"`
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	Generation        int64             `json:"generation"`
	Labels            PurpleLabels      `json:"labels"`
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	ResourceVersion   string            `json:"resourceVersion"`
	Uid               string            `json:"uid"`
}

type PurpleAnnotations struct {
	DeprecatedDaemonsetTemplateGeneration       string `json:"deprecated.daemonset.template.generation"`
	KubectlKubernetesIoLastAppliedConfiguration string `json:"kubectl.kubernetes.io/last-applied-configuration"`
}

type PurpleLabels struct {
	SkaffoldDevRunID string `json:"skaffold.dev/run-id"`
}

type DaemonSetSpec struct {
	RevisionHistoryLimit int64          `json:"revisionHistoryLimit"`
	Selector             Selector       `json:"selector"`
	Template             Template       `json:"template"`
	UpdateStrategy       UpdateStrategy `json:"updateStrategy"`
}

type Selector struct {
	MatchLabels MatchLabels `json:"matchLabels"`
}

type MatchLabels struct {
	App string `json:"app"`
}

type Template struct {
	Metadata TemplateMetadata `json:"metadata"`
	Spec     TemplateSpec     `json:"spec"`
}

type TemplateMetadata struct {
	Annotations       FluffyAnnotations `json:"annotations"`
	CreationTimestamp interface{}       `json:"creationTimestamp"`
	Labels            FluffyLabels      `json:"labels"`
}

type FluffyAnnotations struct {
	SchedulerAlphaKubernetesIoCriticalPod string `json:"scheduler.alpha.kubernetes.io/critical-pod"`
}

type FluffyLabels struct {
	App              string `json:"app"`
	SkaffoldDevRunID string `json:"skaffold.dev/run-id"`
}

type TemplateSpec struct {
	Affinity                      Affinity          `json:"affinity"`
	Containers                    []Container       `json:"containers"`
	DNSPolicy                     string            `json:"dnsPolicy"`
	HostNetwork                   bool              `json:"hostNetwork"`
	PriorityClassName             string            `json:"priorityClassName"`
	RestartPolicy                 string            `json:"restartPolicy"`
	SchedulerName                 string            `json:"schedulerName"`
	SecurityContext               SecurityContext   `json:"securityContext"`
	ServiceAccount                string            `json:"serviceAccount"`
	ServiceAccountName            string            `json:"serviceAccountName"`
	TerminationGracePeriodSeconds int64             `json:"terminationGracePeriodSeconds"`
	Tolerations                   []Toleration      `json:"tolerations"`
	Volumes                       []Volume          `json:"volumes"`
	NodeSelector                  map[string]string `json:"nodeSelector"`
}

type Affinity struct {
	NodeAffinity NodeAffinity `json:"nodeAffinity"`
}

type NodeAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution RequiredDuringSchedulingIgnoredDuringExecution `json:"requiredDuringSchedulingIgnoredDuringExecution"`
}

type RequiredDuringSchedulingIgnoredDuringExecution struct {
	NodeSelectorTerms []NodeSelectorTerm `json:"nodeSelectorTerms"`
}

type NodeSelectorTerm struct {
	MatchExpressions []MatchExpression `json:"matchExpressions"`
}

type MatchExpression struct {
	Key      string `json:"key"`
	Operator string `json:"operator"`
}

type SecretKeyRef struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Resources struct {
	Requests Requests `json:"requests"`
}

type Requests struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type VolumeMount struct {
	MountPath string `json:"mountPath"`
	Name      string `json:"name"`
	ReadOnly  bool   `json:"readOnly"`
}

type SecurityContext struct{}

type Toleration struct {
	Effect   *string `json:"effect,omitempty"`
	Key      string  `json:"key"`
	Value    *string `json:"value,omitempty"`
	Operator *string `json:"operator,omitempty"`
}

type Volume struct {
	ConfigMap ConfigMap `json:"configMap"`
	Name      string    `json:"name"`
}

type UpdateStrategy struct {
	RollingUpdate RollingUpdate `json:"rollingUpdate"`
	Type          string        `json:"type"`
}

type RollingUpdate struct {
	MaxSurge       int64 `json:"maxSurge"`
	MaxUnavailable int64 `json:"maxUnavailable"`
}
