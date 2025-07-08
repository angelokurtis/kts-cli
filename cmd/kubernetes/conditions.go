package kubernetes

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

// kube conditions -A pods
func conditions(cmd *cobra.Command, args []string) error {
	resources := ""
	if len(args) > 0 {
		resources = args[0]
	}

	mapper, err := newMapper()
	if err != nil {
		return err
	}

	list, err := kubectl.ListUnstructureds(resources, namespace, allNamespaces)
	if err != nil {
		return err
	}

	grouped := lo.GroupBy(list.Items, func(u unstructured.Unstructured) schema.GroupVersionKind {
		return u.GroupVersionKind()
	})

	keys := make([]schema.GroupVersionKind, 0, len(grouped))
	for k := range grouped {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		a, b := keys[i], keys[j]
		if a.Group != b.Group {
			return a.Group < b.Group
		}

		if a.Version != b.Version {
			return a.Version < b.Version
		}

		return a.Kind < b.Kind
	})

	for _, k := range keys {
		printTable(grouped[k], mapper)
		fmt.Println()
	}

	return nil
}

func printTable(unstructureds []unstructured.Unstructured, mapper meta.RESTMapper) {
	cols := make(map[string]struct{})
	var rows []map[string]string

	for _, obj := range unstructureds {
		row := map[string]string{}

		if ns := obj.GetNamespace(); ns != "" {
			cols["Namespace"] = struct{}{}
			row["Namespace"] = ns
		}

		cols["Name"] = struct{}{}
		row["Name"] = getFQN(&obj, mapper)

		if conds, err := getStatusConditions(&obj); err == nil {
			for _, cond := range conds {
				cols[cond.Type] = struct{}{}

				if len(cond.Reason) > 0 {
					row[cond.Type] = cond.Reason
				} else {
					row[cond.Type] = string(cond.Status)
				}
			}
		}

		rows = append(rows, row)
	}

	// Build header in desired order: Namespace, Name, then the rest sorted
	var headers []string
	if _, ok := cols["Namespace"]; ok {
		headers = append(headers, "Namespace")
	}

	if _, ok := cols["Name"]; ok {
		headers = append(headers, "Name")
	}

	var otherCols []string

	for col := range cols {
		if col != "Namespace" && col != "Name" {
			otherCols = append(otherCols, col)
		}
	}

	sort.Strings(otherCols)
	headers = append(headers, otherCols...)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader(headers)

	for _, r := range rows {
		line := make([]string, len(headers))
		for i, h := range headers {
			line[i] = r[h]
		}

		table.Append(line)
	}

	table.Render()
}

func newMapper() (meta.RESTMapper, error) {
	// Load kubeconfig from default location (~/.kube/config)
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get users' home: %w", err)
	}

	kubeconfig := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create discovery client: %w", err)
	}

	gr, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		return nil, fmt.Errorf("failed to get API group resources: %w", err)
	}

	mapper := restmapper.NewDiscoveryRESTMapper(gr)

	return mapper, nil
}

func getFQN(u *unstructured.Unstructured, mapper meta.RESTMapper) string {
	gvk := u.GroupVersionKind()

	// Use RESTMapper to find the resource
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		log.Fatalf("Failed to get REST mapping: %w", err)
	}

	name := u.GetName()

	resource := mapping.Resource.Resource // This gives plural + group
	if len(gvk.Group) > 0 {
		resource = fmt.Sprintf("%s.%s", resource, gvk.Group)
	}

	return fmt.Sprintf("%s/%s", resource, name)
}

func getStatusConditions(obj *unstructured.Unstructured) ([]metav1.Condition, error) {
	raw, found, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if err != nil {
		return nil, fmt.Errorf("failed to get status.conditions: %w", err)
	}

	if !found {
		return nil, fmt.Errorf("status.conditions not found")
	}

	var out []metav1.Condition

	for _, r := range raw {
		m, ok := r.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid condition format")
		}

		// Marshal to JSON
		b, err := json.Marshal(m)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal condition: %w", err)
		}

		// Unmarshal to metav1.Condition
		var cond metav1.Condition
		if err := json.Unmarshal(b, &cond); err != nil {
			return nil, fmt.Errorf("failed to unmarshal condition: %w", err)
		}

		out = append(out, cond)
	}

	return out, nil
}
