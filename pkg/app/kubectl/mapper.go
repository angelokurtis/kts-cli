package kubectl

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

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
